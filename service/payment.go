package service

import (
	"context"

	"github.com/labstack/gommon/log"

	"github.com/libsv/go-p4"
)

// payment is a layer on top of the payment services of which we currently support:
// * wallet payments, that are handled by the wallet and transmitted to the network
// * paymail payments, that use the paymail protocol for making the payments.
type payment struct {
	paymentWtr p4.PaymentWriter
}

// NewPayment will create and return a new payment service.
func NewPayment(paymentWtr p4.PaymentWriter) *payment {
	return &payment{
		paymentWtr: paymentWtr,
	}
}

// PaymentCreate will setup a new payment and return the result.
func (p *payment) PaymentCreate(ctx context.Context, args p4.PaymentCreateArgs, req p4.PaymentCreate) (*p4.PaymentACK, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// broadcast it to a wallet for processing.
	if err := p.paymentWtr.PaymentCreate(ctx, args, req); err != nil {
		log.Error(err)
		return &p4.PaymentACK{
			Memo:  err.Error(),
			Error: 1,
		}, err
	}
	return &p4.PaymentACK{
		Payment: &req,
		Memo:    req.Memo,
	}, nil
}
