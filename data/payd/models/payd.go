package models

import (
	"time"

	"github.com/libsv/go-bt/v2"

	"github.com/libsv/go-dpp"
)

// PayDPaymentRequest is used to send a payment to PayD for valdiation and storage.
type PayDPaymentRequest struct {
	Ancestry       *string                      `json:"ancestry"`
	RawTx          *string                      `json:"rawTx"`
	ProofCallbacks map[string]dpp.ProofCallback `json:"proofCallbacks"`
}

// Destination is a payment output with locking script.
type Destination struct {
	Script   string `json:"script"`
	Satoshis uint64 `json:"satoshis"`
}

// DestinationResponse is the response for the destinations api.
type DestinationResponse struct {
	AncestryRequired bool          `json:"ancestryRequired"`
	Network          string        `json:"network"`
	Outputs          []Destination `json:"outputs"`
	Fees             *bt.FeeQuote  `json:"fees"`
	CreatedAt        time.Time     `json:"createdAt"`
	ExpiresAt        time.Time     `json:"expiresAt"`
}
