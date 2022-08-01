package mocks

//go:generate moq -pkg mocks -out payment_writer.go ../ PaymentWriter
//go:generate moq -pkg mocks -out payment_service.go ../ PaymentService
//go:generate moq -pkg mocks -out payment_request_service.go ../ PaymentTermsService
