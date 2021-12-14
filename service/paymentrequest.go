package service

import (
	"context"

	"github.com/pkg/errors"
	validator "github.com/theflyingcodr/govalidator"

	"github.com/libsv/go-p4"
)

type paymentRequest struct {
	prRdr p4.PaymentRequestReader
}

// NewPaymentRequest will setup and return a new PaymentRequest service that will generate outputs
// using the provided outputter which is defined in server config.
func NewPaymentRequest(prRdr p4.PaymentRequestReader) *paymentRequest {
	return &paymentRequest{
		prRdr: prRdr,
	}
}

// PaymentRequest handles setting up a new PaymentRequest response and will validate that we have a paymentID.
func (p *paymentRequest) PaymentRequest(ctx context.Context, args p4.PaymentRequestArgs) (*p4.PaymentRequest, error) {
	if err := validator.New().
		Validate("paymentID", validator.NotEmpty(args.PaymentID)); err.Err() != nil {
		return nil, err
	}

	pReq, err := p.prRdr.PaymentRequest(ctx, args)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get payment request for paymentID %s", args.PaymentID)
	}
	if pReq.MerchantData != nil && pReq.MerchantData.ExtendedData == nil {
		pReq.MerchantData.ExtendedData = map[string]interface{}{
			"paymentReference": args.PaymentID,
		}
	}

	return pReq, nil
}
