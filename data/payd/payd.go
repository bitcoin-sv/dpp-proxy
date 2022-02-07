package payd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/libsv/go-bk/envelope"
	"github.com/pkg/errors"

	"github.com/libsv/go-p4"
	"github.com/libsv/p4-server/config"
	"github.com/libsv/p4-server/data"
	"github.com/libsv/p4-server/data/payd/models"
)

// Known endpoints for the payd wallet implementing the payment protocol interface.
const (
	urlPayments      = "%s/api/v1/payments/%s"
	urlProofs        = "%s/api/v1/proofs/%s"
	protocolInsecure = "http"
	protocolSecure   = "https"
)

type payd struct {
	client data.HTTPClient
	cfg    *config.PayD
}

// NewPayD will setup a new store that can interface with a payd wallet implementing
// the Payment Protocol Interface.
func NewPayD(cfg *config.PayD, client data.HTTPClient) *payd {
	return &payd{
		cfg:    cfg,
		client: client,
	}
}

// PaymentRequest will fetch a payment request message from payd for a given payment.
func (p *payd) PaymentRequest(ctx context.Context, args p4.PaymentRequestArgs) (*p4.PaymentRequest, error) {
	var resp p4.PaymentRequest
	if err := p.client.Do(ctx, http.MethodGet, fmt.Sprintf(urlPayments, p.baseURL(), args.PaymentID), http.StatusOK, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PaymentCreate will post a request to payd to validate and add the txos to the wallet.
//
// If invalid a non 204 status code is returned.
func (p *payd) PaymentCreate(ctx context.Context, args p4.PaymentCreateArgs, req p4.Payment) (*p4.PaymentACK, error) {
	paymentReq := models.PayDPaymentRequest{
		RawTX:          req.RawTX,
		SPVEnvelope:    req.SPVEnvelope,
		ProofCallbacks: req.ProofCallbacks,
	}
	var ack p4.PaymentACK
	if err := p.client.Do(ctx, http.MethodPost, fmt.Sprintf(urlPayments, p.baseURL(), args.PaymentID), http.StatusNoContent, paymentReq, &ack); err != nil {
		return nil, err
	}
	return &ack, nil
}

// ProofCreate will pass on the proof to a payd instance for storage.
func (p *payd) ProofCreate(ctx context.Context, args p4.ProofCreateArgs, req envelope.JSONEnvelope) error {
	return errors.WithStack(p.client.Do(ctx, http.MethodPost, fmt.Sprintf(urlProofs, p.baseURL(), args.TxID), http.StatusCreated, req, nil))
}

// baseURL will return http or https depending on if we're using TLS.
func (p *payd) baseURL() string {
	if p.cfg.Secure {
		return fmt.Sprintf("%s://%s%s", protocolSecure, p.cfg.Host, p.cfg.Port)
	}
	return fmt.Sprintf("%s://%s%s", protocolInsecure, p.cfg.Host, p.cfg.Port)
}
