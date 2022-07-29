package dpp

import (
	"context"
	validator "github.com/theflyingcodr/govalidator"
)

// HybridPaymentModeACK includes data required for hybrid payment mode.
type HybridPaymentModeACK struct {
	TransactionIds []string         `json:"transactionIds"`
	PeerChannel    *PeerChannelData `json:"peerChannel"`
}

// PaymentACK message used in the TSC DPP spec.
// See https://tsc.bitcoinassociation.net/standards/direct_payment_protocol/#PaymentModes
type PaymentACK struct {
	// ModeID the chosen mode.
	ModeID string `json:"modeId" binding:"required" example:"ef63d9775da5"`
	// Mode data required by specific payment mode
	Mode        *HybridPaymentModeACK `json:"mode"`
	PeerChannel *PeerChannelData      `json:"peerChannel"`
	RedirectURL string                `json:"redirectUrl"`

	// Memo may contain information about why there was an error. This field is poorly defined until
	// error reporting is more standardised.
	Memo string
	// A number indicating why the transaction was not accepted. 0 or undefined indicates no error.
	// A 1 or any other positive integer indicates an error. The errors are left undefined for now;
	// it is recommended only to use “1” and to fill the memo with a textual explanation about why
	// the transaction was not accepted until further numbers are defined and standardised.
	Error int `json:"error,omitempty"`
}

// PeerChannelData holds peer channel information for subscribing to and reading from a peer channel.
type PeerChannelData struct {
	Host      string `json:"host"`
	Path      string `json:"path"`
	ChannelID string `json:"channel_id"`
	Token     string `json:"token"`
}

// PaymentCreateArgs identifies the paymentID used for the payment.
type PaymentCreateArgs struct {
	PaymentID string `param:"paymentID"`
}

// Validate will ensure that the PaymentCreateArgs are supplied and correct.
func (p PaymentCreateArgs) Validate() error {
	return validator.New().
		Validate("paymentID", validator.NotEmpty(p.PaymentID)).
		Err()
}

// PaymentService enforces business rules when creating payments.
type PaymentService interface {
	PaymentCreate(ctx context.Context, args PaymentCreateArgs, req Payment) (*PaymentACK, error)
}

// PaymentWriter will write a payment to a data store.
type PaymentWriter interface {
	PaymentCreate(ctx context.Context, args PaymentCreateArgs, req Payment) (*PaymentACK, error)
}
