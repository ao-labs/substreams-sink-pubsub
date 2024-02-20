package substreams_sink_pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/streamingfast/logging"
	"github.com/streamingfast/shutter"
	sink "github.com/streamingfast/substreams-sink"
	pbpubsub "github.com/streamingfast/substreams-sink-pubsub/pb/proto/substreams/sink/pubsub/v1"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
	"go.uber.org/zap"
)

type Sink struct {
	*shutter.Shutter
	*sink.Sinker
	logger     *zap.Logger
	client     *pubsub.Client
	topics     map[string]*pubsub.Topic
	cursorPath string
}

func NewSink(sinker *sink.Sinker, logger *zap.Logger, tracer logging.Tracer) *Sink {
	s := &Sink{
		Shutter: shutter.New(),
		Sinker:  sinker,
		logger:  logger,
	}

	s.OnTerminating(func(err error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.onTerminating(ctx, err)
	})
	return s
}

func (s *Sink) Run(ctx context.Context) {
	s.logger.Info("starting PubSub sink")
	s.Sinker.OnTerminating(s.Shutdown)
	s.OnTerminating(func(err error) {
		s.logger.Info("terminating")
		s.Sinker.Shutdown(err)
	})

	cursor, err := s.loadCursor()
	if err != nil {
		s.Shutdown(fmt.Errorf("loading cursor: %w", err))
	}

	s.logger.Info("starting PubSub sink", zap.Stringer("restarting_at", cursor.Block()))
	s.Sinker.Run(ctx, cursor, sink.NewSinkerHandlers(s.handleBlockScopedData, s.handleBlockUndoSignal))
}

func (s *Sink) onTerminating(ctx context.Context, err error) {
	//if s.lastCursor == nil || err != nil {
	//	return
	//}
	//
	//_ = s.operationDB.WriteCursor(ctx, s.lastCursor)
	panic("implement me")
}

func (s *Sink) handleBlockScopedData(ctx context.Context, data *pbsubstreamsrpc.BlockScopedData, isLive *bool, cursor *sink.Cursor) error {
	publishOperations := &pbpubsub.PublishOperations{}

	err := data.Output.MapOutput.UnmarshalTo(publishOperations)
	if err != nil {
		return fmt.Errorf("unmarshalling output: %w", err)
	}

	err = s.handlePublishOperations(ctx, publishOperations, data.Clock.Number)
	if err != nil {
		return fmt.Errorf("handling publish operations: %w", err)
	}
	err = s.saveCursor(cursor)
	if err != nil {
		return fmt.Errorf("saving cursor: %w", err)
	}

	return nil
}

func (s *Sink) handleBlockUndoSignal(ctx context.Context, data *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	panic("implement me")
}

func (s *Sink) loadCursor() (*sink.Cursor, error) {
	_, err := os.Stat(s.cursorPath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	cursorData, err := os.ReadFile(s.cursorPath)

	if err != nil {
		return nil, fmt.Errorf("reading cursor file: %w", err)
	}

	cursor := &sink.Cursor{}
	err = json.Unmarshal(cursorData, cursor)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling cursor: %w", err)
	}

	return cursor, nil
}

func (s *Sink) saveCursor(cursor *sink.Cursor) error {

	cursorBytes, err := json.Marshal(cursor)
	if err != nil {
		return fmt.Errorf("marshalling cursor: %w", err)
	}

	err = os.WriteFile(s.cursorPath, cursorBytes, 0644)
	if err != nil {
		return fmt.Errorf("writing cursor file: %w", err)
	}

	return nil
}

func (s *Sink) handlePublishOperations(ctx context.Context, publishOperations *pbpubsub.PublishOperations, blockNum uint64) error {
	var results []*pubsub.PublishResult

	for _, operation := range publishOperations.PublishOperations {
		if t, ok := s.topics[operation.TopicID]; ok {

			msg := &pubsub.Message{
				Data:        operation.Message.Data,
				OrderingKey: fmt.Sprintf("%d", blockNum),
			}
			if operation.OrderingKey != "" {
				msg.OrderingKey = operation.OrderingKey
			}

			result := t.Publish(ctx, msg)
			s.logger.Info("publishing message", zap.String("topic_id", operation.TopicID), zap.String("message", string(operation.Message.Data)))
			results = append(results, result)
			continue
		}

		return fmt.Errorf("topic %s not found", operation.TopicID)
	}

	var resultErrors []error

	for _, res := range results {
		_, err := res.Get(ctx)
		if err != nil {
			resultErrors = append(resultErrors, err)
		}
	}

	if len(resultErrors) != 0 {
		return fmt.Errorf("handling result error: %w", resultErrors[len(resultErrors)-1])
	}

	return nil
}

//TODO: Configure retry

//TODO: Activate compression when publishing
