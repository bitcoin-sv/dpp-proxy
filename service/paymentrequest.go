package service

import (
	"context"

	"github.com/libsv/go-dpp"
	"github.com/pkg/errors"
	validator "github.com/theflyingcodr/govalidator"
)

type paymentRequest struct {
	prRdr dpp.PaymentRequestReader
}

// NewPaymentRequest will setup and return a new PaymentRequest service that will generate outputs
// using the provided outputter which is defined in server config.
func NewPaymentRequest(prRdr dpp.PaymentRequestReader) *paymentRequest {
	return &paymentRequest{
		prRdr: prRdr,
	}
}

// PaymentRequest handles setting up a new PaymentRequest response and will validate that we have a paymentID.
func (p *paymentRequest) PaymentRequest(ctx context.Context, args dpp.PaymentRequestArgs) (*dpp.PaymentRequest, error) {
	if err := validator.New().
		Validate("paymentID", validator.NotEmpty(args.PaymentID)); err.Err() != nil {
		return nil, err
	}

	pReq, err := p.prRdr.PaymentRequest(ctx, args)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get payment request for paymentID %s", args.PaymentID)
	}
	if pReq.Beneficiary != nil && pReq.Beneficiary.ExtendedData == nil {
		pReq.Beneficiary.ExtendedData = map[string]interface{}{
			"paymentReference": args.PaymentID,
		}
	}

	return pReq, nil
}
