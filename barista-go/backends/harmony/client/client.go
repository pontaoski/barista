package client

import (
	"context"

	corev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/core"
	foundationv1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/foundation"
	profilev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/profile"
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
	onceHandlers  []func(*corev1.Event)
}

func (c Client) Context() context.Context {
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "auth", c.sessionToken)

	return ctx
}

func (c *Client) HandleOnce(f func(*corev1.Event)) {
	c.onceHandlers = append(c.onceHandlers, f)
}
