package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/libsv/go-bc/spv"
	"github.com/libsv/go-p4"
	"github.com/libsv/go-p4/mocks"
	"github.com/libsv/go-p4/service"
	"github.com/stretchr/testify/assert"
)

func TestPayment_Create(t *testing.T) {
	tests := map[string]struct {
		paymentCreateFn func(context.Context, p4.PaymentCreateArgs, p4.PaymentCreate) error
		args            p4.PaymentCreateArgs
		req             p4.PaymentCreate
		expResp         *p4.PaymentACK
		expErr          error
	}{
		"successful payment create": {
			paymentCreateFn: func(context.Context, p4.PaymentCreateArgs, p4.PaymentCreate) error {
				return nil
			},
			args: p4.PaymentCreateArgs{
				PaymentID: "abc123",
			},
			req: p4.PaymentCreate{
				SPVEnvelope: &spv.Envelope{
					RawTx: "01000000000000000000",
					TxID:  "d21633ba23f70118185227be58a63527675641ad37967e2aa461559f577aec43",
				},
				MerchantData: p4.MerchantData{
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			expResp: &p4.PaymentACK{
				Payment: &p4.PaymentCreate{
					SPVEnvelope: &spv.Envelope{
						RawTx: "01000000000000000000",
						TxID:  "d21633ba23f70118185227be58a63527675641ad37967e2aa461559f577aec43",
					},
					MerchantData: p4.MerchantData{
						ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
					},
				},
			},
		},
		"invalid args errors": {
			paymentCreateFn: func(context.Context, p4.PaymentCreateArgs, p4.PaymentCreate) error {
				return nil
			},
			args: p4.PaymentCreateArgs{},
			req: p4.PaymentCreate{
				SPVEnvelope: &spv.Envelope{
					RawTx: "01000000000000000000",
					TxID:  "d21633ba23f70118185227be58a63527675641ad37967e2aa461559f577aec43",
				},
				MerchantData: p4.MerchantData{
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			expErr: errors.New("[paymentID: value cannot be empty]"),
		},
		"missing raw tx errors": {
			paymentCreateFn: func(context.Context, p4.PaymentCreateArgs, p4.PaymentCreate) error {
				return nil
			},
			args: p4.PaymentCreateArgs{
				PaymentID: "abc123",
			},
			req: p4.PaymentCreate{
				MerchantData: p4.MerchantData{
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			expErr: errors.New("[spvEnvelope/rawTx: either an SPVEnvelope or a rawTX are required]"),
		},
		"error on payment create is handled": {
			paymentCreateFn: func(context.Context, p4.PaymentCreateArgs, p4.PaymentCreate) error {
				return errors.New("lol oh boi")
			},
			args: p4.PaymentCreateArgs{
				PaymentID: "abc123",
			},
			req: p4.PaymentCreate{
				SPVEnvelope: &spv.Envelope{
					RawTx: "01000000000000000000",
					TxID:  "d21633ba23f70118185227be58a63527675641ad37967e2aa461559f577aec43",
				},
				MerchantData: p4.MerchantData{
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			expResp: &p4.PaymentACK{
				Memo:  "lol oh boi",
				Error: 1,
			},
			expErr: errors.New("lol oh boi"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svc := service.NewPayment(&mocks.PaymentWriterMock{
				PaymentCreateFunc: test.paymentCreateFn,
			})

			ack, err := svc.PaymentCreate(context.TODO(), test.args, test.req)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expResp != nil {
				assert.NotNil(t, ack)
				assert.Equal(t, *test.expResp, *ack)
			} else {
				assert.Nil(t, ack)
			}
		})
	}
}
