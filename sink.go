package substreams_sink_pubsub

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/hashicorp/go-multierror"
	"github.com/streamingfast/shutter"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
	"go.uber.org/zap"

	pbpubsub "github.com/streamingfast/substreams-sink-pubsub/pb/sf/substreams/sink/pubsub/v1"
)

type Sink struct {
	*shutter.Shutter
	*sink.Sinker
	logger     *zap.Logger
	client     *pubsub.Client
	topic      *pubsub.Topic
	cursorPath string
}

type Message struct {
	Data        []byte
	Attributes  map[string]string
	OrderingKey string
}

func NewSink(sinker *sink.Sinker, logger *zap.Logger, cursorPath string, client *pubsub.Client, topic *pubsub.Topic) *Sink {
	s := &Sink{
		Shutter:    shutter.New(),
		Sinker:     sinker,
		logger:     logger,
		client:     client,
		cursorPath: cursorPath,
		topic:      topic,
	}

	return s
}

func (s *Sink) Run(ctx context.Context) {
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

func (s *Sink) handleBlockScopedData(ctx context.Context, data *pbsubstreamsrpc.BlockScopedData, isLive *bool, cursor *sink.Cursor) error {

	publish := &pbpubsub.Publish{}

	err := data.Output.MapOutput.UnmarshalTo(publish)
	if err != nil {
		return fmt.Errorf("unmarshalling output: %w", err)
	}

	blockNum := data.Clock.Number
	messages := generateBlockScopedMessages(publish, cursor, blockNum)

	err = s.publishMessages(ctx, messages)
	if err != nil {
		return fmt.Errorf("publishing messages: %w", err)
	}

	err = s.saveCursor(cursor)
	if err != nil {
		return fmt.Errorf("saving cursor: %w", err)
	}

	return nil
}

func generateBlockScopedMessages(publish *pbpubsub.Publish, cursor *sink.Cursor, blockNum uint64) []*pubsub.Message {
	var messages []*pubsub.Message
	var indexCounter int
	for _, message := range publish.Messages {
		attributes := make(map[string]string)
		for _, attribute := range message.Attributes {
			attributes[attribute.Key] = attribute.Value
		}

		attributes["cursor"] = cursor.String()

		msg := &pubsub.Message{
			Data:       message.Data,
			Attributes: attributes,
		}

		messages = append(messages, msg)
		indexCounter++
	}

	return messages
}

func (s *Sink) handleBlockUndoSignal(ctx context.Context, data *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	lastValidBlockNum := data.LastValidBlock.Number

	messages := generateUndoBlockMessages(lastValidBlockNum, cursor)

	err := s.publishMessages(ctx, messages)
	if err != nil {
		return fmt.Errorf("publishing messages: %w", err)
	}

	err = s.saveCursor(cursor)
	if err != nil {
		return fmt.Errorf("saving cursor: %w", err)
	}

	return nil
}

func (s *Sink) loadCursor() (*sink.Cursor, error) {
	fpath := filepath.Join(s.cursorPath, "cursor.json")

	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	cursorData, err := os.ReadFile(fpath)

	if err != nil {
		return nil, fmt.Errorf("reading cursor file: %w", err)
	}

	cursorString := string(cursorData)
	cursor, err := sink.NewCursor(cursorString)
	if err != nil {
		return nil, fmt.Errorf("parsing cursor: %w", err)
	}

	return cursor, nil
}

func (s *Sink) saveCursor(c *sink.Cursor) error {
	cursorString := c.String()

	err := os.MkdirAll(s.cursorPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("making state store path: %w", err)
	}

	fpath := filepath.Join(s.cursorPath, "cursor.json")

	err = os.WriteFile(fpath, []byte(cursorString), os.ModePerm)
	if err != nil {
		return fmt.Errorf("writing cursor file: %w", err)
	}

	return nil
}

func (s *Sink) publishMessages(ctx context.Context, messages []*pubsub.Message) error {
	var results []*pubsub.PublishResult

	for _, message := range messages {
		result := s.topic.Publish(ctx, message)
		results = append(results, result)
	}

	meg := multierror.Group{}
	for _, res := range results {
		res := res
		meg.Go(func() error {
			_, err := res.Get(ctx)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := meg.Wait(); err != nil {
		return fmt.Errorf("handling result error: %w", err)
	}
	return nil
}
func generateUndoBlockMessages(lastValidBlockNum uint64, cursor *sink.Cursor) []*pubsub.Message {
	attributes := make(map[string]string)
	attributes["LastValidBlock"] = strconv.FormatUint(lastValidBlockNum, 10)
	attributes["Step"] = "Undo"
	attributes["Cursor"] = cursor.String()

	msg := &pubsub.Message{
		Data:       nil,
		Attributes: attributes,
	}

	messages := []*pubsub.Message{msg}

	return messages
}
