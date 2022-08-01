package noop

import (
	"context"
	"github.com/libsv/go-dpp/modes/hybridmode"
	"time"

	"github.com/bitcoin-sv/dpp-proxy/log"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-dpp"
	"github.com/libsv/go-dpp/nativetypes"
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

func (n noop) PaymentTerms(ctx context.Context, args dpp.PaymentTermsArgs) (*dpp.PaymentTerms, error) {
	return &dpp.PaymentTerms{
		Network:             "noop",
		Version:			 "1.0",
		CreationTimestamp:   time.Now().Unix(),
		ExpirationTimestamp: time.Now().Add(time.Hour).Unix(),
		Memo:                "noop",
		PaymentURL:          "noop",
		Beneficiary: &dpp.Beneficiary{
			AvatarURL:        "noop",
			Name:             "noop",
			Email:            "noop",
			Address:          "noop",
			ExtendedData:     nil,
			PaymentReference: "noop",
		},
		Modes: &dpp.PaymentTermsModes{
			Hybrid: hybridmode.PaymentTerms{
				"choiceID0": {
					"transactions": {
						hybridmode.TransactionTerms{
							Outputs: hybridmode.Outputs{ NativeOutputs: []nativetypes.NativeOutput{
								{
									Amount:        1000,
									LockingScript: func() *bscript.Script {
										ls, _ := bscript.NewFromHexString(
											"76a91493d0d43918a5df78f08cfe22a4e022846b6736c288ac")
										return ls
									}(),
									Description:   "noop description",
								},
							} },
							Inputs: hybridmode.Inputs{},
							Policies: &hybridmode.Policies{},
						},
					},
				},

			},
		},
	}, nil
}
