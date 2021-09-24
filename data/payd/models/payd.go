package models

import (
	"github.com/libsv/go-bc/spv"
	"github.com/libsv/go-bt/v2"

	"github.com/libsv/go-p4"
)

// PayDPaymentRequest is used to send a payment to PayD for valdiation and storage.
type PayDPaymentRequest struct {
	SPVEnvelope    *spv.Envelope
	ProofCallbacks map[string]p4.ProofCallback `json:"proofCallbacks"`
}

// Destination is a payment output with locking script.
type Destination struct {
	Script   string `json:"script"`
	Satoshis uint64 `json:"satoshis"`
}

// DestinationResponse is the response for the destinations api.
type DestinationResponse struct {
	Outputs []Destination `json:"outputs"`
	Fees    *bt.FeeQuote  `json:"fees"`
}
