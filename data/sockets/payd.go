package sockets

import (
	"context"

	"github.com/libsv/go-bk/envelope"
	"github.com/theflyingcodr/sockets"

	"github.com/libsv/go-p4"
)

type payd struct {
	s sockets.ServerChannelBroadcaster
}

// NewPayd will setup and return a new payd socket data store.
func NewPayd(b sockets.ServerChannelBroadcaster) *payd {
	return &payd{s: b}
}

// ProofCreate will broadcast the proof to all currently listening clients on the socket channel.
func (p *payd) ProofCreate(ctx context.Context, args p4.ProofCreateArgs, req envelope.JSONEnvelope) error {
	msg := sockets.NewMessage("proof.create", "", args.PaymentReference)
	msg.AppID = "p4"
	msg.CorrelationID = args.TxID
	if err := msg.WithBody(req); err != nil {
		return err
	}
	msg.Headers.Add("x-tx-id", args.TxID)
	p.s.Broadcast(args.PaymentReference, msg)
	return nil
}
