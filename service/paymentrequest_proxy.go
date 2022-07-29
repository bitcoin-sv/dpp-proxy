package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/dpp-proxy/config"
	"github.com/libsv/go-dpp"
	"github.com/pkg/errors"
	validator "github.com/theflyingcodr/govalidator"
)

// paymentRequestProxy simply acts as a pass-through to the data layer
// where another service will create the paymentRequest.
// TODO - remove the other payment request service.
type paymentRequestProxy struct {
	preqRdr   dpp.PaymentRequestReader
	transCfg  *config.Transports
	walletCfg *config.Server
}

// NewPaymentRequestProxy will setup and return a new PaymentRequest service that will generate outputs
// using the provided outputter which is defined in server config.
func NewPaymentRequestProxy(preqRdr dpp.PaymentRequestReader, transCfg *config.Transports, walletCfg *config.Server) *paymentRequestProxy {
	return &paymentRequestProxy{
		preqRdr:   preqRdr,
		transCfg:  transCfg,
		walletCfg: walletCfg,
	}
}

// PaymentRequest will call to the data layer to return a full payment request.
func (p *paymentRequestProxy) PaymentRequest(ctx context.Context, args dpp.PaymentRequestArgs) (*dpp.PaymentTerms, error) {
	if err := validator.New().
		Validate("paymentID", validator.NotEmpty(args.PaymentID)); err.Err() != nil {
		return nil, err
	}
	resp, err := p.preqRdr.PaymentRequest(ctx, args)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read payment request for paymentID %s", args.PaymentID)
	}

	if len(resp.Modes.HybridPaymentMode["choiceID0"]["transactions"][0].Outputs.NativeOutputs) == 0 {
		return nil, fmt.Errorf("no outputs received for paymentID %s", args.PaymentID)
	}

	if resp.Modes.HybridPaymentMode["choiceID0"]["transactions"][0].Policies.FeeRate == nil {
		return nil, fmt.Errorf("no fees received for paymentID %s", args.PaymentID)
	}

	if p.transCfg.Mode == config.TransportModeHybrid {
		u := url.URL{
			Scheme: "http",
			Host:   p.walletCfg.FQDN,
			Path:   "/api/v1/payment/" + args.PaymentID,
		}
		resp.PaymentURL = u.String()
	}

	return resp, nil
}
