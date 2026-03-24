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
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/scrape"
	"github.com/prometheus/prometheus/storage"
)

// MinimalStorage is a no-op implementation for minimal builds.
type MinimalStorage struct{}

// startTimeCallback is a callback func that return the oldest timestamp stored in a storage.
type startTimeCallback func() (int64, error)

// ReadyScrapeManager is an interface for getting the scrape manager.
type ReadyScrapeManager interface {
	Get() (*scrape.Manager, error)
}

// NewStorage returns a no-op storage for minimal builds.
func NewStorage(_ *slog.Logger, _ prometheus.Registerer, _ startTimeCallback, _ string, _ time.Duration, _ ReadyScrapeManager, _ bool) *MinimalStorage {
	return &MinimalStorage{}
}

// ApplyConfig is a no-op for minimal builds.
func (s *MinimalStorage) ApplyConfig(_ *config.Config) error {
	return nil
}

// Close is a no-op for minimal builds.
func (s *MinimalStorage) Close() error {
	return nil
}

// Querier returns a no-op querier for minimal builds.
func (s *MinimalStorage) Querier(_ int64, _ int64) (storage.Querier, error) {
	return storage.NoopQuerier(), nil
}

// ChunkQuerier returns a no-op chunk querier for minimal builds.
func (s *MinimalStorage) ChunkQuerier(_ int64, _ int64) (storage.ChunkQuerier, error) {
	return storage.NoopChunkedQuerier(), nil
}

// StartTime returns 0 for minimal builds.
func (s *MinimalStorage) StartTime() (int64, error) {
	return 0, nil
}

// Appender returns a no-op appender for minimal builds.
func (s *MinimalStorage) Appender(_ context.Context) storage.Appender {
	return noopAppender{}
}

// ExemplarQuerier returns nil for minimal builds.
func (s *MinimalStorage) ExemplarQuerier(_ context.Context) (storage.ExemplarQuerier, error) {
	return nil, nil
}

// noopAppender is a no-op appender.
type noopAppender struct{}

func (noopAppender) Append(storage.SeriesRef, ...interface{}) (storage.SeriesRef, error) {
	return 0, nil
}

func (noopAppender) Commit() error {
	return nil
}

func (noopAppender) Rollback() error {
	return nil
}
