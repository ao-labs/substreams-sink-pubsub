package substreams_sink_pubsub

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"context"
	"fmt"
	"github.com/streamingfast/logging"
	"github.com/streamingfast/shutter"
	pbpubsub "github.com/streamingfast/substreams-sink-pubsub/pb/github.com/streamingfast/substreams-sink-pubsub/pb/substreams/sink/pubsub/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

var logger, _ = logging.ApplicationLogger("test", "test")

func TestHandlePublishOperations(t *testing.T) {

	ctx := context.Background()
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()

	// Connect to the server without using TLS.
	conn, err := grpc.Dial(srv.Addr, grpc.WithInsecure())
	require.NoError(t, err)

	defer conn.Close()

	cases := []struct {
		name            string
		operations      *pbpubsub.PublishOperations
		expectedResults []string
	}{
		{
			name: "sunny path",
			operations: &pbpubsub.PublishOperations{
				PublishOperations: []*pbpubsub.PublishOperation{
					{
						TopicID: "topic",
						Message: &pbpubsub.Message{
							Data: []byte("data"),
						},
					},
					{
						TopicID: "topic",
						Message: &pbpubsub.Message{
							Data: []byte("fiouu"),
						},
					},
				},
			},
			expectedResults: []string{"data", "fiouu"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, err := pubsub.NewClient(ctx, "project", option.WithGRPCConn(conn))

			client.Topic("topic")
			require.NoError(t, err)
			defer client.Close()

			topic, err := client.CreateTopic(ctx, "topic")
			if err != nil {
				require.NoError(t, err)
			}

			topics := make(map[string]*pubsub.Topic)
			topics["topic"] = topic

			sink := &Sink{
				Shutter:    shutter.New(),
				Sinker:     nil,
				logger:     logger,
				client:     client,
				topics:     topics,
				cursorPath: "",
			}

			subscription, err := client.CreateSubscription(ctx, "sub", pubsub.SubscriptionConfig{
				Topic:              topic,
				PushConfig:         pubsub.PushConfig{},
				BigQueryConfig:     pubsub.BigQueryConfig{},
				CloudStorageConfig: pubsub.CloudStorageConfig{},
				AckDeadline:        0,
			})
			if err != nil {
				require.NoError(t, err)
			}

			results := []string{}
			done := make(chan interface{})

			go func() {
				err = subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
					results = append(results, string(m.Data))
					m.Ack()
					if len(results) == len(c.expectedResults) {
						close(done)
					}
				})
				if err != nil {
					if s, ok := status.FromError(err); ok {
						if s.Code() == codes.Canceled {
							return
						}
					}

					fmt.Printf("Error: %T\n", err)
					require.NoError(t, err)
				}

			}()

			err = sink.handlePublishOperations(ctx, c.operations)
			if err != nil {
				require.NoError(t, err)
			}

			expire := time.After(1 * time.Second)

			select {
			case <-done:
			case <-expire:
				t.Fatal(fmt.Sprintf("expected %d messages, got %d", len(c.expectedResults), len(results)))
			}

			for _, res := range results {
				fmt.Println("Result ", res)
			}

			require.NoError(t, err)
		})

	}

}
