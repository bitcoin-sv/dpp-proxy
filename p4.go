package p4

import validator "github.com/theflyingcodr/govalidator"

// ClientError defines an error type that can be returned to handle client errors.
type ClientError struct {
	ID      string `json:"id" example:"e97970bf-2a88-4bc8-90e6-2f597a80b93d"`
	Code    string `json:"code" example:"N01"`
	Title   string `json:"title" example:"not found"`
	Message string `json:"message" example:"unable to find foo when loading bar"`
}

// BadRequestError defines an error type to handle validation errors.
type BadRequestError struct {
	Errors validator.ErrValidation `json:"errors"`
}