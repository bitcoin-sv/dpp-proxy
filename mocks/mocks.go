package mocks

//go:generate moq -pkg mocks -out payment_writer.go ../ PaymentWriter
//go:generate moq -pkg mocks -out merchant_reader.go ../ MerchantReader
//go:generate moq -pkg mocks -out destination_reader.go ../ DestinationReader
