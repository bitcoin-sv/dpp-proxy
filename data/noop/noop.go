package noop

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/libsv/go-bt/v2"

	"github.com/libsv/pptcl"
)

type noop struct {
}

func NewNoOp() *noop {
	log.Info("using NOOP data store")
	return &noop{}
}

// PaymentCreate will post a request to payd to validate and add the txos to the wallet.
//
// If invalid a non 204 status code is returned.
func (n *noop) PaymentCreate(ctx context.Context, args pptcl.PaymentCreateArgs, req pptcl.PaymentCreate) error {
	log.Info("hit noop.PaymentCreate")
	return nil
}

// Owner will return information regarding the owner of a payd wallet.
//
// In this example, the payd wallet has no auth, in proper implementations auth would
// be enabled and a cookie / oauth / bearer token etc would be passed down.
func (n *noop) Owner(ctx context.Context) (*pptcl.MerchantData, error) {
	log.Info("hit noop.Owner")
	return &pptcl.MerchantData{
		AvatarURL:    "noop",
		MerchantName: "noop",
		Email:        "noop",
		Address:      "noop",
		ExtendedData: nil,
	}, nil
}

// Outputs will return outputs for payment requests, the sender will then fulfil these outputs
// and send a tx for broadcast.
func (n *noop) Outputs(ctx context.Context, args pptcl.PaymentRequestArgs) ([]pptcl.Output, error) {
	log.Info("hit noop.Outputs")
	return []pptcl.Output{
		{
			Amount:      0,
			Script:      "noop",
			Description: "noop",
		},
	}, nil
}

// Fees will return current fees that a payd wallet is using.
func (n *noop) Fees(ctx context.Context) (*bt.FeeQuote, error) {
	log.Info("hit noop.Fees")
	return bt.NewFeeQuote(), nil
}
