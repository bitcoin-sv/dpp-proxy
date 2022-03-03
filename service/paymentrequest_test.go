package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/libsv/dpp-proxy/service"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-dpp"
	dppMocks "github.com/libsv/go-dpp/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRequest_PaymentRequest(t *testing.T) {
	created := time.Now()
	expired := created.Add(time.Hour * 24)
	tests := map[string]struct {
		paymentRequestFunc func(context.Context, dpp.PaymentRequestArgs) (*dpp.PaymentRequest, error)
		args               dpp.PaymentRequestArgs
		expResp            *dpp.PaymentRequest
		expErr             error
	}{
		"successful request": {
			args: dpp.PaymentRequestArgs{
				PaymentID: "abc123",
			},
			paymentRequestFunc: func(context.Context, dpp.PaymentRequestArgs) (*dpp.PaymentRequest, error) {
				return &dpp.PaymentRequest{
					SPVRequired:         false,
					CreationTimestamp:   created,
					ExpirationTimestamp: expired,
					Destinations: dpp.PaymentDestinations{
						Outputs: []dpp.Output{{
							Amount: 500,
							LockingScript: func() *bscript.Script {
								ls, err := bscript.NewFromHexString("abc123")
								assert.NoError(t, err)
								return ls
							}(),
						}},
					},
					PaymentURL: "http://iamsotest/api/v1/payment/abc123",
					Memo:       "invoice abc123",
					MerchantData: &dpp.Merchant{
						ExtendedData: map[string]interface{}{"paymentReference": "abc123"},
					},
				}, nil
			},
			expResp: &dpp.PaymentRequest{
				SPVRequired:         false,
				CreationTimestamp:   created,
				ExpirationTimestamp: expired,
				Destinations: dpp.PaymentDestinations{
					Outputs: []dpp.Output{{
						Amount: 500,
						LockingScript: func() *bscript.Script {
							ls, err := bscript.NewFromHexString("abc123")
							assert.NoError(t, err)
							return ls
						}(),
					}},
				},
				PaymentURL: "http://iamsotest/api/v1/payment/abc123",
				Memo:       "invoice abc123",
				MerchantData: &dpp.Merchant{
					ExtendedData: map[string]interface{}{"paymentReference": "abc123"},
				},
			},
		},
		"successful request with nil extended data": {
			args: dpp.PaymentRequestArgs{
				PaymentID: "abc123",
			},
			paymentRequestFunc: func(context.Context, dpp.PaymentRequestArgs) (*dpp.PaymentRequest, error) {
				return &dpp.PaymentRequest{
					SPVRequired:         false,
					CreationTimestamp:   created,
					ExpirationTimestamp: expired,
					Destinations: dpp.PaymentDestinations{
						Outputs: []dpp.Output{{
							Amount: 500,
							LockingScript: func() *bscript.Script {
								ls, err := bscript.NewFromHexString("abc123")
								assert.NoError(t, err)
								return ls
							}(),
						}},
					},
					MerchantData: &dpp.Merchant{},
					PaymentURL:   "http://iamsotest/api/v1/payment/abc123",
					Memo:         "invoice abc123",
				}, nil
			},
			expResp: &dpp.PaymentRequest{
				SPVRequired:         false,
				CreationTimestamp:   created,
				ExpirationTimestamp: expired,
				Destinations: dpp.PaymentDestinations{
					Outputs: []dpp.Output{{
						Amount: 500,
						LockingScript: func() *bscript.Script {
							ls, err := bscript.NewFromHexString("abc123")
							assert.NoError(t, err)
							return ls
						}(),
					}},
				},
				PaymentURL: "http://iamsotest/api/v1/payment/abc123",
				Memo:       "invoice abc123",
				MerchantData: &dpp.Merchant{
					ExtendedData: map[string]interface{}{"paymentReference": "abc123"},
				},
			},
		},
		"invalid args rejected": {
			expErr: errors.New("[paymentID: value cannot be empty]"),
		},
		"payment request reader error handled and reported": {
			args: dpp.PaymentRequestArgs{
				PaymentID: "abc123",
			},
			paymentRequestFunc: func(context.Context, dpp.PaymentRequestArgs) (*dpp.PaymentRequest, error) {
				return nil, errors.New("oh boi")
			},
			expErr: errors.New("failed to get payment request for paymentID abc123: oh boi"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svc := service.NewPaymentRequest(&dppMocks.PaymentRequestServiceMock{
				PaymentRequestFunc: test.paymentRequestFunc,
			})

			resp, err := svc.PaymentRequest(context.TODO(), test.args)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, *test.expResp, *resp)
		})
	}
}
