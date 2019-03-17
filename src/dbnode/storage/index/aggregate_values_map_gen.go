// Copyright (c) 2019 Uber Technologies, Inc.
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

// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package index

import (
	"github.com/m3db/m3x/ident"
)

// Copyright (c) 2018 Uber Technologies, Inc.
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

// AggregateValuesMapHash is the hash for a given map entry, this is public to support
// iterating over the map using a native Go for loop.
type AggregateValuesMapHash uint64

// AggregateValuesMapHashFn is the hash function to execute when hashing a key.
type AggregateValuesMapHashFn func(ident.ID) AggregateValuesMapHash

// AggregateValuesMapEqualsFn is the equals key function to execute when detecting equality of a key.
type AggregateValuesMapEqualsFn func(ident.ID, ident.ID) bool

// AggregateValuesMapCopyFn is the copy key function to execute when copying the key.
type AggregateValuesMapCopyFn func(ident.ID) ident.ID

// AggregateValuesMapFinalizeFn is the finalize key function to execute when finished with a key.
type AggregateValuesMapFinalizeFn func(ident.ID)

// AggregateValuesMap uses the genny package to provide a generic hash map that can be specialized
// by running the following command from this root of the repository:
// ```
// make hashmap-gen pkg=outpkg key_type=Type value_type=Type out_dir=/tmp
// ```
// Or if you would like to use bytes or ident.ID as keys you can use the
// partially specialized maps to generate your own maps as well:
// ```
// make byteshashmap-gen pkg=outpkg value_type=Type out_dir=/tmp
// make idhashmap-gen pkg=outpkg value_type=Type out_dir=/tmp
// ```
// This will output to stdout the generated source file to use for your map.
// It uses linear probing by incrementing the number of the hash created when
// hashing the identifier if there is a collision.
// AggregateValuesMap is a value type and not an interface to allow for less painful
// upgrades when adding/removing methods, it is not likely to need mocking so
// an interface would not be super useful either.
type AggregateValuesMap struct {
	_AggregateValuesMapOptions

	// lookup uses hash of the identifier for the key and the MapEntry value
	// wraps the value type and the key (used to ensure lookup is correct
	// when dealing with collisions), we use uint64 for the hash partially
	// because lookups of maps with uint64 keys has a fast path for Go.
	lookup map[AggregateValuesMapHash]AggregateValuesMapEntry
}

// _AggregateValuesMapOptions is a set of options used when creating an identifier map, it is kept
// private so that implementers of the generated map can specify their own options
// that partially fulfill these options.
type _AggregateValuesMapOptions struct {
	// hash is the hash function to execute when hashing a key.
	hash AggregateValuesMapHashFn
	// equals is the equals key function to execute when detecting equality.
	equals AggregateValuesMapEqualsFn
	// copy is the copy key function to execute when copying the key.
	copy AggregateValuesMapCopyFn
	// finalize is the finalize key function to execute when finished with a
	// key, this is optional to specify.
	finalize AggregateValuesMapFinalizeFn
	// initialSize is the initial size for the map, use zero to use Go's std map
	// initial size and consequently is optional to specify.
	initialSize int
}

// AggregateValuesMapEntry is an entry in the map, this is public to support iterating
// over the map using a native Go for loop.
type AggregateValuesMapEntry struct {
	// key is used to check equality on lookups to resolve collisions
	key _AggregateValuesMapKey
	// value type stored
	value []ident.ID
}

type _AggregateValuesMapKey struct {
	key      ident.ID
	finalize bool
}

// Key returns the map entry key.
func (e AggregateValuesMapEntry) Key() ident.ID {
	return e.key.key
}

// Value returns the map entry value.
func (e AggregateValuesMapEntry) Value() []ident.ID {
	return e.value
}

// _AggregateValuesMapAlloc is a non-exported function so that when generating the source code
// for the map you can supply a public constructor that sets the correct
// hash, equals, copy, finalize options without users of the map needing to
// implement them themselves.
func _AggregateValuesMapAlloc(opts _AggregateValuesMapOptions) *AggregateValuesMap {
	m := &AggregateValuesMap{_AggregateValuesMapOptions: opts}
	m.Reallocate()
	return m
}

