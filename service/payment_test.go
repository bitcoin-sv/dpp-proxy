package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/libsv/dpp-proxy/log"
	"github.com/libsv/dpp-proxy/service"
	"github.com/libsv/go-bc/spv"
	"github.com/libsv/go-dpp"
	dppMocks "github.com/libsv/go-dpp/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPayment_Create(t *testing.T) {
	tests := map[string]struct {
		paymentCreateFn func(context.Context, dpp.PaymentCreateArgs, dpp.Payment) (*dpp.PaymentACK, error)
		args            dpp.PaymentCreateArgs
		req             dpp.Payment
		expErr          error
	}{
		"successful payment create": {
			paymentCreateFn: func(context.Context, dpp.PaymentCreateArgs, dpp.Payment) (*dpp.PaymentACK, error) {
				return &dpp.PaymentACK{}, nil
			},
			req: dpp.Payment{
				SPVEnvelope: &spv.Envelope{
					RawTx: "01000000000000000000",
					TxID:  "d21633ba23f70118185227be58a63527675641ad37967e2aa461559f577aec43",
				},
				MerchantData: dpp.Merchant{
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			args: dpp.PaymentCreateArgs{
				PaymentID: "abc123",
			},
		},
		"invalid args errors": {
			paymentCreateFn: func(context.Context, dpp.PaymentCreateArgs, dpp.Payment) (*dpp.PaymentACK, error) {
				return &dpp.PaymentACK{}, nil
			},
			args: dpp.PaymentCreateArgs{},
			req: dpp.Payment{
				SPVEnvelope: &spv.Envelope{
					RawTx: "01000000000000000000",
					TxID:  "d21633ba23f70118185227be58a63527675641ad37967e2aa461559f577aec43",
				},
				MerchantData: dpp.Merchant{
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			expErr: errors.New("[paymentID: value cannot be empty]"),
		},
		"missing raw tx errors": {
			paymentCreateFn: func(context.Context, dpp.PaymentCreateArgs, dpp.Payment) (*dpp.PaymentACK, error) {
				return &dpp.PaymentACK{}, nil
			},
			args: dpp.PaymentCreateArgs{
				PaymentID: "abc123",
			},
			req: dpp.Payment{
				MerchantData: dpp.Merchant{
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			expErr: errors.New("[spvEnvelope/rawTx: either an SPVEnvelope or a rawTX are required]"),
		},
		"error on payment create is handled": {
			paymentCreateFn: func(context.Context, dpp.PaymentCreateArgs, dpp.Payment) (*dpp.PaymentACK, error) {
				return nil, errors.New("lol oh boi")
			},
			args: dpp.PaymentCreateArgs{
				PaymentID: "abc123",
			},
			req: dpp.Payment{
				SPVEnvelope: &spv.Envelope{
					RawTx: "01000000000000000000",
					TxID:  "d21633ba23f70118185227be58a63527675641ad37967e2aa461559f577aec43",
				},
				MerchantData: dpp.Merchant{
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			expErr: errors.New("lol oh boi"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svc := service.NewPayment(
				log.Noop{},
				&dppMocks.PaymentWriterMock{
					PaymentCreateFunc: test.paymentCreateFn,
				})

			_, err := svc.PaymentCreate(context.TODO(), test.args, test.req)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
