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
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/m3db/m3metrics/metric/aggregated"
	"github.com/m3db/m3metrics/metric/unaggregated"
	"github.com/m3db/m3metrics/policy"
	"github.com/m3db/m3x/clock"

	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
)

func TestMetricListPushBack(t *testing.T) {
	l, err := newMetricList(testShard, time.Second, testOptions())
	require.NoError(t, err)
	elem := NewCounterElem(nil, policy.EmptyStoragePolicy, policy.DefaultAggregationTypes, l.opts)

	// Push a counter to the list
	e, err := l.PushBack(elem)
	require.NoError(t, err)
	require.Equal(t, 1, l.aggregations.Len())
	require.Equal(t, elem, e.Value.(*CounterElem))

	// Push a counter to a closed list should result in an error.
	l.Lock()
	l.closed = true
	l.Unlock()

	_, err = l.PushBack(elem)
	require.Equal(t, err, errListClosed)
}

func TestMetricListClose(t *testing.T) {
	var (
		registered   int
		unregistered int
	)
	mockFlushManager := &mockFlushManager{
		registerFn:   func(PeriodicFlusher) error { registered++; return nil },
		unregisterFn: func(PeriodicFlusher) error { unregistered++; return nil },
	}
	opts := testOptions().SetFlushManager(mockFlushManager)
	l, err := newMetricList(testShard, time.Second, opts)
	require.NoError(t, err)

	l.RLock()
	require.False(t, l.closed)
	l.RUnlock()
	require.Equal(t, 1, registered)

	l.Close()
	require.True(t, l.closed)
	require.Equal(t, 1, unregistered)

	// Close for a second time should have no impact.
	l.Close()
	require.True(t, l.closed)
	require.Equal(t, 1, registered)
	require.Equal(t, 1, unregistered)
}

func TestMetricListFlushWithRequests(t *testing.T) {
	var (
		now     = time.Unix(12345, 0)
		nowFn   = func() time.Time { return now }
		results []flushBeforeResult
	)
	opts := testOptions().SetClockOptions(clock.NewOptions().SetNowFn(nowFn))
	l, err := newMetricList(testShard, time.Second, opts)
	require.NoError(t, err)
	l.flushBeforeFn = func(beforeNanos int64, flushType flushType) {
		results = append(results, flushBeforeResult{
			beforeNanos: beforeNanos,
			flushType:   flushType,
		})
	}

	inputs := []struct {
		request  FlushRequest
		expected []flushBeforeResult
	}{
		{
			request: FlushRequest{
				CutoverNanos:      20000 * int64(time.Second),
				CutoffNanos:       30000 * int64(time.Second),
				BufferAfterCutoff: time.Second,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 12345 * int64(time.Second),
					flushType:   discardType,
				},
			},
		},
		{
			request: FlushRequest{
				CutoverNanos:      10000 * int64(time.Second),
				CutoffNanos:       30000 * int64(time.Second),
				BufferAfterCutoff: time.Second,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 10000 * int64(time.Second),
					flushType:   discardType,
				},
				{
					beforeNanos: 12345 * int64(time.Second),
					flushType:   consumeType,
				},
			},
		},
		{
			request: FlushRequest{
				CutoverNanos:      10000 * int64(time.Second),
				CutoffNanos:       12300 * int64(time.Second),
				BufferAfterCutoff: time.Minute,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 10000 * int64(time.Second),
					flushType:   discardType,
				},
				{
					beforeNanos: 12300 * int64(time.Second),
					flushType:   consumeType,
				},
			},
		},
		{
			request: FlushRequest{
				CutoverNanos:      10000 * int64(time.Second),
				CutoffNanos:       12300 * int64(time.Second),
				BufferAfterCutoff: 10 * time.Second,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 10000 * int64(time.Second),
					flushType:   discardType,
				},
				{
					beforeNanos: 12300 * int64(time.Second),
					flushType:   consumeType,
				},
				{
					beforeNanos: 12335 * int64(time.Second),
					flushType:   discardType,
				},
			},
		},
		{
			request: FlushRequest{
				CutoverNanos:      0,
				CutoffNanos:       30000 * int64(time.Second),
				BufferAfterCutoff: time.Second,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 12345 * int64(time.Second),
					flushType:   consumeType,
				},
			},
		},
	}
	for _, input := range inputs {
		results = results[:0]
		l.Flush(input.request)
		require.Equal(t, input.expected, results)
	}
}

