package noop

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/libsv/go-bt/v2"

	"github.com/libsv/go-p4"
)

type noop struct {
}

// NewNoOp will setup and return a new no operational data store for
// testing purposes. Useful if you want to explore endpoints without
// integrating with a wallet.
func NewNoOp() *noop {
	log.Info("using NOOP data store")
	return &noop{}
}

// PaymentCreate will post a request to payd to validate and add the txos to the wallet.
//
// If invalid a non 204 status code is returned.
func (n *noop) PaymentCreate(ctx context.Context, args p4.PaymentCreateArgs, req p4.PaymentCreate) error {
	log.Info("hit noop.PaymentCreate")
	return nil
}

// Owner will return information regarding the owner of a payd wallet.
//
// In this example, the payd wallet has no auth, in proper implementations auth would
// be enabled and a cookie / oauth / bearer token etc would be passed down.
func (n *noop) Owner(ctx context.Context) (*p4.MerchantData, error) {
	log.Info("hit noop.Owner")
	return &p4.MerchantData{
		AvatarURL:    "noop",
		MerchantName: "noop",
		Email:        "noop",
		Address:      "noop",
		ExtendedData: nil,
	}, nil
}

func (n *noop) Destinations(ctx context.Context, args p4.PaymentRequestArgs) (*p4.Destinations, error) {
	log.Info("hit noop.Destinations")
	return &p4.Destinations{
		Outputs: []p4.Output{{
			Amount:      0,
			Script:      "noop",
			Description: "noop",
		}},
		Fees: bt.NewFeeQuote(),
	}, nil
}
