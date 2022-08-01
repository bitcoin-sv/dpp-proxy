package hybridmode


// PeerChannelData holds peer channel information for subscribing to and reading from a peer channel.
type PeerChannelData struct {
	Host      string `json:"host"`
	Path      string `json:"path"`
	ChannelID string `json:"channel_id"`
	Token     string `json:"token"`
}


// PaymentACK includes data required for hybridmode payment mode.
type PaymentACK struct {
	TransactionIds []string         `json:"transactionIds"`
	PeerChannel    *PeerChannelData `json:"peerChannel"`
}
