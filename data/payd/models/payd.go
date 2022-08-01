package models

import (
	"github.com/libsv/go-dpp"
	"github.com/libsv/go-dpp/modes/hybridmode"
	"time"

	"github.com/libsv/go-bt/v2"
)

// PayDPayment is used to send a payment to PayD for valdiation and storage.
type PayDPayment struct {
	// ModeID chosen from possible modes of PaymentTerms.
	ModeID string `json:"modeId" binding:"required" example:"ef63d9775da5"`
	// Mode Object with data required by specific mode, e.g. HybridPaymentMode
	Mode hybridmode.Payment `json:"mode" binding:"required"`
	// Originator Data about payer. This data might be needed in many cases, e.g. refund, tract data for later loyalty points processing etc.
	Originator dpp.Originator `json:"originator"`
	// Transaction A single valid, signed Bitcoin transaction that fully pays the PaymentTerms. This field is deprecated.
	Transaction *string `json:"transaction,omitempty"`
	// Memo A plain-text note from the customer to the payment host.
	Memo string `json:"memo,omitempty"`
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
