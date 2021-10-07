package p4

import (
	"context"
	"time"

	"github.com/libsv/go-bt/v2"
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

// PaymentDestinations contains the supported destinations
// by this P4 server.
type PaymentDestinations struct {
	Outputs []Output `json:"outputs"`
}

// Destinations message containing outputs and their fees.
type Destinations struct {
	SPVRequired bool
	Outputs     []Output
	Fees        *bt.FeeQuote
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

// DestinationReader interfaces retrieving payment destinations.
type DestinationReader interface {
	Destinations(ctx context.Context, args PaymentRequestArgs) (*Destinations, error)
}
