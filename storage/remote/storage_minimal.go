// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build minimal

package remote

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/model/exemplar"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/metadata"
	"github.com/prometheus/prometheus/scrape"
	"github.com/prometheus/prometheus/storage"
)

// Storage is a no-op implementation for minimal builds.
type Storage struct{}

// startTimeCallback is a callback func that return the oldest timestamp stored in a storage.
type startTimeCallback func() (int64, error)

// ReadyScrapeManager is an interface for getting the scrape manager.
type ReadyScrapeManager interface {
	Get() (*scrape.Manager, error)
}

// NewStorage returns a no-op storage for minimal builds.
func NewStorage(_ *slog.Logger, _ prometheus.Registerer, _ startTimeCallback, _ string, _ time.Duration, _ ReadyScrapeManager, _ bool) *Storage {
	return &Storage{}
}

// ApplyConfig is a no-op for minimal builds.
func (s *Storage) ApplyConfig(_ *config.Config) error {
	return nil
}

// Close is a no-op for minimal builds.
func (s *Storage) Close() error {
	return nil
}

// Querier returns a no-op querier for minimal builds.
func (s *Storage) Querier(_ int64, _ int64) (storage.Querier, error) {
	return storage.NoopQuerier(), nil
}

// ChunkQuerier returns a no-op chunk querier for minimal builds.
func (s *Storage) ChunkQuerier(_ int64, _ int64) (storage.ChunkQuerier, error) {
	return storage.NoopChunkedQuerier(), nil
}

// StartTime returns 0 for minimal builds.
func (s *Storage) StartTime() (int64, error) {
	return 0, nil
}

// Appender returns a no-op appender for minimal builds.
func (s *Storage) Appender(_ context.Context) storage.Appender {
	return noopAppender{}
}

// ExemplarQuerier returns nil for minimal builds.
func (s *Storage) ExemplarQuerier(_ context.Context) (storage.ExemplarQuerier, error) {
	return nil, nil
}

// LowestSentTimestamp returns 0 for minimal builds.
func (s *Storage) LowestSentTimestamp() int64 {
	return 0
}

// noopAppender is a no-op appender.
type noopAppender struct{}

func (noopAppender) Append(_ storage.SeriesRef, _ labels.Labels, _ int64, _ float64) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppender) SetOptions(_ *storage.AppendOptions) {}

func (noopAppender) AppendExemplar(_ storage.SeriesRef, _ labels.Labels, _ exemplar.Exemplar) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppender) AppendHistogram(_ storage.SeriesRef, _ labels.Labels, _ int64, _ *histogram.Histogram, _ *histogram.FloatHistogram) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppender) AppendHistogramSTZeroSample(_ storage.SeriesRef, _ labels.Labels, _, _ int64, _ *histogram.Histogram, _ *histogram.FloatHistogram) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppender) UpdateMetadata(_ storage.SeriesRef, _ labels.Labels, _ metadata.Metadata) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppender) AppendSTZeroSample(_ storage.SeriesRef, _ labels.Labels, _, _ int64) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppender) AppendCTZeroSample(_ storage.SeriesRef, _ labels.Labels, _, _ int64) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppender) Commit() error {
	return nil
}

func (noopAppender) Rollback() error {
	return nil
}

// OTLPOptions is a placeholder type for minimal builds.
type OTLPOptions struct {
	ConvertDelta            bool
	NativeDelta             bool
	LookbackDelta           time.Duration
	EnableTypeAndUnitLabels bool
}

// NewReadHandler returns a no-op handler for minimal builds.
func NewReadHandler(_ *slog.Logger, _ prometheus.Registerer, _ storage.SampleAndChunkQueryable, _ func() config.Config, _, _, _ int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
}

// NewWriteHandler returns a no-op handler for minimal builds.
func NewWriteHandler(_ *slog.Logger, _ prometheus.Registerer, _ storage.Appendable, _ interface{}, _, _, _ bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
}

// NewOTLPWriteHandler returns a no-op handler for minimal builds.
func NewOTLPWriteHandler(_ *slog.Logger, _ prometheus.Registerer, _ storage.AppendableV2, _ func() config.Config, _ OTLPOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
}

// Notify is a no-op for minimal builds.
func (s *Storage) Notify() {}

// AppenderV2 returns a no-op appenderV2 for minimal builds.
func (s *Storage) AppenderV2(_ context.Context) storage.AppenderV2 {
	return noopAppenderV2{}
}

// noopAppenderV2 is a no-op appenderV2.
type noopAppenderV2 struct{}

func (noopAppenderV2) Append(_ storage.SeriesRef, _ labels.Labels, _, _ int64, _ float64, _ *histogram.Histogram, _ *histogram.FloatHistogram, _ storage.AppendV2Options) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppenderV2) Commit() error {
	return nil
}

func (noopAppenderV2) Rollback() error {
	return nil
}
