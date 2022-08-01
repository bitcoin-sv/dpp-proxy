package dpp

import (
	"github.com/libsv/go-dpp/modes/hybridmode"
)

// These structures are defined in the TSC spec:
// See https://tsc.bitcoinassociation.net/standards/direct_payment_protocol

// PaymentACK message used in the TSC DPP spec.
// See https://tsc.bitcoinassociation.net/standards/direct_payment_protocol/#PaymentModes
type PaymentACK struct {
	// ModeID the chosen mode.
	ModeID string `json:"modeId" binding:"required" example:"ef63d9775da5"`
	// Mode data required by specific payment mode
	Mode        *hybridmode.PaymentACK      `json:"mode"`
	PeerChannel *hybridmode.PeerChannelData `json:"peerChannel"`
	RedirectURL string                      `json:"redirectUrl"`

	// Memo may contain information about why there was an error. This field is poorly defined until
	// error reporting is more standardised.
	Memo string
	// A number indicating why the transaction was not accepted. 0 or undefined indicates no error.
	// A 1 or any other positive integer indicates an error. The errors are left undefined for now;
	// it is recommended only to use “1” and to fill the memo with a textual explanation about why
	// the transaction was not accepted until further numbers are defined and standardised.
	Error int `json:"error,omitempty"`
}
