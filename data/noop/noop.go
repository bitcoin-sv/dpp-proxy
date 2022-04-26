package noop

import (
	"context"
	"time"

	"github.com/libsv/dp3/log"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"

	"github.com/libsv/go-dpp"
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
func (n *noop) PaymentCreate(ctx context.Context, args dpp.PaymentCreateArgs, req dpp.Payment) (*dpp.PaymentACK, error) {
	n.l.Info("hit noop.PaymentCreate")
	return &dpp.PaymentACK{}, nil
}

func (n noop) PaymentRequest(ctx context.Context, args dpp.PaymentRequestArgs) (*dpp.PaymentRequest, error) {
	return &dpp.PaymentRequest{
		Network:             "noop",
		CreationTimestamp:   time.Now(),
		ExpirationTimestamp: time.Now().Add(time.Hour),
		FeeRate: func() *bt.FeeQuote {
			fq := bt.NewFeeQuote()
			fq.UpdateExpiry(time.Now().Add(10 * time.Hour))
			return fq
		}(),
		Memo:             "noop",
		PaymentURL:       "noop",
		AncestryRequired: true,
		MerchantData: &dpp.Merchant{
			AvatarURL:    "noop",
			Name:         "noop",
			Email:        "noop",
			Address:      "noop",
			ExtendedData: nil,
		},
		Destinations: dpp.PaymentDestinations{
			Outputs: []dpp.Output{{
				Amount:        0,
				LockingScript: &bscript.Script{},
				Description:   "noop",
			}},
		},
	}, nil
}