func (m *AggregateValuesMap) newMapKey(k ident.ID, opts _AggregateValuesMapKeyOptions) _AggregateValuesMapKey {
	key := _AggregateValuesMapKey{key: k, finalize: opts.finalizeKey}
	if !opts.copyKey {
		return key
	}

	key.key = m.copy(k)
	return key
}

func (m *AggregateValuesMap) removeMapKey(hash AggregateValuesMapHash, key _AggregateValuesMapKey) {
	delete(m.lookup, hash)
	if key.finalize {
		m.finalize(key.key)
	}
}

// Get returns a value in the map for an identifier if found.
func (m *AggregateValuesMap) Get(k ident.ID) ([]ident.ID, bool) {
	hash := m.hash(k)
	for entry, ok := m.lookup[hash]; ok; entry, ok = m.lookup[hash] {
		if m.equals(entry.key.key, k) {
			return entry.value, true
		}
		// Linear probe to "next" to this entry (really a rehash)
		hash++
	}
	var empty []ident.ID
	return empty, false
}

// Set will set the value for an identifier.
func (m *AggregateValuesMap) Set(k ident.ID, v []ident.ID) {
	m.set(k, v, _AggregateValuesMapKeyOptions{
		copyKey:     true,
		finalizeKey: m.finalize != nil,
	})
}

// AggregateValuesMapSetUnsafeOptions is a set of options to use when setting a value with
// the SetUnsafe method.
type AggregateValuesMapSetUnsafeOptions struct {
	NoCopyKey     bool
	NoFinalizeKey bool
}

// SetUnsafe will set the value for an identifier with unsafe options for how
// the map treats the key.
func (m *AggregateValuesMap) SetUnsafe(k ident.ID, v []ident.ID, opts AggregateValuesMapSetUnsafeOptions) {
	m.set(k, v, _AggregateValuesMapKeyOptions{
		copyKey:     !opts.NoCopyKey,
		finalizeKey: !opts.NoFinalizeKey,
	})
}

type _AggregateValuesMapKeyOptions struct {
	copyKey     bool
	finalizeKey bool
}

func (m *AggregateValuesMap) set(k ident.ID, v []ident.ID, opts _AggregateValuesMapKeyOptions) {
	hash := m.hash(k)
	for entry, ok := m.lookup[hash]; ok; entry, ok = m.lookup[hash] {
		if m.equals(entry.key.key, k) {
			m.lookup[hash] = AggregateValuesMapEntry{
				key:   entry.key,
				value: v,
			}
			return
		}
		// Linear probe to "next" to this entry (really a rehash)
		hash++
	}

	m.lookup[hash] = AggregateValuesMapEntry{
		key:   m.newMapKey(k, opts),
		value: v,
	}
}

// Iter provides the underlying map to allow for using a native Go for loop
// to iterate the map, however callers should only ever read and not write
// the map.
func (m *AggregateValuesMap) Iter() map[AggregateValuesMapHash]AggregateValuesMapEntry {
	return m.lookup
}

// Len returns the number of map entries in the map.
func (m *AggregateValuesMap) Len() int {
	return len(m.lookup)
}

// Contains returns true if value exists for key, false otherwise, it is
// shorthand for a call to Get that doesn't return the value.
func (m *AggregateValuesMap) Contains(k ident.ID) bool {
	_, ok := m.Get(k)
	return ok
}

// Delete will remove a value set in the map for the specified key.
func (m *AggregateValuesMap) Delete(k ident.ID) {
	hash := m.hash(k)
	for entry, ok := m.lookup[hash]; ok; entry, ok = m.lookup[hash] {
		if m.equals(entry.key.key, k) {
			m.removeMapKey(hash, entry.key)
			return
		}
		// Linear probe to "next" to this entry (really a rehash)
		hash++
	}
}

// Reset will reset the map by simply deleting all keys to avoid
// allocating a new map.
func (m *AggregateValuesMap) Reset() {
	for hash, entry := range m.lookup {
		m.removeMapKey(hash, entry.key)
	}
}

// Reallocate will avoid deleting all keys and reallocate a new
// map, this is useful if you believe you have a large map and
// will not need to grow back to a similar size.
func (m *AggregateValuesMap) Reallocate() {
	if m.initialSize > 0 {
		m.lookup = make(map[AggregateValuesMapHash]AggregateValuesMapEntry, m.initialSize)
	} else {
		m.lookup = make(map[AggregateValuesMapHash]AggregateValuesMapEntry)
	}
}