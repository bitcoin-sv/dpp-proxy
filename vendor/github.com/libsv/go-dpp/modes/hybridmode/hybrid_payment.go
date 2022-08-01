package hybridmode

import "github.com/libsv/go-bc/spv"


// Payment includes data required for hybridmode payment mode.
type Payment struct {
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
