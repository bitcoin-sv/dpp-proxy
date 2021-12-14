package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/libsv/go-p4"
	p4mocks "github.com/libsv/go-p4/mocks"
	"github.com/libsv/p4-server/service"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRequest_PaymentRequest(t *testing.T) {
	created := time.Now()
	expired := created.Add(time.Hour * 24)
	tests := map[string]struct {
		paymentRequestFunc func(context.Context, p4.PaymentRequestArgs) (*p4.PaymentRequest, error)
		args               p4.PaymentRequestArgs
		expResp            *p4.PaymentRequest
		expErr             error
	}{
		"successful request": {
			args: p4.PaymentRequestArgs{
				PaymentID: "abc123",
			},
			paymentRequestFunc: func(context.Context, p4.PaymentRequestArgs) (*p4.PaymentRequest, error) {
				return &p4.PaymentRequest{
					SPVRequired:         false,
					CreationTimestamp:   created,
					ExpirationTimestamp: expired,
					Destinations: p4.PaymentDestinations{
						Outputs: []p4.Output{{
							Amount: 500,
							Script: "abc123",
						}},
					},
					PaymentURL: "http://iamsotest/api/v1/payment/abc123",
					Memo:       "invoice abc123",
					MerchantData: &p4.Merchant{
						ExtendedData: map[string]interface{}{"paymentReference": "abc123"},
					},
				}, nil
			},
			expResp: &p4.PaymentRequest{
				SPVRequired:         false,
				CreationTimestamp:   created,
				ExpirationTimestamp: expired,
				Destinations: p4.PaymentDestinations{
					Outputs: []p4.Output{{
						Amount: 500,
						Script: "abc123",
					}},
				},
				PaymentURL: "http://iamsotest/api/v1/payment/abc123",
				Memo:       "invoice abc123",
				MerchantData: &p4.Merchant{
					ExtendedData: map[string]interface{}{"paymentReference": "abc123"},
				},
			},
		},
		"successful request with nil extended data": {
			args: p4.PaymentRequestArgs{
				PaymentID: "abc123",
			},
			paymentRequestFunc: func(context.Context, p4.PaymentRequestArgs) (*p4.PaymentRequest, error) {
				return &p4.PaymentRequest{
					SPVRequired:         false,
					CreationTimestamp:   created,
					ExpirationTimestamp: expired,
					Destinations: p4.PaymentDestinations{
						Outputs: []p4.Output{{
							Amount: 500,
							Script: "abc123",
						}},
					},
					MerchantData: &p4.Merchant{},
					PaymentURL:   "http://iamsotest/api/v1/payment/abc123",
					Memo:         "invoice abc123",
				}, nil
			},
			expResp: &p4.PaymentRequest{
				SPVRequired:         false,
				CreationTimestamp:   created,
				ExpirationTimestamp: expired,
				Destinations: p4.PaymentDestinations{
					Outputs: []p4.Output{{
						Amount: 500,
						Script: "abc123",
					}},
				},
				PaymentURL: "http://iamsotest/api/v1/payment/abc123",
				Memo:       "invoice abc123",
				MerchantData: &p4.Merchant{
					ExtendedData: map[string]interface{}{"paymentReference": "abc123"},
				},
			},
		},
		"invalid args rejected": {
			expErr: errors.New("[paymentID: value cannot be empty]"),
		},
		"payment request reader error handled and reported": {
			args: p4.PaymentRequestArgs{
				PaymentID: "abc123",
			},
			paymentRequestFunc: func(context.Context, p4.PaymentRequestArgs) (*p4.PaymentRequest, error) {
				return nil, errors.New("oh boi")
			},
			expErr: errors.New("failed to get payment request for paymentID abc123: oh boi"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svc := service.NewPaymentRequest(&p4mocks.PaymentRequestServiceMock{
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
