package client

import (
	"context"
	"reflect"
	"strings"

	corev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/core"
	foundationv1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/foundation"
	profilev1 "github.com/appadeia/barista/barista-go/backends/harmony/gen/profile"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type GuildEvent struct {
	*corev1.GuildEvent
	Client *Client
}

type MainClient struct {
	Client
	homeserver string
	subclients map[string]*Client
	streams    map[chan *corev1.GuildEvent]*Client
}

func NewClient(homeserver, email, password string) (client *MainClient, err error) {
	client = &MainClient{
		Client: Client{},
	}

	client.homeserver = homeserver
	client.subclients = make(map[string]*Client)
	client.streams = make(map[chan *corev1.GuildEvent]*Client)

	client.conn, err = grpc.Dial(homeserver, grpc.WithInsecure())
	if err != nil {
		client = nil
		err = errors.Wrap(err, "NewClient: failed to dial grpc")
		return
	}

	client.CoreKit = corev1.NewCoreServiceClient(client.conn)
	client.FoundationKit = foundationv1.NewFoundationServiceClient(client.conn)
	client.Profilekit = profilev1.NewProfileServiceClient(client.conn)

	session, err := client.FoundationKit.Login(context.Background(), &foundationv1.LoginRequest{
		Login: &foundationv1.LoginRequest_Local_{
			Local: &foundationv1.LoginRequest_Local{
				Email:    email,
				Password: []byte(password),
			},
		},
	})
	if err != nil {
		client = nil
		err = errors.Wrap(err, "NewClient: failed to login")
		return
	}

	client.sessionToken = session.SessionToken
	client.userID = session.UserId

	return
}

func (m *MainClient) ClientFor(homeserver string) (*Client, error) {
	if m.homeserver == homeserver || strings.Split(homeserver, ":")[0] == "localhost" {
		return &m.Client, nil
	}

	if val, ok := m.subclients[homeserver]; ok {
		return val, nil
	}

	federatedSession, err := m.FoundationKit.Federate(m.Context(), &foundationv1.FederateRequest{
		Target: homeserver,
	})

	if err != nil {
		return nil, err
	}

	client := new(Client)
	client.conn, err = grpc.Dial(homeserver, grpc.WithInsecure())
	if err != nil {
		err = errors.Wrap(err, "ClientFor: failed to dial grpc")
		return nil, err
	}

	client.CoreKit = corev1.NewCoreServiceClient(client.conn)
	client.FoundationKit = foundationv1.NewFoundationServiceClient(client.conn)
	client.Profilekit = profilev1.NewProfileServiceClient(client.conn)

	session, err := client.FoundationKit.Login(context.Background(), &foundationv1.LoginRequest{
		Login: &foundationv1.LoginRequest_Federated_{
			Federated: &foundationv1.LoginRequest_Federated{
				AuthToken: federatedSession.Token,
				Domain:    m.homeserver,
			},
		},
	})
	if err != nil {
		err = errors.Wrap(err, "ClientFor: failed to login")
		return nil, err
	}

	client.sessionToken = session.SessionToken
	client.userID = session.UserId
	return client, nil
}

func (m *MainClient) Start() (chan GuildEvent, error) {
	list, err := m.CoreKit.GetGuildList(m.Context(), &corev1.GetGuildListRequest{})
	if err != nil {
		err = errors.Wrap(err, "Start: failed to get guild list")
		return nil, err
	}

	cases := []reflect.SelectCase{}

	for _, guild := range list.Guilds {
		client, err := m.ClientFor(guild.Host)
		if err != nil {
			err = errors.Wrap(err, "Start: failed to get client for guild")
			return nil, err
		}

		stream, err := client.GuildEvents(guild.GuildId)
		if err != nil {
			err = errors.Wrap(err, "Start: failed to get guild events stream")
			return nil, err
		}

		m.streams[stream] = client
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(stream),
		})
	}

	channel := make(chan GuildEvent)
	go func() {
		for {
			i, v, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases[:i], cases[i+1:]...)
			}

			val := v.Interface().(*corev1.GuildEvent)

			channel <- GuildEvent{
				GuildEvent: val,
				Client:     m.streams[cases[i].Chan.Interface().(chan *corev1.GuildEvent)],
			}
		}
	}()

	return channel, nil
}
