package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bitcoin-sv/dpp-proxy/service"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-dpp"
	dppMocks "github.com/libsv/go-dpp/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRequest_PaymentRequest(t *testing.T) {
	created := time.Now()
	expired := created.Add(time.Hour * 24)
	tests := map[string]struct {
		paymentRequestFunc func(context.Context, dpp.PaymentRequestArgs) (*dpp.PaymentTerms, error)
		args               dpp.PaymentRequestArgs
		expResp            *dpp.PaymentTerms
		expErr             error
	}{
		"successful request": {
			args: dpp.PaymentRequestArgs{
				PaymentID: "abc123",
			},
			paymentRequestFunc: func(context.Context, dpp.PaymentRequestArgs) (*dpp.PaymentTerms, error) {
				return &dpp.PaymentTerms{
					Network:             "regtest",
					Version:			 "1.0",
					CreationTimestamp:   created.Unix(),
					ExpirationTimestamp: expired.Unix(),
					Modes: &dpp.PaymentModes{
						HybridPaymentMode: map[string]map[string][]dpp.TransactionTerms{
							"choiceID0": {
								"transactions": {
									dpp.TransactionTerms{
										Outputs: dpp.Outputs{ NativeOutputs: []dpp.NativeOutput{
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
										Inputs: dpp.Inputs{},
										Policies: &dpp.Policies{},
									},
								},
							},

						},
					},
					PaymentURL: "http://iamsotest/api/v1/payment/abc123",
					Memo:       "invoice abc123",
					Beneficiary: &dpp.Beneficiary{
						ExtendedData: map[string]interface{}{"paymentReference": "abc123"},
					},
				}, nil
			},
			expResp: &dpp.PaymentTerms{
				Network:             "regtest",
				Version:			 "1.0",
				CreationTimestamp:   created.Unix(),
				ExpirationTimestamp: expired.Unix(),
				Modes: &dpp.PaymentModes{
					HybridPaymentMode: map[string]map[string][]dpp.TransactionTerms{
						"choiceID0": {
							"transactions": {
								dpp.TransactionTerms{
									Outputs: dpp.Outputs{ NativeOutputs: []dpp.NativeOutput{
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
									Inputs: dpp.Inputs{},
									Policies: &dpp.Policies{},
								},
							},
						},

					},
				},
				PaymentURL: "http://iamsotest/api/v1/payment/abc123",
				Memo:       "invoice abc123",
				Beneficiary: &dpp.Beneficiary{
					ExtendedData: map[string]interface{}{"paymentReference": "abc123"},
				},
			},
		},
		"successful request with nil extended data": {
			args: dpp.PaymentRequestArgs{
				PaymentID: "abc123",
			},
			paymentRequestFunc: func(context.Context, dpp.PaymentRequestArgs) (*dpp.PaymentTerms, error) {
				return &dpp.PaymentTerms{
					Network:             "regtest",
					Version:			 "1.0",
					CreationTimestamp:   created.Unix(),
					ExpirationTimestamp: expired.Unix(),
					Modes: &dpp.PaymentModes{
						HybridPaymentMode: map[string]map[string][]dpp.TransactionTerms{
							"choiceID0": {
								"transactions": {
									dpp.TransactionTerms{
										Outputs: dpp.Outputs{ NativeOutputs: []dpp.NativeOutput{
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
										Inputs: dpp.Inputs{},
										Policies: &dpp.Policies{},
									},
								},
							},

						},
					},
					Beneficiary: &dpp.Beneficiary{},
					PaymentURL:   "http://iamsotest/api/v1/payment/abc123",
					Memo:         "invoice abc123",
				}, nil
			},
			expResp: &dpp.PaymentTerms{
				Network:             "regtest",
				Version:			 "1.0",
				CreationTimestamp:   created.Unix(),
				ExpirationTimestamp: expired.Unix(),
				Modes: &dpp.PaymentModes{
					HybridPaymentMode: map[string]map[string][]dpp.TransactionTerms{
						"choiceID0": {
							"transactions": {
								dpp.TransactionTerms{
									Outputs: dpp.Outputs{ NativeOutputs: []dpp.NativeOutput{
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
									Inputs: dpp.Inputs{},
									Policies: &dpp.Policies{},
								},
							},
						},

					},
				},
				PaymentURL: "http://iamsotest/api/v1/payment/abc123",
				Memo:       "invoice abc123",
				Beneficiary: &dpp.Beneficiary{
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
			paymentRequestFunc: func(context.Context, dpp.PaymentRequestArgs) (*dpp.PaymentTerms, error) {
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
