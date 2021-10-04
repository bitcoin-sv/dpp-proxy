package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	validator "github.com/theflyingcodr/govalidator"

	"github.com/libsv/go-p4"
	"github.com/libsv/go-p4/config"
)

type paymentRequest struct {
	walletCfg   *config.Server
	destRdr     p4.DestinationReader
	merchantRdr p4.MerchantReader
}

// NewPaymentRequest will setup and return a new PaymentRequest service that will generate outputs
// using the provided outputter which is defined in server config.
func NewPaymentRequest(walletCfg *config.Server, destRdr p4.DestinationReader, merchantRdr p4.MerchantReader) *paymentRequest {
	return &paymentRequest{
		walletCfg:   walletCfg,
		destRdr:     destRdr,
		merchantRdr: merchantRdr,
	}
}

// CreatePaymentRequest handles setting up a new PaymentRequest response and will validate that we have a paymentID.
func (p *paymentRequest) CreatePaymentRequest(ctx context.Context, args p4.PaymentRequestArgs) (*p4.PaymentRequest, error) {
	if err := validator.New().
		Validate("paymentID", validator.NotEmpty(args.PaymentID)); err.Err() != nil {
		return nil, err
	}

	dests, err := p.destRdr.Destinations(ctx, args)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to geet destinations for paymentID %s", args.PaymentID)
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
	return &p4.PaymentRequest{
		Network:             "mainnet",
		SPVRequired:         dests.SPVRequired,
		Destinations:        p4.PaymentDestinations{Outputs: dests.Outputs},
		FeeRate:             dests.Fees,
		CreationTimestamp:   time.Now().UTC(),
		ExpirationTimestamp: time.Now().Add(24 * time.Hour).UTC(),
		PaymentURL:          fmt.Sprintf("http://%s/api/v1/payment/%s", p.walletCfg.FQDN, args.PaymentID),
		Memo:                fmt.Sprintf("invoice %s", args.PaymentID),
		MerchantData: &p4.MerchantData{
			AvatarURL:    merchant.AvatarURL,
			MerchantName: merchant.MerchantName,
			Email:        merchant.Email,
			Address:      merchant.Address,
			ExtendedData: merchant.ExtendedData,
		},
	}, nil
}
