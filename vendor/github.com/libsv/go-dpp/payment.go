package dpp

import (
	"github.com/libsv/go-bc/spv"
	validator "github.com/theflyingcodr/govalidator"
)

// These structures are defined in the TSC spec:
// See https://tsc.bitcoinassociation.net/standards/direct_payment_protocol

// HybridPaymentModePayment includes data required for hybrid payment mode.
type HybridPaymentModePayment struct {
	// OptionID ID of chosen payment options
	OptionID string `json:"optionId"`
	// Transactions A list of valid, signed Bitcoin transactions that fully pays the PaymentTerms.
	// The transaction is hex-encoded and must NOT be prefixed with “0x”.
	// The order of transactions should match the order from PaymentTerms for this mode.
	Transactions []string `json:"transactions"`
	// Ancestors a map of txid to ancestry transaction info for the transactions in <optionID> above
	// each ancestor contains the TX together with the MerkleProof needed when SPVRequired is true.
	// See: https://tsc.bitcoinassociation.net/standards/transaction-ancestors/
	Ancestors map[string]spv.TSCAncestryJSON `json:"ancestors"`
}

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
	// ModeID chosen from possible modes of PaymentTerms.
	ModeID string `json:"modeId" binding:"required" example:"ef63d9775da5"`
	// Mode Object with data required by specific mode, e.g. HybridPaymentMode
	Mode HybridPaymentModePayment `json:"mode" binding:"required"`
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
