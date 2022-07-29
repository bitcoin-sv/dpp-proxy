package dpp

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
