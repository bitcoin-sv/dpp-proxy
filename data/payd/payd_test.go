package payd_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/libsv/dpp-proxy/config"
	"github.com/libsv/dpp-proxy/data/payd"
	"github.com/libsv/dpp-proxy/data/payd/models"
	"github.com/libsv/dpp-proxy/mocks"
	"github.com/libsv/go-bc/spv"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-dpp"
	"github.com/stretchr/testify/assert"
)

func TestPayd_PaymentCreate(t *testing.T) {
	tests := map[string]struct {
		doFunc func(context.Context, string, string, int, interface{}, interface{}) error
		args   dpp.PaymentCreateArgs
		req    dpp.Payment
		cfg    *config.PayD
		expURL string
		expReq models.PayDPaymentRequest
		expErr error
	}{
		"successful payment created": {
			doFunc: func(context.Context, string, string, int, interface{}, interface{}) error {
				return nil
			},
			args: dpp.PaymentCreateArgs{
				PaymentID: "qwe123",
			},
			req: dpp.Payment{
				RawTX:       func() *string { s := "rawrawraw"; return &s }(),
				SPVEnvelope: &spv.Envelope{},
				ProofCallbacks: map[string]dpp.ProofCallback{
					"abc.com": {Token: "mYtOkEn"},
				},
			},
			expReq: models.PayDPaymentRequest{
				RawTX:       func() *string { s := "rawrawraw"; return &s }(),
				SPVEnvelope: &spv.Envelope{},
				ProofCallbacks: map[string]dpp.ProofCallback{
					"abc.com": {Token: "mYtOkEn"},
				},
			},
			cfg: &config.PayD{
				Host: "paydhost",
				Port: ":8080",
			},
			expURL: "http://paydhost:8080/api/v1/payments/qwe123",
		},
		"successful https payment created": {
			doFunc: func(context.Context, string, string, int, interface{}, interface{}) error {
				return nil
			},
			args: dpp.PaymentCreateArgs{
				PaymentID: "qwe123",
			},
			req: dpp.Payment{
				RawTX:       func() *string { s := "rawrawraw"; return &s }(),
				SPVEnvelope: &spv.Envelope{},
				ProofCallbacks: map[string]dpp.ProofCallback{
					"abc.com": {Token: "mYtOkEn"},
				},
			},
			expReq: models.PayDPaymentRequest{
				RawTX:       func() *string { s := "rawrawraw"; return &s }(),
				SPVEnvelope: &spv.Envelope{},
				ProofCallbacks: map[string]dpp.ProofCallback{
					"abc.com": {Token: "mYtOkEn"},
				},
			},
			cfg: &config.PayD{
				Host:   "securepaydhost",
				Port:   ":8081",
				Secure: true,
			},
			expURL: "https://securepaydhost:8081/api/v1/payments/qwe123",
		},
		"error is handled and returned": {
			doFunc: func(context.Context, string, string, int, interface{}, interface{}) error {
				return errors.New("i tried so hard")
			},
			args: dpp.PaymentCreateArgs{
				PaymentID: "qwe123",
			},
			req: dpp.Payment{
				RawTX:       func() *string { s := "rawrawraw"; return &s }(),
				SPVEnvelope: &spv.Envelope{},
				ProofCallbacks: map[string]dpp.ProofCallback{
					"abc.com": {Token: "mYtOkEn"},
				},
			},
			expReq: models.PayDPaymentRequest{
				RawTX:       func() *string { s := "rawrawraw"; return &s }(),
				SPVEnvelope: &spv.Envelope{},
				ProofCallbacks: map[string]dpp.ProofCallback{
					"abc.com": {Token: "mYtOkEn"},
				},
			},
			cfg: &config.PayD{
				Host:   "securepaydhost",
				Port:   ":8081",
				Secure: true,
			},
			expURL: "https://securepaydhost:8081/api/v1/payments/qwe123",
			expErr: errors.New("i tried so hard"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			pd := payd.NewPayD(test.cfg, &mocks.HTTPClientMock{
				DoFunc: func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error {
					assert.Equal(t, test.expURL, url)
					assert.Equal(t, test.expReq, req)
					return test.doFunc(ctx, method, url, statusCode, req, out)
				},
			})
			_, err := pd.PaymentCreate(context.Background(), test.args, test.req)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPayd_PaymentRequest(t *testing.T) {
	tests := map[string]struct {
		doFunc        func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error
		args          dpp.PaymentRequestArgs
		cfg           *config.PayD
		expURL        string
		expPaymentReq *dpp.PaymentRequest
		expErr        error
	}{
		"successful payment request": {
			doFunc: func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error {
				return json.Unmarshal([]byte(`
					{
						"network": "mainnet",
						"spvRequired": true,
						"destinations": {
							"outputs": [
								{
									"amount": 100,
									"script": "525252"
								},
								{
									"amount": 400,
									"script": "535353"
								}
							]
						},
						"creationTimestamp": "2021-12-13T10:37:15.7946831Z",
						"expirationTimestamp": "2021-12-14T10:37:15.7946831Z",
						"paymentUrl": "http://dpp:8445/api/v1/payment/6K9oZq9",
						"memo": "invoice 6K9oZq9",
						"merchantData": {
							"avatar": "http://url.com",
							"name": "Merchant Name",
							"email": "merchant@demo.com",
							"address": "123 Street Fake",
							"extendedData": {
								"dislikes": "trying to think up placeholder data",
								"likes": "walks in the park at night",
								"paymentReference": "6K9oZq9"
							}
						},
						"fees": {
							"data": {
								"miningFee": {
									"satoshis": 5,
									"bytes": 10
								},
								"relayFee": {
									"satoshis": 5,
									"bytes": 10
								}
							},
							"standard": {
								"miningFee": {
									"satoshis": 5,
									"bytes": 10
								},
								"relayFee": {
									"satoshis": 5,
									"bytes": 10
								}
							}
						}
					}
				`), &out)
			},
			args: dpp.PaymentRequestArgs{
				PaymentID: "qwe123",
			},
			cfg: &config.PayD{
				Host: "payddest",
				Port: ":445",
			},
			expURL: "http://payddest:445/api/v1/payments/qwe123",
			expPaymentReq: &dpp.PaymentRequest{
				SPVRequired: true,
				Network:     "mainnet",
				Memo:        "invoice 6K9oZq9",
				PaymentURL:  "http://dpp:8445/api/v1/payment/6K9oZq9",
				Destinations: dpp.PaymentDestinations{
					Outputs: []dpp.Output{{
						LockingScript: func() *bscript.Script {
							ls, err := bscript.NewFromHexString("525252")
							assert.NoError(t, err)
							return ls
						}(),
						Amount: 100,
					}, {
						LockingScript: func() *bscript.Script {
							ls, err := bscript.NewFromHexString("535353")
							assert.NoError(t, err)
							return ls
						}(),
						Amount: 400,
					}},
				},
				MerchantData: &dpp.Merchant{
					AvatarURL: "http://url.com",
					Name:      "Merchant Name",
					Email:     "merchant@demo.com",
					Address:   "123 Street Fake",
					ExtendedData: map[string]interface{}{
						"dislikes":         "trying to think up placeholder data",
						"likes":            "walks in the park at night",
						"paymentReference": "6K9oZq9",
					},
				},
				FeeRate:             bt.NewFeeQuote(),
				CreationTimestamp:   func() time.Time { t, _ := time.Parse(time.RFC3339, "2021-12-13T10:37:15.7946831Z"); return t }(),
				ExpirationTimestamp: func() time.Time { t, _ := time.Parse(time.RFC3339, "2021-12-14T10:37:15.7946831Z"); return t }(),
			},
		},
		"successful https payment request": {
			doFunc: func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error {
				return json.Unmarshal([]byte(`
					{
						"network": "mainnet",
						"spvRequired": true,
						"destinations": {
							"outputs": [
								{
									"amount": 100,
									"script": "525252"
								},
								{
									"amount": 400,
									"script": "535353"
								}
							]
						},
						"creationTimestamp": "2021-12-13T10:37:15.7946831Z",
						"expirationTimestamp": "2021-12-14T10:37:15.7946831Z",
						"paymentUrl": "https://dpp:8445/api/v1/payment/6K9oZq9",
						"memo": "invoice 6K9oZq9",
						"merchantData": {
							"avatar": "http://url.com",
							"name": "Merchant Name",
							"email": "merchant@demo.com",
							"address": "123 Street Fake",
							"extendedData": {
								"dislikes": "trying to think up placeholder data",
								"likes": "walks in the park at night",
								"paymentReference": "6K9oZq9"
							}
						},
						"fees": {
							"data": {
								"miningFee": {
									"satoshis": 5,
									"bytes": 10
								},
								"relayFee": {
									"satoshis": 5,
									"bytes": 10
								}
							},
							"standard": {
								"miningFee": {
									"satoshis": 5,
									"bytes": 10
								},
								"relayFee": {
									"satoshis": 5,
									"bytes": 10
								}
							}
						}
					}
				`), &out)
			},
			args: dpp.PaymentRequestArgs{
				PaymentID: "bwe123",
			},
			cfg: &config.PayD{
				Host:   "securepayddest",
				Port:   ":4445",
				Secure: true,
			},
			expURL: "https://securepayddest:4445/api/v1/payments/bwe123",
			expPaymentReq: &dpp.PaymentRequest{
				SPVRequired: true,
				Network:     "mainnet",
				Memo:        "invoice 6K9oZq9",
				PaymentURL:  "https://dpp:8445/api/v1/payment/6K9oZq9",
				Destinations: dpp.PaymentDestinations{
					Outputs: []dpp.Output{{
						LockingScript: func() *bscript.Script {
							ls, err := bscript.NewFromHexString("525252")
							assert.NoError(t, err)
							return ls
						}(),
						Amount: 100,
					}, {
						LockingScript: func() *bscript.Script {
							ls, err := bscript.NewFromHexString("535353")
							assert.NoError(t, err)
							return ls
						}(),
						Amount: 400,
					}},
				},
				MerchantData: &dpp.Merchant{
					AvatarURL: "http://url.com",
					Name:      "Merchant Name",
					Email:     "merchant@demo.com",
					Address:   "123 Street Fake",
					ExtendedData: map[string]interface{}{
						"dislikes":         "trying to think up placeholder data",
						"likes":            "walks in the park at night",
						"paymentReference": "6K9oZq9",
					},
				},
				FeeRate:             bt.NewFeeQuote(),
				CreationTimestamp:   func() time.Time { t, _ := time.Parse(time.RFC3339, "2021-12-13T10:37:15.7946831Z"); return t }(),
				ExpirationTimestamp: func() time.Time { t, _ := time.Parse(time.RFC3339, "2021-12-14T10:37:15.7946831Z"); return t }(),
			},
		},
		"error is handled": {
			doFunc: func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error {
				return errors.New("yikes")
			},
			args: dpp.PaymentRequestArgs{
				PaymentID: "bwe123",
			},
			cfg: &config.PayD{
				Host:   "securepayddest",
				Port:   ":4445",
				Secure: true,
			},
			expURL: "https://securepayddest:4445/api/v1/payments/bwe123",
			expErr: errors.New("yikes"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			pd := payd.NewPayD(test.cfg, &mocks.HTTPClientMock{
				DoFunc: func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error {
					assert.Equal(t, test.expURL, url)
					return test.doFunc(ctx, method, url, statusCode, req, out)
				},
			})
			pr, err := pd.PaymentRequest(context.Background(), test.args)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expPaymentReq != nil {
				assert.NotNil(t, pr)
				assert.Equal(t, test.expPaymentReq.CreationTimestamp.String(), pr.CreationTimestamp.String())
				assert.Equal(t, test.expPaymentReq.ExpirationTimestamp.String(), pr.ExpirationTimestamp.String())

				ts := time.Now()
				pr.FeeRate.UpdateExpiry(ts)
				test.expPaymentReq.FeeRate.UpdateExpiry(ts)

				assert.Equal(t, *test.expPaymentReq, *pr)
			} else {
				assert.Nil(t, pr)
			}
		})
	}
}
