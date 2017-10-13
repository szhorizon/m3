// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aggregator

import (
	"sync"
	"time"

	"github.com/m3db/m3aggregator/aggregation/quantile/cm"
	"github.com/m3db/m3aggregator/sharding"
	"github.com/m3db/m3metrics/metric/aggregated"
	"github.com/m3db/m3metrics/policy"
	"github.com/m3db/m3x/clock"
	"github.com/m3db/m3x/instrument"
	"github.com/m3db/m3x/pool"

	"github.com/uber-go/tally"
)

// CounterElemAlloc allocates a new counter element
type CounterElemAlloc func() *CounterElem

// CounterElemPool provides a pool of counter elements
type CounterElemPool interface {
	// Init initializes the counter element pool
	Init(alloc CounterElemAlloc)

	// Get gets a counter element from the pool
	Get() *CounterElem

	// Put returns a counter element to the pool
	Put(value *CounterElem)
}

// TimerElemAlloc allocates a new timer element
type TimerElemAlloc func() *TimerElem

// TimerElemPool provides a pool of timer elements
type TimerElemPool interface {
	// Init initializes the timer element pool
	Init(alloc TimerElemAlloc)

	// Get gets a timer element from the pool
	Get() *TimerElem

	// Put returns a timer element to the pool
	Put(value *TimerElem)
}

// GaugeElemAlloc allocates a new gauge element
type GaugeElemAlloc func() *GaugeElem

// GaugeElemPool provides a pool of gauge elements
type GaugeElemPool interface {
	// Init initializes the gauge element pool
	Init(alloc GaugeElemAlloc)

	// Get gets a gauge element from the pool
	Get() *GaugeElem

	// Put returns a gauge element to the pool
	Put(value *GaugeElem)
}

// EntryAlloc allocates a new entry
type EntryAlloc func() *Entry

// EntryPool provides a pool of entries
type EntryPool interface {
	// Init initializes the entry pool
	Init(alloc EntryAlloc)

	// Get gets a entry from the pool
	Get() *Entry

	// Put returns a entry to the pool
	Put(value *Entry)
}

// QuantileSuffixFn returns the byte-slice suffix for a quantile value
type QuantileSuffixFn func(quantile float64) []byte

// Handler handles aggregated metrics alongside their policies.
type Handler interface {
	// NewWriter creates a new writer for writing aggregated metrics and policies.
	NewWriter(scope tally.Scope) (Writer, error)

	// Close closes the handler.
	Close()
}

// Writer writes aggregated metrics alongside their policies.
type Writer interface {
	// Write writes an aggregated metric alongside its storage policy.
	Write(mp aggregated.ChunkedMetricWithStoragePolicy) error

	// Flush flushes data buffered in the writer to backend.
	Flush() error

	// Close closes the writer.
	Close() error
}

