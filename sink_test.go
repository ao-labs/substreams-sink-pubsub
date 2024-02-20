package substreams_sink_pubsub

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/streamingfast/logging"
	"github.com/streamingfast/shutter"
	pbpubsub "github.com/streamingfast/substreams-sink-pubsub/pb/substreams/sink/pubsub/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
						TopicID:     "topic",
						OrderingKey: "key.2",
						Message: &pbpubsub.Message{
							Data: []byte("data.99"),
						},
					},
					{
						TopicID:     "topic",
						OrderingKey: "key.1",
						Message: &pbpubsub.Message{
							Data: []byte("data.1"),
						},
					},
				},
			},
			expectedResults: []string{"data.99", "data.1"},
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
			topic.EnableMessageOrdering = true

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
				Topic:                 topic,
				AckDeadline:           10 * time.Second,
				EnableMessageOrdering: true,
			})
			if err != nil {
				require.NoError(t, err)
			}

			expectedResultCount := len(c.expectedResults)
			var results []string
			var lock sync.Mutex
			done := make(chan interface{})

			go func() {
				err = subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
					lock.Lock()
					results = append(results, string(m.Data))
					if len(results) == expectedResultCount {
						fmt.Println("Closing done channel")
						close(done)
						fmt.Println("channel closed")
					}
					m.Ack()
					lock.Unlock()
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

			err = sink.handlePublishOperations(ctx, c.operations, 1)
			require.NoError(t, err)

			expire := time.After(1 * time.Second)
			select {
			case <-done:
			case <-expire:
				t.Fatal(fmt.Sprintf("Timeout: expected %d messages, got %d", len(c.expectedResults), len(results)))
			}

			sort.Strings(results)
			sort.Strings(c.expectedResults)

			require.Equal(t, c.expectedResults, results)
		})
	}
}
