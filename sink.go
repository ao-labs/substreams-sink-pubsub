package substreams_sink_pubsub

import (
	"context"
	"time"

	"github.com/streamingfast/logging"
	"github.com/streamingfast/shutter"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
	"go.uber.org/zap"
)

type Sink struct {
	*shutter.Shutter
	*sink.Sinker
	logger *zap.Logger
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

	var cursor *sink.Cursor
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
	panic("implement me")
}

func (s *Sink) handleBlockUndoSignal(ctx context.Context, data *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	panic("implement me")
}