// Options provide a set of base and derived options for the aggregator
type Options interface {
	/// Read-write base options

	// SetDefaultCounterAggregationTypes sets the counter aggregation types.
	SetDefaultCounterAggregationTypes(value policy.AggregationTypes) Options

	// DefaultCounterAggregationTypes returns the aggregation types for counters.
	DefaultCounterAggregationTypes() policy.AggregationTypes

	// SetDefaultTimerAggregationTypes sets the timer aggregation types.
	SetDefaultTimerAggregationTypes(value policy.AggregationTypes) Options

	// DefaultTimerAggregationTypes returns the aggregation types for timers.
	DefaultTimerAggregationTypes() policy.AggregationTypes

	// SetDefaultGaugeAggregationTypes sets the gauge aggregation types.
	SetDefaultGaugeAggregationTypes(value policy.AggregationTypes) Options

	// DefaultGaugeAggregationTypes returns the aggregation types for gauges.
	DefaultGaugeAggregationTypes() policy.AggregationTypes

	// SetMetricPrefix sets the common prefix for all metric types
	SetMetricPrefix(value []byte) Options

	// MetricPrefix returns the common prefix for all metric types
	MetricPrefix() []byte

	// SetCounterPrefix sets the prefix for counters
	SetCounterPrefix(value []byte) Options

	// CounterPrefix returns the prefix for counters
	CounterPrefix() []byte

	// SetTimerPrefix sets the prefix for timers
	SetTimerPrefix(value []byte) Options

	// TimerPrefix returns the prefix for timers
	TimerPrefix() []byte

	// SetGaugePrefix sets the prefix for gauges
	SetGaugePrefix(value []byte) Options

	// GaugePrefix returns the prefix for gauges
	GaugePrefix() []byte

	// SetAggregationLastSuffix sets the suffix for aggregation type last.
	SetAggregationLastSuffix(value []byte) Options

	// AggregationLastSuffix returns the suffix for aggregation type last.
	AggregationLastSuffix() []byte

	// SetAggregationMinSuffix sets the suffix for aggregation type min.
	SetAggregationMinSuffix(value []byte) Options

	// AggregationMinSuffix returns the suffix for aggregation type min.
	AggregationMinSuffix() []byte

	// SetAggregationMaxSuffix sets the suffix for aggregation type max.
	SetAggregationMaxSuffix(value []byte) Options

	// AggregationMaxSuffix returns the suffix for aggregation type max.
	AggregationMaxSuffix() []byte

	// SetAggregationMeanSuffix sets the suffix for aggregation type mean.
	SetAggregationMeanSuffix(value []byte) Options

	// AggregationMeanSuffix returns the suffix for aggregation type mean.
	AggregationMeanSuffix() []byte

	// SetAggregationMedianSuffix sets the suffix for aggregation type median.
	SetAggregationMedianSuffix(value []byte) Options

	// AggregationMedianSuffix returns the suffix for aggregation type median.
	AggregationMedianSuffix() []byte

	// SetAggregationCountSuffix sets the suffix for aggregation type count.
	SetAggregationCountSuffix(value []byte) Options

	// AggregationCountSuffix returns the suffix for aggregation type count.
	AggregationCountSuffix() []byte

	// SetAggregationSumSuffix sets the suffix for aggregation type sum.
	SetAggregationSumSuffix(value []byte) Options

	// AggregationSumSuffix returns the suffix for aggregation type sum.
	AggregationSumSuffix() []byte

	// SetAggregationSumSqSuffix sets the suffix for aggregation type sum square.
	SetAggregationSumSqSuffix(value []byte) Options

	// AggregationSumSqSuffix returns the suffix for aggregation type sum square.
	AggregationSumSqSuffix() []byte

	// SetAggregationStdevSuffix sets the suffix for aggregation type standard deviation.
	SetAggregationStdevSuffix(value []byte) Options

	// AggregationStdevSuffix returns the suffix for aggregation type standard deviation.
	AggregationStdevSuffix() []byte

	// SetTimerQuantileSuffixFn sets the quantile suffix function for timers.
	SetTimerQuantileSuffixFn(value QuantileSuffixFn) Options

	// TimerQuantileSuffixFn returns the quantile suffix function for timers.
	TimerQuantileSuffixFn() QuantileSuffixFn

	// SetTimeLock sets the time lock
	SetTimeLock(value *sync.RWMutex) Options

	// TimeLock returns the time lock
	TimeLock() *sync.RWMutex

	// SetClockOptions sets the clock options
	SetClockOptions(value clock.Options) Options

	// ClockOptions returns the clock options
	ClockOptions() clock.Options

	// SetInstrumentOptions sets the instrument options
	SetInstrumentOptions(value instrument.Options) Options

	// InstrumentOptions returns the instrument options
	InstrumentOptions() instrument.Options

	// SetStreamOptions sets the stream options
	SetStreamOptions(value cm.Options) Options

	// StreamOptions returns the stream options
	StreamOptions() cm.Options

	// SetPlacementManager sets the placement manager.
	SetPlacementManager(value PlacementManager) Options

	// PlacementManager returns the placement manager.
	PlacementManager() PlacementManager

	// SetShardFn sets the sharding function.
	SetShardFn(value sharding.ShardFn) Options

	// ShardFn returns the sharding function.
	ShardFn() sharding.ShardFn

	// SetBufferDurationBeforeShardCutover sets the duration for buffering writes before shard cutover.
	SetBufferDurationBeforeShardCutover(value time.Duration) Options

	// BufferDurationBeforeShardCutover returns the duration for buffering writes before shard cutover.
	BufferDurationBeforeShardCutover() time.Duration

	// SetBufferDurationAfterShardCutoff sets the duration for buffering writes after shard cutoff.
	SetBufferDurationAfterShardCutoff(value time.Duration) Options

	// BufferDurationAfterShardCutoff returns the duration for buffering writes after shard cutoff.
	BufferDurationAfterShardCutoff() time.Duration

	// SetFlushTimesManager sets the flush times manager.
	SetFlushTimesManager(value FlushTimesManager) Options

	// FlushTimesManager returns the flush times manager.
	FlushTimesManager() FlushTimesManager

	// SetElectionManager sets the election manager.
	SetElectionManager(value ElectionManager) Options

	// ElectionManager returns the election manager.
	ElectionManager() ElectionManager

	// SetFlushManager sets the flush manager.
	SetFlushManager(value FlushManager) Options

	// FlushManager returns the flush manager.
	FlushManager() FlushManager

	// SetMinFlushInterval sets the minimum flush interval
	SetMinFlushInterval(value time.Duration) Options

	// MinFlushInterval returns the minimum flush interval
	MinFlushInterval() time.Duration

	// SetFlushHandler sets the handler that flushes buffered encoders
	SetFlushHandler(value Handler) Options

	// FlushHandler returns the handler that flushes buffered encoders
	FlushHandler() Handler

	// SetEntryTTL sets the ttl for expiring stale entries
	SetEntryTTL(value time.Duration) Options

	// EntryTTL returns the ttl for expiring stale entries
	EntryTTL() time.Duration

	// SetEntryCheckInterval sets the interval for checking expired entries
	SetEntryCheckInterval(value time.Duration) Options

	// EntryCheckInterval returns the interval for checking expired entries
	EntryCheckInterval() time.Duration

	// SetEntryCheckBatchPercent sets the batch percentage for checking expired entries
	SetEntryCheckBatchPercent(value float64) Options

	// EntryCheckBatchPercent returns the batch percentage for checking expired entries
	EntryCheckBatchPercent() float64

	// SetMaxTimerBatchSizePerWrite sets the maximum timer batch size for each batched write.
	SetMaxTimerBatchSizePerWrite(value int) Options

	// MaxTimerBatchSizePerWrite returns the maximum timer batch size for each batched write.
	MaxTimerBatchSizePerWrite() int

	// SetDefaultPolicies sets the default policies.
	SetDefaultPolicies(value []policy.Policy) Options

	// DefaultPolicies returns the default policies.
	DefaultPolicies() []policy.Policy

	// SetResignTimeout sets the resign timeout.
	SetResignTimeout(value time.Duration) Options

	// ResignTimeout returns the resign timeout.
	ResignTimeout() time.Duration

	// SetEntryPool sets the entry pool
	SetEntryPool(value EntryPool) Options

	// EntryPool returns the entry pool
	EntryPool() EntryPool

	// SetCounterElemPool sets the counter element pool.
	SetCounterElemPool(value CounterElemPool) Options

	// CounterElemPool returns the counter element pool.
	CounterElemPool() CounterElemPool

	// SetTimerElemPool sets the timer element pool.
	SetTimerElemPool(value TimerElemPool) Options

	// TimerElemPool returns the timer element pool.
	TimerElemPool() TimerElemPool

	// SetGaugeElemPool sets the gauge element pool.
	SetGaugeElemPool(value GaugeElemPool) Options

	// GaugeElemPool returns the gauge element pool.
	GaugeElemPool() GaugeElemPool

	// SetAggregationTypesPool sets the aggregation types pool.
	SetAggregationTypesPool(pool policy.AggregationTypesPool) Options

	// AggregationTypesPool returns the aggregation types pool.
	AggregationTypesPool() policy.AggregationTypesPool

	// SetQuantilesPool sets the timer quantiles pool.
	SetQuantilesPool(pool pool.FloatsPool) Options

	// QuantilesPool returns the timer quantiles pool.
	QuantilesPool() pool.FloatsPool

	/// Write-only options.

	// SetCounterSuffixOverride sets the overrides for counter suffixes.
	SetCounterSuffixOverride(m map[policy.AggregationType][]byte) Options

	// SetTimerSuffixOverride sets the overrides for timer suffixes.
	SetTimerSuffixOverride(m map[policy.AggregationType][]byte) Options

	// SetGaugeSuffixOverride sets the overrides for gauge suffixes.
	SetGaugeSuffixOverride(m map[policy.AggregationType][]byte) Options

	/// Read-only derived options.

	// FullCounterPrefix returns the full prefix for counters
	FullCounterPrefix() []byte

	// FullTimerPrefix returns the full prefix for timers
	FullTimerPrefix() []byte

	// FullGaugePrefix returns the full prefix for gauges
	FullGaugePrefix() []byte

	// TimerQuantiles returns the quantiles for timers.
	TimerQuantiles() []float64

	// DefaultCounterAggregationSuffixes returns the suffix for
	// default counter aggregation types.
	DefaultCounterAggregationSuffixes() [][]byte

	// DefaultTimerAggregationSuffixes returns the suffix for
	// default timer aggregation types.
	DefaultTimerAggregationSuffixes() [][]byte

	// DefaultGaugeAggregationSuffixes returns the suffix for
	// default gauge aggregation types.
	DefaultGaugeAggregationSuffixes() [][]byte

	// Suffix returns the suffix for the aggregation type for counters.
	SuffixForCounter(value policy.AggregationType) []byte

	// Suffix returns the suffix for the aggregation type for timers.
	SuffixForTimer(value policy.AggregationType) []byte

	// Suffix returns the suffix for the aggregation type for gauges.
	SuffixForGauge(value policy.AggregationType) []byte
}