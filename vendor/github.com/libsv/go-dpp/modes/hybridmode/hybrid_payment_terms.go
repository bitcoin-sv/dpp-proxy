package hybridmode

import (
	"github.com/libsv/go-dpp/nativetypes"
)

// These structures are defined in the TSC spec:
// See https://tsc.bitcoinassociation.net/standards/direct_payment_protocol


// Policies An object containing some policy information like fees or whether Ancestors are
// required in the `Payment`.
type Policies struct {
	// FeeRate defines the amount of fees a users wallet should add to the payment
	// when submitting their final payments.
	FeeRate     map[string]map[string]int `json:"fees,omitempty"`
	SPVRequired bool                      `json:"SPVRequired,omitempty"`
	LockTime    uint32                    `json:"lockTime,omitempty"`
}

// Inputs provides options of different arrays of input script types.
// Currently, only "native" type input are supported.
type Inputs struct {
	NativeOutputs []nativetypes.NativeInput `json:"native"`
}

// Outputs provides options of different arrays of output script types.
// Currently, only "native" type outputs are supported.
type Outputs struct {
	NativeOutputs []nativetypes.NativeOutput `json:"native"`
}

// TransactionTerms a single definition of requested transaction format for the standard payment mode:
// "ef63d9775da5" in the DPP TSC spec.
type TransactionTerms struct {
	Outputs  Outputs   `json:"outputs"`
	Inputs   Inputs    `json:"inputs,omitempty"`
	Policies *Policies `json:"policies"`
}

// PaymentTerms message used in DPP TSC spec. for the `PaymentTerms` message.
type PaymentTerms map[string]map[string][]TransactionTerms
