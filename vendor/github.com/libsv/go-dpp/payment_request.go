package dpp

import (
	"context"
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

// NativeInput a way of declaring requirements for the inputs which should be used.
type NativeInput struct {
	ScriptSig string `json:"scriptSig" binding:"required"` // string. required.
	TxID      string `json:"txid" binding:"required"`      // string. required.
	Vout      uint32 `json:"vout" binding:"required"`      // integer. required.
	Value     uint64 `json:"value" binding:"required"`     // integer. required.
	NSequence int    `json:"nSequence,omitempty"`          // number. optional.
}

// Inputs provides options of different arrays of input script types.
// Currently, only "native" type input are supported.
type Inputs struct {
	NativeOutputs []NativeInput `json:"native"`
}

// Outputs provides options of different arrays of output script types.
// Currently, only "native" type outputs are supported.
type Outputs struct {
	NativeOutputs []NativeOutput `json:"native"`
}

// TransactionTerms a single definition of requested transaction format for the standard payment mode:
// "ef63d9775da5" in the DPP TSC spec.
type TransactionTerms struct {
	Outputs  Outputs   `json:"outputs"`
	Inputs   Inputs     `json:"inputs,omitempty"`
	Policies *Policies `json:"policies"`
}

// PaymentModes message used in DPP TSC spec.
// At present we will strictly only allow the "standard" mode of payment with native bitcoins (satoshis). Handling
// of tokens is left for a later date.
type PaymentModes struct {
	HybridPaymentMode map[string]map[string][]TransactionTerms `json:"ef63d9775da5"`
}

// PaymentRequest message as defined in the DPP T$C spec.
type PaymentRequest struct {
	// Network  Always set to "bitcoin" (but seems to be set to 'bitcoin-sv'
	// outside bip270 spec, see https://handcash.github.io/handcash-merchant-integration/#/merchant-payments)
	// {enum: bitcoin, bitcoin-sv, test}
	// Required.
	Network string `json:"network" binding:"required" example:"mainnet" enums:"mainnet,testnet,stn,regtest"`
	// Version version of DPP TSC spec.
	// Required.
	Version string `json:"version" binding:"required" example:"1.0"`
	// Outputs an array of outputs. DEPRECATED but included for backward compatibility.
	// Optional.
	Outputs []NativeOutput `json:"outputs,omitempty"`
	// CreationTimestamp Unix timestamp (seconds since 1-Jan-1970 UTC) when the PaymentRequest was created.
	// Required.
	CreationTimestamp int64 `json:"creationTimestamp" binding:"required" swaggertype:"primitive,int" example:"1648163657"`
	// ExpirationTimestamp Unix timestamp (UTC) after which the PaymentRequest should be considered invalid.
	// Optional.
	ExpirationTimestamp int64 `json:"expirationTimestamp" binding:"required" swaggertype:"primitive,int" example:"1648164657"`
	// PaymentURL secure HTTPS location where a Payment message (see below) will be sent to obtain a PaymentACK.
	// Maximum length is 4000 characters
	PaymentURL string `json:"paymentUrl" binding:"required" example:"http://localhost:3443/api/v1/payment/123456"`
	// Memo note that should be displayed to the customer, explaining what this PaymentRequest is for.
	// Maximum length is 50 characters.
	// Optional.
	Memo string `json:"memo,omitempty" example:"invoice number 123456"`
	// Beneficiary Arbitrary data that may be used by the payment host to identify the PaymentRequest
	// May be omitted if the payment host does not need to associate Payments with PaymentRequest
	// or if they associate each PaymentRequest with a separate payment address.
	// Maximum length is 10000 characters.
	// Optional.
	Beneficiary *Merchant `json:"beneficiary,omitempty"`
	// Modes TSC payment modes specified by ID (and well defined) modes customer can choose to pay
	// A key-value map. required field but not if legacy BIP270 outputs are provided
	Modes *PaymentModes `json:"modes"`
}

// PaymentRequestArgs are request arguments that can be passed to the service.
type PaymentRequestArgs struct {
	// PaymentID is an identifier for an invoice.
	PaymentID string `param:"paymentID"`
}

// PaymentRequestService can be implemented to enforce business rules
// and process in order to fulfil a PaymentRequest.
type PaymentRequestService interface {
	PaymentRequestReader
}

// PaymentRequestReader will return a new payment request.
type PaymentRequestReader interface {
	PaymentRequest(ctx context.Context, args PaymentRequestArgs) (*PaymentRequest, error)
}
