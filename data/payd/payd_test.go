package payd_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/libsv/go-bc/spv"
	"github.com/libsv/go-dpp/modes/hybridmode"
	"github.com/libsv/go-dpp/nativetypes"
	"testing"
	"time"

	"github.com/bitcoin-sv/dpp-proxy/config"
	"github.com/bitcoin-sv/dpp-proxy/data/payd"
	"github.com/bitcoin-sv/dpp-proxy/data/payd/models"
	"github.com/bitcoin-sv/dpp-proxy/mocks"
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
		expReq models.PayDPayment
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
				ModeID: "ef63d9775da5",
				Mode: hybridmode.Payment{
					OptionID:     "choiceID0",
					Transactions: []string{"tx1 hex", "tx2 hex"},
					Ancestors:    map[string]spv.TSCAncestryJSON{},
				},
			},
			expReq: models.PayDPayment{
				ModeID: "ef63d9775da5",
				Mode: hybridmode.Payment{
					OptionID:     "choiceID0",
					Transactions: []string{"tx1 hex", "tx2 hex"},
					Ancestors:    map[string]spv.TSCAncestryJSON{},
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
				ModeID: "ef63d9775da5",
				Mode: hybridmode.Payment{
					OptionID:     "choiceID0",
					Transactions: []string{"tx1 hex", "tx2 hex"},
					Ancestors:    map[string]spv.TSCAncestryJSON{},
				},
			},
			expReq: models.PayDPayment{
				ModeID: "ef63d9775da5",
				Mode: hybridmode.Payment{
					OptionID:     "choiceID0",
					Transactions: []string{"tx1 hex", "tx2 hex"},
					Ancestors:    map[string]spv.TSCAncestryJSON{},
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
				ModeID: "ef63d9775da5",
				Mode: hybridmode.Payment{
					OptionID:     "choiceID0",
					Transactions: []string{"tx1 hex", "tx2 hex"},
					Ancestors:    map[string]spv.TSCAncestryJSON{},
				},
				Originator: dpp.Originator{
					Name: 		"Bob the builder",
					Paymail: 	"bob@bestpaymail.com",
					Avatar:  	"https://iamges.com/vtwe4eerf",
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
				},
			},
			expReq: models.PayDPayment{
				ModeID: "ef63d9775da5",
				Mode: hybridmode.Payment{
					OptionID:     "choiceID0",
					Transactions: []string{"tx1 hex", "tx2 hex"},
					Ancestors:    map[string]spv.TSCAncestryJSON{},
				},
				Originator: dpp.Originator{
					Name: 		"Bob the builder",
					Paymail: 	"bob@bestpaymail.com",
					Avatar:  	"https://iamges.com/vtwe4eerf",
					ExtendedData: map[string]interface{}{"paymentReference": "omgwow"},
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

func TestPayd_PaymentTerms(t *testing.T) {
	tests := map[string]struct {
		doFunc        func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error
		args          dpp.PaymentTermsArgs
		cfg           *config.PayD
		expURL        string
		expPaymentReq *dpp.PaymentTerms
		expErr        error
	}{
		"successful payment request": {
			doFunc: func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error {
				return json.Unmarshal([]byte(`
					{
						"network": "mainnet",
						"version": "1.0",
						"ancestryRequired": true,
						"modes": {"ef63d9775da5": {
							"choiceID0": {
								"transactions": [
									{
										"outputs": {
											"native": [
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
										"policies": {
											"fees":
												{"standard": {"satoshis": 100, "bytes": 200},
													"data": {"satoshis": 100, "bytes": 200}},
											"SPVRequired": false
										}
									}
								]
							}
						}},
						"creationTimestamp": 1639391835,
						"expirationTimestamp": 1639478235,
						"paymentUrl": "http://dpp:8445/api/v1/payment/6K9oZq9",
						"memo": "invoice 6K9oZq9",
						"beneficiary": {
							"avatar": "http://url.com",
							"name": "Merchant Name",
							"email": "merchant@demo.com",
							"address": "123 Street Fake",
							"extendedData": {
								"dislikes": "trying to think up placeholder data",
								"likes": "walks in the park at night",
								"paymentReference": "6K9oZq9"
							}
						}
					}
				`), &out)
			},
			args: dpp.PaymentTermsArgs{
				PaymentID: "qwe123",
			},
			cfg: &config.PayD{
				Host: "payddest",
				Port: ":445",
			},
			expURL: "http://payddest:445/api/v1/payments/qwe123",
			expPaymentReq: &dpp.PaymentTerms{
				Network:          "mainnet",
				Version:		  "1.0",
				Memo:             "invoice 6K9oZq9",
				PaymentURL:       "http://dpp:8445/api/v1/payment/6K9oZq9",
				Modes: &dpp.PaymentTermsModes{
					Hybrid: hybridmode.PaymentTerms{
						"choiceID0": {
							"transactions": {
								hybridmode.TransactionTerms{
									Outputs: hybridmode.Outputs{ NativeOutputs: []nativetypes.NativeOutput{
										{
											Amount:        100,
											LockingScript: func() *bscript.Script {
												ls, _ := bscript.NewFromHexString("525252")
												return ls
											}(),
										},
										{
											Amount:        400,
											LockingScript: func() *bscript.Script {
												ls, _ := bscript.NewFromHexString("535353")
												return ls
											}(),
										},
									} },
									Inputs: hybridmode.Inputs{},
									Policies: &hybridmode.Policies{
										FeeRate: map[string]map[string]int{
											"data":
												{"bytes":200,"satoshis":100},
												"standard":
												{"bytes":200,"satoshis":100},
										},
										SPVRequired: false,
										LockTime:    0,
									},
								},
							},
						},

					},
				},
				Beneficiary: &dpp.Beneficiary{
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
				CreationTimestamp:   func() int64 { t, _ := time.Parse(time.RFC3339, "2021-12-13T10:37:15.7946831Z"); return t.Unix() }(),
				ExpirationTimestamp: func() int64 { t, _ := time.Parse(time.RFC3339, "2021-12-14T10:37:15.7946831Z"); return t.Unix() }(),
			},
		},
		"successful https payment request": {
			doFunc: func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error {
				return json.Unmarshal([]byte(`
					{
						"network": "mainnet",
						"version": "1.0",
						"ancestryRequired": true,
						"modes": {"ef63d9775da5": {
							"choiceID0": {
								"transactions": [
									{
										"outputs": {
											"native": [
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
										"policies": {
											"fees":
												{"standard": {"satoshis": 100, "bytes": 200},
													"data": {"satoshis": 100, "bytes": 200}},
											"SPVRequired": false
										}
									}
								]
							}
						}},
						"creationTimestamp": 1639391835,
						"expirationTimestamp": 1639478235,
						"paymentUrl": "https://dpp:8445/api/v1/payment/6K9oZq9",
						"memo": "invoice 6K9oZq9",
						"beneficiary": {
							"avatar": "http://url.com",
							"name": "Merchant Name",
							"email": "merchant@demo.com",
							"address": "123 Street Fake",
							"extendedData": {
								"dislikes": "trying to think up placeholder data",
								"likes": "walks in the park at night",
								"paymentReference": "6K9oZq9"
							}
						}
					}
				`), &out)
			},
			args: dpp.PaymentTermsArgs{
				PaymentID: "bwe123",
			},
			cfg: &config.PayD{
				Host:   "securepayddest",
				Port:   ":4445",
				Secure: true,
			},
			expURL: "https://securepayddest:4445/api/v1/payments/bwe123",
			expPaymentReq: &dpp.PaymentTerms{
				Network:          "mainnet",
				Version:		  "1.0",
				Memo:             "invoice 6K9oZq9",
				PaymentURL:       "https://dpp:8445/api/v1/payment/6K9oZq9",
				Modes: &dpp.PaymentTermsModes{
					Hybrid: hybridmode.PaymentTerms{
						"choiceID0": {
							"transactions": {
								hybridmode.TransactionTerms{
									Outputs: hybridmode.Outputs{ NativeOutputs: []nativetypes.NativeOutput{
										{
											Amount:        100,
											LockingScript: func() *bscript.Script {
												ls, _ := bscript.NewFromHexString("525252")
												return ls
											}(),
										},
										{
											Amount:        400,
											LockingScript: func() *bscript.Script {
												ls, _ := bscript.NewFromHexString("535353")
												return ls
											}(),
										},
									} },
									Inputs: hybridmode.Inputs{},
									Policies: &hybridmode.Policies{
										FeeRate: map[string]map[string]int{
											"data":
											{"bytes":200,"satoshis":100},
											"standard":
											{"bytes":200,"satoshis":100},
										},
										SPVRequired: false,
										LockTime:    0,
									},
								},
							},
						},

					},
				},
				Beneficiary: &dpp.Beneficiary{
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
				CreationTimestamp:   func() int64 { t, _ := time.Parse(time.RFC3339, "2021-12-13T10:37:15.7946831Z"); return t.Unix() }(),
				ExpirationTimestamp: func() int64 { t, _ := time.Parse(time.RFC3339, "2021-12-14T10:37:15.7946831Z"); return t.Unix() }(),
			},
		},
		"error is handled": {
			doFunc: func(ctx context.Context, method string, url string, statusCode int, req, out interface{}) error {
				return errors.New("yikes")
			},
			args: dpp.PaymentTermsArgs{
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
			pr, err := pd.PaymentTerms(context.Background(), test.args)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expPaymentReq != nil {
				assert.NotNil(t, pr)
				assert.Equal(t, test.expPaymentReq.CreationTimestamp, pr.CreationTimestamp)
				assert.Equal(t, test.expPaymentReq.ExpirationTimestamp, pr.ExpirationTimestamp)

				assert.Equal(t, *test.expPaymentReq, *pr)
			} else {
				assert.Nil(t, pr)
			}
		})
	}
}