func TestMetricListFlushConsumingAndCollectingElems(t *testing.T) {
	var (
		cutoverNanos = int64(0)
		cutoffNanos  = int64(math.MaxInt64)
		count        int
		flushLock    sync.Mutex
		flushed      []aggregated.ChunkedMetricWithStoragePolicy
	)

	// Intentionally cause a one-time error during encoding.
	writeFn := func(mp aggregated.ChunkedMetricWithStoragePolicy) error {
		flushLock.Lock()
		defer flushLock.Unlock()

		if count == 0 {
			count++
			return errors.New("foo")
		}
		flushed = append(flushed, mp)
		return nil
	}
	writer := &mockWriter{
		writeFn: writeFn,
		flushFn: func() error { return nil },
	}
	handler := &mockHandler{
		newWriterFn: func(tally.Scope) (Writer, error) { return writer, nil },
	}

	var now = time.Unix(216, 0).UnixNano()
	nowTs := time.Unix(0, now)
	clockOpts := clock.NewOptions().SetNowFn(func() time.Time {
		return time.Unix(0, atomic.LoadInt64(&now))
	})
	opts := testOptions().
		SetClockOptions(clockOpts).
		SetMinFlushInterval(0).
		SetFlushHandler(handler)

	l, err := newMetricList(testShard, 0, opts)
	require.NoError(t, err)
	l.resolution = testStoragePolicy.Resolution().Window

	// Intentionally cause a one-time error during encoding.
	elemPairs := []testElemPair{
		{
			elem:   NewCounterElem(testCounterID, testStoragePolicy, policy.DefaultAggregationTypes, opts),
			metric: testCounter,
		},
		{
			elem:   NewTimerElem(testBatchTimerID, testStoragePolicy, policy.DefaultAggregationTypes, opts),
			metric: testBatchTimer,
		},
		{
			elem:   NewGaugeElem(testGaugeID, testStoragePolicy, policy.DefaultAggregationTypes, opts),
			metric: testGauge,
		},
	}

	for _, ep := range elemPairs {
		require.NoError(t, ep.elem.AddMetric(nowTs, ep.metric))
		require.NoError(t, ep.elem.AddMetric(nowTs.Add(l.resolution), ep.metric))
		_, err := l.PushBack(ep.elem)
		require.NoError(t, err)
	}

	// Force a flush.
	l.Flush(FlushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert nothing has been flushed.
	flushLock.Lock()
	require.Equal(t, 0, len(flushed))
	flushLock.Unlock()

	for i := 0; i < 2; i++ {
		// Move the time forward by one aggregation interval.
		nowTs = nowTs.Add(l.resolution)
		atomic.StoreInt64(&now, nowTs.UnixNano())

		// Force a flush.
		l.Flush(FlushRequest{
			CutoverNanos: cutoverNanos,
			CutoffNanos:  cutoffNanos,
		})

		var expected []testAggMetric
		alignedStart := nowTs.Truncate(l.resolution).UnixNano()
		expected = append(expected, expectedAggMetricsForCounter(alignedStart, testStoragePolicy, policy.DefaultAggregationTypes)...)
		expected = append(expected, expectedAggMetricsForTimer(alignedStart, testStoragePolicy, policy.DefaultAggregationTypes)...)
		expected = append(expected, expectedAggMetricsForGauge(alignedStart, testStoragePolicy, policy.DefaultAggregationTypes)...)

		// Skip the first item because we intentionally triggered
		// an encoder error when encoding the first item.
		if i == 0 {
			expected = expected[1:]
		}

		flushLock.Lock()
		require.NotNil(t, flushed)
		validateFlushed(t, expected, flushed)
		flushed = flushed[:0]
		flushLock.Unlock()
	}

	// Move the time forward by one aggregation interval.
	nowTs = nowTs.Add(l.resolution)
	atomic.StoreInt64(&now, nowTs.UnixNano())

	// Force a flush.
	l.Flush(FlushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert nothing has been flushed.
	flushLock.Lock()
	require.Equal(t, 0, len(flushed))
	flushLock.Unlock()
	require.Equal(t, 3, l.aggregations.Len())

	// Mark all elements as tombstoned.
	for e := l.aggregations.Front(); e != nil; e = e.Next() {
		e.Value.(metricElem).MarkAsTombstoned()
	}

	// Move the time forward and force a flush.
	nowTs = nowTs.Add(l.resolution)
	atomic.StoreInt64(&now, nowTs.UnixNano())
	l.Flush(FlushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert all elements have been collected.
	require.Equal(t, 0, l.aggregations.Len())

	require.Equal(t, l.lastFlushedNanos, nowTs.Truncate(l.resolution).UnixNano())
}

func TestMetricListFlushBeforeStale(t *testing.T) {
	opts := testOptions()
	l, err := newMetricList(testShard, 0, opts)
	require.NoError(t, err)
	l.lastFlushedNanos = 1234
	l.flushBefore(1000, discardType)
	require.Equal(t, int64(1234), l.LastFlushedNanos())
}

func TestMetricLists(t *testing.T) {
	lists := newMetricLists(testShard, testOptions())
	require.False(t, lists.closed)

	// Create a new list.
	l, err := lists.FindOrCreate(time.Second)
	require.NoError(t, err)
	require.NotNil(t, l)
	require.Equal(t, 1, lists.Len())

	// Find the same list.
	l2, err := lists.FindOrCreate(time.Second)
	require.NoError(t, err)
	require.Equal(t, l, l2)
	require.Equal(t, 1, lists.Len())

	// Finding or creating in a closed list should result in an error.
	lists.Close()
	_, err = lists.FindOrCreate(time.Second)
	require.Equal(t, errListsClosed, err)
	require.True(t, lists.closed)

	// Closing a second time should have no impact.
	lists.Close()
	require.True(t, lists.closed)
}

type testElemPair struct {
	elem   metricElem
	metric unaggregated.MetricUnion
}

func validateFlushed(
	t *testing.T,
	expected []testAggMetric,
	flushed []aggregated.ChunkedMetricWithStoragePolicy,
) {
	require.Equal(t, len(expected), len(flushed))
	for i := 0; i < len(flushed); i++ {
		require.Equal(t, expected[i].idPrefix, flushed[i].ChunkedID.Prefix)
		require.Equal(t, []byte(expected[i].id), flushed[i].ChunkedID.Data)
		require.Equal(t, expected[i].idSuffix, flushed[i].ChunkedID.Suffix)
		require.Equal(t, expected[i].timeNanos, flushed[i].TimeNanos)
		require.Equal(t, expected[i].value, flushed[i].Value)
		require.Equal(t, expected[i].sp, flushed[i].StoragePolicy)
	}
}

type flushBeforeResult struct {
	beforeNanos int64
	flushType   flushType
}

type writeFn func(mp aggregated.ChunkedMetricWithStoragePolicy) error
type writeflushFn func() error

type mockWriter struct {
	writeFn writeFn
	flushFn writeflushFn
}

func (w *mockWriter) Write(mp aggregated.ChunkedMetricWithStoragePolicy) error {
	return w.writeFn(mp)
}

func (w *mockWriter) Flush() error { return w.flushFn() }
func (w *mockWriter) Close() error { return nil }

type newWriterFn func(scope tally.Scope) (Writer, error)

type mockHandler struct{ newWriterFn newWriterFn }

func (h *mockHandler) NewWriter(scope tally.Scope) (Writer, error) { return h.newWriterFn(scope) }
func (h *mockHandler) Close()                                      {}