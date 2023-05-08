// SPDX-License-Identifier: AGPL-3.0-only

package querier

import (
	"errors"
	"fmt"
	"io"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb/chunkenc"

	"github.com/grafana/mimir/pkg/ingester/client"
	"github.com/grafana/mimir/pkg/storage/series"
	"github.com/grafana/mimir/pkg/util/chunkcompat"
)

type streamingChunkSeries struct {
	labels            labels.Labels
	chunkIteratorFunc chunkIteratorFunc
	mint, maxt        int64
	sources           []StreamingSeriesSource
}

func (s *streamingChunkSeries) Labels() labels.Labels {
	return s.labels
}

func (s *streamingChunkSeries) Iterator(_ chunkenc.Iterator) chunkenc.Iterator {
	var rawChunks []client.Chunk // TODO: guess a size for this?

	for _, source := range s.sources {
		c, err := source.StreamReader.GetChunks(source.SeriesIndex)

		if err != nil {
			return series.NewErrIterator(err)
		}

		rawChunks = client.AccumulateChunks(rawChunks, c)
	}

	chunks, err := chunkcompat.FromChunks(s.labels, rawChunks)
	if err != nil {
		return series.NewErrIterator(err)
	}

	return s.chunkIteratorFunc(chunks, model.Time(s.mint), model.Time(s.maxt))
}

type SeriesChunksStreamReader struct {
	client              client.Ingester_QueryStreamClient
	seriesBatchChan     chan []streamedIngesterSeries
	seriesBatch         []streamedIngesterSeries
	expectedSeriesCount int
	seriesBufferSize    int
}

func NewSeriesStreamer(client client.Ingester_QueryStreamClient, expectedSeriesCount int, seriesBufferSize int) *SeriesChunksStreamReader {
	return &SeriesChunksStreamReader{
		client:              client,
		expectedSeriesCount: expectedSeriesCount,
		seriesBufferSize:    seriesBufferSize,
	}
}

type streamedIngesterSeries struct {
	chunks      []client.Chunk
	seriesIndex int
	err         error
}

// StartBuffering begins streaming series' chunks from the ingester associated with
// this SeriesChunksStreamReader. Once all series have been consumed with GetChunks, all resources
// associated with this SeriesChunksStreamReader are cleaned up.
// If an error occurs while streaming, a subsequent call to GetChunks will return an error.
// To cancel buffering, cancel the context associated with this SeriesChunksStreamReader's client.Ingester_QueryStreamClient.
func (s *SeriesChunksStreamReader) StartBuffering() {
	s.seriesBatchChan = make(chan []streamedIngesterSeries, 1)
	ctxDone := s.client.Context().Done()

	// Why does this exist?
	// We want to make sure that the goroutine below is never leaked.
	// The goroutine below could be leaked if nothing is reading from the buffer but it's still trying to send
	// more series to a full buffer: it would block forever.
	// So, here, we try to send the series to the buffer if we can, but if the context is cancelled, then we give up.
	// This only works correctly if the context is cancelled when the query request is complete (or cancelled),
	// which is true at the time of writing.
	writeToBufferOrAbort := func(batch []streamedIngesterSeries) bool {
		select {
		case <-ctxDone:
			return false
		case s.seriesBatchChan <- batch:
			return true
		}
	}

	tryToWriteErrorToBuffer := func(err error) {
		writeToBufferOrAbort([]streamedIngesterSeries{{err: err}})
	}

	go func() {
		defer s.client.CloseSend() //nolint:errcheck
		defer close(s.seriesBatchChan)

		nextSeriesIndex := 0
		currentBatch := make([]streamedIngesterSeries, 0, s.seriesBufferSize)

		for {
			msg, err := s.client.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					if nextSeriesIndex < s.expectedSeriesCount {
						tryToWriteErrorToBuffer(fmt.Errorf("expected to receive %v series, but got EOF after receiving %v series", s.expectedSeriesCount, nextSeriesIndex))
					} else if len(currentBatch) > 0 {
						writeToBufferOrAbort(currentBatch)
					}
				} else {
					tryToWriteErrorToBuffer(fmt.Errorf("expected to receive %v series, but got EOF after receiving %v series", s.expectedSeriesCount, nextSeriesIndex))
				}

				return
			}

			for _, series := range msg.SeriesChunks {
				if nextSeriesIndex >= s.expectedSeriesCount {
					tryToWriteErrorToBuffer(fmt.Errorf("expected to receive only %v series, but received more than this", s.expectedSeriesCount))
					return
				}

				currentBatch = append(currentBatch, streamedIngesterSeries{chunks: series.Chunks, seriesIndex: nextSeriesIndex})

				if len(currentBatch) == s.seriesBufferSize {
					if !writeToBufferOrAbort(currentBatch) {
						return
					}

					currentBatch = currentBatch[:0]
				}

				nextSeriesIndex++
			}
		}
	}()
}

// GetChunks returns the chunks for the series with index seriesIndex.
// This method must be called with monotonically increasing values of seriesIndex.
func (s *SeriesChunksStreamReader) GetChunks(seriesIndex int) ([]client.Chunk, error) {
	if len(s.seriesBatch) == 0 {
		s.seriesBatch = <-s.seriesBatchChan

		if len(s.seriesBatch) == 0 {
			// If the context has been cancelled, report the cancellation.
			// Note that we only check this if there are no series in the buffer as the context is always cancelled
			// at the end of a successful request - so if we checked for an error even if there are series in the
			// buffer, we might incorrectly report that the context has been cancelled, when in fact the request
			// has concluded as expected.
			if err := s.client.Context().Err(); err != nil {
				return nil, err
			}

			return nil, fmt.Errorf("attempted to read series at index %v from stream, but the stream has already been exhausted", seriesIndex)
		}
	}

	series := s.seriesBatch[0]

	if len(s.seriesBatch) > 1 {
		s.seriesBatch = s.seriesBatch[1:]
	} else {
		s.seriesBatch = nil
	}

	if series.err != nil {
		return nil, fmt.Errorf("attempted to read series at index %v from stream, but the stream has failed: %w", seriesIndex, series.err)
	}

	if series.seriesIndex != seriesIndex {
		return nil, fmt.Errorf("attempted to read series at index %v from stream, but the stream has series with index %v", seriesIndex, series.seriesIndex)
	}

	return series.chunks, nil
}
