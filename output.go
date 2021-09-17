package pptcl

import (
	"context"
)

// Output message used in BIP270.
// See https://github.com/moneybutton/bips/blob/master/bip-0270.mediawiki#output
type Output struct {
	// Amount is the number of satoshis to be paid.
	Amount uint64 `json:"amount" example:"100000"`
	// Script is a locking script where payment should be sent, formatted as a hexadecimal string.
	Script string `json:"script" example:"76a91455b61be43392125d127f1780fb038437cd67ef9c88ac"`
	// Description, an optional description such as "tip" or "sales tax". Maximum length is 100 chars.
	Description string `json:"description" example:"paymentReference 123456"`
}

// OutputReader will read outputs from another system.
type OutputReader interface {
	// Outputs will be used to get locking scripts from an underlying service or data store.
	//
	// This use case is that the underlying system has the invoice and therefor the amount.
	// We then send the invoice / paymentID for lookup and the service will create
	// n outputs to equal the amount of satoshis. It may also add additional outputs
	// for merchant fees or tax etc.
	Outputs(ctx context.Context, args PaymentRequestArgs) ([]Output, error)
}
