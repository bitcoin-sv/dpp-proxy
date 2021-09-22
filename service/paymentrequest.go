package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	validator "github.com/theflyingcodr/govalidator"

	"github.com/libsv/pptcl"
	"github.com/libsv/pptcl/config"
)

type paymentRequest struct {
	walletCfg   *config.Server
	outputter   pptcl.OutputReader
	merchantRdr pptcl.MerchantReader
	feeRdr      pptcl.FeeReader
}

// NewPaymentRequest will setup and return a new PaymentRequest service that will generate outputs
// using the provided outputter which is defined in server config.
func NewPaymentRequest(walletCfg *config.Server, outputter pptcl.OutputReader, merchantRdr pptcl.MerchantReader, feeRdr pptcl.FeeReader) *paymentRequest {
	return &paymentRequest{
		walletCfg:   walletCfg,
		outputter:   outputter,
		merchantRdr: merchantRdr,
		feeRdr:      feeRdr,
	}
}

// CreatePaymentRequest handles setting up a new PaymentRequest response and will validate that we have a paymentID.
func (p *paymentRequest) CreatePaymentRequest(ctx context.Context, args pptcl.PaymentRequestArgs) (*pptcl.PaymentRequest, error) {
	if err := validator.New().
		Validate("paymentID", validator.NotEmpty(args.PaymentID)); err.Err() != nil {
		return nil, err
	}

	// get payment destinations from merchant
	oo, err := p.outputter.Outputs(ctx, args)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate outputs for paymentID %s", args.PaymentID)
	}
	// get fees from merchant
	fees, err := p.feeRdr.Fees(ctx, args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read fees when constructing payment request")
	}
	// get merchant information
	merchant, err := p.merchantRdr.Owner(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read merchant data when constructing payment request")
	}
	if merchant.ExtendedData == nil {
		merchant.ExtendedData = map[string]interface{}{}
	}
	// here we store paymentRef in extended data to allow some validation in payment flow
	merchant.ExtendedData["paymentReference"] = args.PaymentID
	return &pptcl.PaymentRequest{
		Network:             "mainnet",
		Outputs:             oo,
		CreationTimestamp:   time.Now().UTC(),
		ExpirationTimestamp: time.Now().Add(24 * time.Hour).UTC(),
		PaymentURL:          fmt.Sprintf("http://%s%s/api/v1/payment/%s", p.walletCfg.Hostname, p.walletCfg.Port, args.PaymentID),
		Memo:                fmt.Sprintf("invoice %s", args.PaymentID),
		MerchantData: &pptcl.MerchantData{
			AvatarURL:    merchant.AvatarURL,
			MerchantName: merchant.MerchantName,
			Email:        merchant.Email,
			Address:      merchant.Address,
			ExtendedData: merchant.ExtendedData,
		},
		FeeRate: fees,
	}, nil
}
