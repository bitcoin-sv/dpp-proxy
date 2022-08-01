package nativetypes

import (
	"github.com/libsv/go-bt/v2/bscript"
)

// NativeOutput defines a native type output as opposed to a token for example.
type NativeOutput struct {
	// Amount is the number of satoshis to be paid.
	Amount uint64 `json:"amount" binding:"required" example:"100000"`
	// Script is a locking script where payment should be sent, formatted as a hexadecimal string.
	LockingScript *bscript.Script `json:"script" binding:"required" swaggertype:"primitive,string" example:"76a91455b61be43392125d127f1780fb038437cd67ef9c88ac"`
	// Description, an optional description such as "tip" or "sales tax". Maximum length is 100 chars.
	Description string `json:"description,omitempty" example:"paymentReference 123456"`
}


// NativeInput a way of declaring requirements for the inputs which should be used. It is "native" to distinguish it
// from a token input in the hybridmode payment mode.
type NativeInput struct {
	ScriptSig string `json:"scriptSig" binding:"required"` // string. required.
	TxID      string `json:"txid" binding:"required"`      // string. required.
	Vout      uint32 `json:"vout" binding:"required"`      // integer. required.
	Value     uint64 `json:"value" binding:"required"`     // integer. required.
	NSequence int    `json:"nSequence,omitempty"`          // number. optional.
}
