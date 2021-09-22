package models

import (
	"github.com/libsv/go-bc/spv"
	"github.com/libsv/go-bt/v2"

	"github.com/libsv/pptcl"
)

// PayDPaymentRequest is used to send a payment to PayD for valdiation and storage.
type PayDPaymentRequest struct {
	SPVEnvelope    *spv.Envelope
	ProofCallbacks map[string]pptcl.ProofCallback `json:"proofCallbacks"`
}

// Destination is a payment output with locking script.
type Destination struct {
	Script   string `json:"script"`
	Satoshis uint64 `json:"satoshis"`
}

type Fees struct {
	Standard *bt.Fee `json:"standard"`
	Data     *bt.Fee `json:"data"`
}

// DestinationResponse is the response for the destinations api.
type DestinationResponse struct {
	Outputs []Destination `json:"outputs"`
	Fees    Fees          `json:"fees"`
}
