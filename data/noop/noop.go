package noop

import (
	"context"
	"time"

	"github.com/libsv/go-bt/v2"
	"github.com/libsv/p4-server/log"

	"github.com/libsv/go-p4"
)

type noop struct {
	l log.Logger
}

// NewNoOp will setup and return a new no operational data store for
// testing purposes. Useful if you want to explore endpoints without
// integrating with a wallet.
func NewNoOp(l log.Logger) *noop {
	l.Info("using NOOP data store")
	return &noop{}
}

// PaymentCreate will post a request to payd to validate and add the txos to the wallet.
//
// If invalid a non 204 status code is returned.
func (n *noop) PaymentCreate(ctx context.Context, args p4.PaymentCreateArgs, req p4.Payment) (*p4.PaymentACK, error) {
	n.l.Info("hit noop.PaymentCreate")
	return &p4.PaymentACK{}, nil
}

func (n noop) PaymentRequest(ctx context.Context, args p4.PaymentRequestArgs) (*p4.PaymentRequest, error) {
	return &p4.PaymentRequest{
		Network:             "noop",
		CreationTimestamp:   time.Now(),
		ExpirationTimestamp: time.Now().Add(time.Hour),
		FeeRate: func() *bt.FeeQuote {
			fq := bt.NewFeeQuote()
			fq.UpdateExpiry(time.Now().Add(10 * time.Hour))
			return fq
		}(),
		Memo:        "noop",
		PaymentURL:  "noop",
		SPVRequired: true,
		MerchantData: &p4.Merchant{
			AvatarURL:    "noop",
			Name:         "noop",
			Email:        "noop",
			Address:      "noop",
			ExtendedData: nil,
		},
		Destinations: p4.PaymentDestinations{
			Outputs: []p4.Output{{
				Amount:      0,
				Script:      "noop",
				Description: "noop",
			}},
		},
	}, nil
}
