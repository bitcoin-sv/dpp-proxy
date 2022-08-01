package dpp

import (
	"context"
	"github.com/libsv/go-dpp/modes/hybridmode"
	validator "github.com/theflyingcodr/govalidator"
)

// These structures are defined in the TSC spec:
// See https://tsc.bitcoinassociation.net/standards/direct_payment_protocol

// Originator Data about payer. This data might be needed in many cases, e.g. tracking data for later loyalty
// points processing etc.
type Originator struct {
	// Name name of payer.
	Name string `json:"name"`
	// Paymail Payer’s paymail (where e.g. refunds will be send, identity can be use somehow etc.).
	Paymail string `json:"paymail"`
	// Avatar URL to an avatar.
	Avatar string `json:"avatar"`
	// ExtendedData additional optional data.
	ExtendedData map[string]interface{} `json:"extendedData"`
}

// Payment is a Payment message used in BIP270.
// See https://github.com/moneybutton/bips/blob/master/bip-0270.mediawiki#payment
type Payment struct {
	// ModeID chosen from possible messages of PaymentTerms.
	ModeID string `json:"modeId" binding:"required" example:"ef63d9775da5"`
	// Mode Object with data required by specific mode, e.g. HybridPaymentTerms
	Mode hybridmode.Payment `json:"mode" binding:"required"`
	// Originator Data about payer. This data might be needed in many cases, e.g. refund, tract data for later loyalty points processing etc.
	Originator Originator `json:"originator"`
	// Transaction A single valid, signed Bitcoin transaction that fully pays the PaymentTerms. This field is deprecated.
	Transaction *string `json:"transaction,omitempty"`
	// Memo A plain-text note from the customer to the payment host.
	Memo string `json:"memo,omitempty"`
}

// Validate will ensure the users request is correct.
func (p Payment) Validate() error {
	v := validator.New().
		Validate("modeId", validator.NotEmpty(p.ModeID)).
		Validate("mode", validator.NotEmpty(p.Mode)).
		Validate("mode.optionId", validator.NotEmpty(p.Mode.OptionID)).
		Validate("mode.transactions", validator.NotEmpty(p.Mode.Transactions))
	return v.Err()
}

// ProofCallback is used by a payee to request a merkle proof is sent to them
// as proof of acceptance of the tx they have provided in the ancestry.
type ProofCallback struct {
	Token string `json:"token"`
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
