package client

import (
	"context"

	corev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/core"
	foundationv1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/foundation"
	profilev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/profile"
	"github.com/appadeia/barista/barista-go/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn          *grpc.ClientConn
	CoreKit       corev1.CoreServiceClient
	FoundationKit foundationv1.FoundationServiceClient
	Profilekit    profilev1.ProfileServiceClient
	sessionToken  string
	userID        uint64
	onceHandlers  []func(*corev1.GuildEvent)
}

func (c Client) Context() context.Context {
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "auth", c.sessionToken)

	return ctx
}

func (c *Client) HandleOnce(f func(*corev1.GuildEvent)) {
	c.onceHandlers = append(c.onceHandlers, f)
}

func (c *Client) GuildEvents(guildID uint64) (chan *corev1.GuildEvent, error) {
	stream, err := c.CoreKit.StreamGuildEvents(c.Context(), &corev1.StreamGuildEventsRequest{
		Location: &corev1.Location{
			GuildId: guildID,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "GuildEvents: failed to open guild events stream")
		return nil, err
	}

	channel := make(chan *corev1.GuildEvent)
	go func() {
		for {
			ev, err := stream.Recv()
			if err != nil {
				log.Error("%+v", err)
				close(channel)
				return
			}
			for _, handler := range c.onceHandlers {
				handler(ev)
			}
			c.onceHandlers = []func(*corev1.GuildEvent){}
			channel <- ev
		}
	}()

	return channel, nil
}
