package mocks

//go:generate moq -pkg mocks -out payment_writer.go ../ PaymentWriter
//go:generate moq -pkg mocks -out fee_reader.go ../ FeeReader
//go:generate moq -pkg mocks -out merchant_reader.go ../ MerchantReader
