// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compliance

import (
	"math"
	"strings"
	"testing"
)

func checkExpectHas(s IntervalSet, expectHas []int, t *testing.T) {
	for _, i := range expectHas {
		if !s.Has(i) {
			t.Errorf("actual.Has(%d): got false, want true", i)
		}
	}
}

func checkExpectHasNot(s IntervalSet, expectHasNot []int, t *testing.T) {
	for _, i := range expectHasNot {
		if s.Has(i) {
			t.Errorf("actual.Has(%d): got true, want false", i)
		}
	}
}

func checkExpectedSlice(s IntervalSet, expectSlice []int, t *testing.T) {
	actualSlice := make([]int, 0)
	actualSlice = s.AppendTo(actualSlice)

	t.Logf("actualSlice: %v", actualSlice)
	t.Logf("expectSlice: %v", expectSlice)

	if s.Len() != len(expectSlice) {
		t.Errorf("actual.Len(): got %d, want %d", s.Len(), len(expectSlice))
	}

	if len(expectSlice) != len(actualSlice) {
		t.Errorf("len(actualSlice) got %d, want %d", len(actualSlice), len(expectSlice))
	} else {
		for i := 0; i < len(actualSlice); i++ {
			if expectSlice[i] != actualSlice[i] {
				t.Errorf("actualSlice[%d]: got %d, want %d", i, actualSlice[i], expectSlice[i])
				break
			}
		}
	}
	var other IntervalSet
	s.Copy(&other)
	i := 0
	var elem int
	for other.TakeMin(&elem) {
		if i >= len(expectSlice) {
			t.Errorf("TakeMin(&elem %d): got %d, want nothing", i, elem)
			continue
		}
		if expectSlice[i] != elem {
			t.Errorf("TakeMin(&elem %d): got %d, want %d", i, elem, expectSlice[i])
		}
		i++
	}
	if i != len(expectSlice) {
		t.Errorf("TakeMin ended early: got %d elements, want %d elements", i, len(expectSlice))
	}
}

func checkExpectedRanges(s IntervalSet, expectRanges []interval, t *testing.T) {

	t.Logf("actual set: %s", s.String())
	t.Logf("expectRanges: %v", expectRanges)

	i := 0
	var p *block
	for p = s.head(); p != nil && i < len(expectRanges); p = s.next(p) {
		if p.r != expectRanges[i] {
			t.Errorf("actual range[i]: got %s, want %s", i, p.r, expectRanges[i])
			return
		}
		i++
	}
	if i < len(expectRanges) {
		t.Errorf("len(actual ranges): got %d, want %d", i, len(expectRanges))
	}
	if p != nil {
		for p1 := p; p1 != nil; p1 = s.next(p1) {
			i++
		}
		t.Errorf("len(actual ranges): got %d with extras starting at %s, want %d", i, p.r.String(), len(expectRanges))
	}
}

type expectedRanges []interval

func TestIntSet(t *testing.T) {
	tests := []struct {
		name         string
		inserts      []int
		ranges       []interval
		differences  []interval
		remove       []int
		removeranges []interval
		expectHas    []int
		expectHasNot []int
		expectSlice  []int
		expectRanges []interval
		expectErr    string
	}{
		{
			name:         "empty",
			expectHasNot: []int{0},
			expectSlice:  []int{},
			expectRanges: []interval{},
		},
		{
			name:         "emptyremoveone",
			remove: []int{2, 2, 2, 2},
			expectHasNot: []int{0},
			expectSlice:  []int{},
			expectRanges: []interval{},
		},
		{
			name:         "emptyremoverange",
			removeranges: []interval{{1, 10}},
			expectHasNot: []int{0},
			expectSlice:  []int{},
			expectRanges: []interval{},
		},
		{
			name: "singlesingle",
			inserts: []int{1},
			expectHas: []int{1},
			expectHasNot: []int{0, 2},
			expectSlice:  []int{1},
			expectRanges: []interval{{1, 1}},
		},
		{
			name: "singlesingle",
			inserts: []int{2, 2, 2, 2},
			expectSlice:  []int{2},
			expectRanges: []interval{{2, 2}},
		},
		{
			name: "disjointsingles",
			inserts: []int{1, 3},
			expectSlice:  []int{1, 3},
			expectRanges: []interval{{1, 1}, {3, 3}},
		},
		{
			name: "adjacentsingles",
			inserts: []int{1, 2},
			expectSlice:  []int{1, 2},
			expectRanges: []interval{{1, 2}},
		},
		{
			name: "adjacentsinglesremove1st",
			inserts: []int{1, 2},
			remove: []int{1},
			expectSlice:  []int{2},
			expectRanges: []interval{{2, 2}},
		},
		{
			name: "adjacentsinglesremovelast",
			inserts: []int{1, 2},
			remove: []int{2},
			expectSlice:  []int{1},
			expectRanges: []interval{{1, 1}},
		},
		{
			name: "contiguoussingles",
			inserts: []int{1, 3, 5, 2, 4},
			expectSlice:  []int{1, 2, 3, 4, 5},
			expectRanges: []interval{{1, 5}},
		},
		{
			name: "singlerange",
			ranges: []interval{{1, 10}},
			expectRanges: []interval{{1, 10}},
		},
		{
			name: "disjointranges",
			ranges: []interval{{5, 7}, {2, 3}},
			expectHas:    []int{2, 6},
			expectHasNot: []int{0, 4, 8},
			expectRanges: []interval{{2, 3}, {5, 7}},
		},
		{
			name: "adjacentranges",
			ranges: []interval{{-1, 200}, {200, 211}},
			expectRanges: []interval{{-1, 211}},
		},
		{
			name: "contiguousranges",
			ranges: []interval{{1, 3}, {5, 7}, {3, 5}},
			expectRanges: []interval{{1, 7}},
		},
		{
			name: "overlappingranges",
			ranges: []interval{{1, 3}, {1, 3}, {2, 5}},
			expectRanges: []interval{{1, 5}},
		},
		{
			name: "mix",
			inserts: []int{2, 4},
			ranges: []interval{{1, 3}, {5, 8}},
			expectRanges: []interval{{1, 8}},
		},
		{
			name: "diffnooverlap",
			ranges: []interval{{1, 100}},
			differences: []interval{{150, 200}},
			expectRanges: []interval{{1, 100}},
		},
		{
			name: "diffmiddle",
			ranges: []interval{{1, 100}},
			differences: []interval{{40, 50}},
			expectRanges: []interval{{1, 39}, {51, 100}},
		},
		{
			name: "difffirst",
			ranges: []interval{{1, 100}},
			differences: []interval{{1, 2}},
			expectRanges: []interval{{3, 100}},
		},
		{
			name: "difflast",
			ranges: []interval{{1, 100}},
			differences: []interval{{90, 100}},
			expectRanges: []interval{{1, 89}},
		},
		{
			name: "diffleft",
			ranges: []interval{{1, 100}},
			differences: []interval{{-10, 10}},
			expectRanges: []interval{{11, 100}},
		},
		{
			name: "diffright",
			ranges: []interval{{1, 100}},
			differences: []interval{{95, 105}},
			expectRanges: []interval{{1, 94}},
		},
		{
			name: "diffequals",
			ranges: []interval{{1, 100}},
			differences: []interval{{1, 100}},
			expectRanges: []interval{},
		},
		{
			name: "diffcontains",
			ranges: []interval{{1, 100}},
			differences: []interval{{-1, 101}},
			expectRanges: []interval{},
		},
		{
			name: "regress1",
			inserts: []int{12, 0, 1, 3, 13, 5, 6, 14, 8, 10, 11, 15},
			expectRanges: []interval{{0, 1}, {3, 3}, {5, 6}, {8, 8}, {10, 15}},
		},
	}
	for _, tt := range tests {
		populate := func(s *IntervalSet) {	
			for _, i := range tt.inserts {
				s.Insert(i)
			}
			for _, r := range tt.ranges {
				s.InsertInterval(r.start, r.end)
			}
			if len(tt.differences) > 0 {
				var other IntervalSet
				for _, r := range tt.differences {
					other.InsertInterval(r.start, r.end)
				}
				s.DifferenceWith(&other)
			}
			for _, i := range tt.remove {
				s.Remove(i)
			}
			for _, r := range tt.removeranges {
				s.RemoveInterval(r.start, r.end)
			}
		}
		checkError := func(t *testing.T) {
			if err, _ := recover().(error); err != nil {
				if len(tt.expectErr) == 0 {
					t.Errorf("unexpected error: got %w, want no error", err)
				} else if !strings.Contains(err.Error(), tt.expectErr) {
					t.Errorf("unexpected error: got %w, want %q", err, tt.expectErr)
				}
			} else {
				if len(tt.expectErr) != 0 {
					t.Errorf("no expected error: got no error, want %q", tt.expectErr)
				}
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			defer checkError(t)
			var s IntervalSet
			populate(&s)
			t.Logf("checking IntervalSet %s, len=%d", s.String(), s.Len())
			if len(tt.expectRanges) == 0 && len(tt.expectSlice) == 0 {
				if !s.IsEmpty() {
					t.Errorf("actual.IsEmpty(): got false, want true")
				}
			} else if s.IsEmpty() {
				t.Errorf("actual.IsEmpty(): got true, want false")
			}
			if len(tt.expectHas) > 0 {
				checkExpectHas(s, tt.expectHas, t)
			}
			if len(tt.expectHasNot) > 0 {
				checkExpectHasNot(s, tt.expectHasNot, t)
			}
			if len(tt.expectSlice) > 0 {
				checkExpectedSlice(s, tt.expectSlice, t)
			}
			if len(tt.expectRanges) > 0 {
				checkExpectedRanges(s, tt.expectRanges, t)
			}
			t.Logf("before actual.Remove(%d, %d): %s", math.MinInt, math.MaxInt, s.String())
			s.RemoveInterval(math.MinInt, math.MaxInt)
			t.Logf("after actual.Remove(%d, %d): %s", math.MinInt, math.MaxInt, s.String())
			if !s.IsEmpty() {
				t.Errorf("removeAllInts actual.IsEmpty(): got false, want true")
			}
		})
		t.Run(tt.name+"_Copy", func(t *testing.T) {
			defer checkError(t)
			var s0 IntervalSet
			populate(&s0)
			var s IntervalSet
			s0.Copy(&s)
			t.Logf("checking IntervalSet %s.Copy() %s, len=%d", s0.String(), s.String(), s.Len())
			if len(tt.expectRanges) == 0 && len(tt.expectSlice) == 0 {
				if !s.IsEmpty() {
					t.Errorf("actual.IsEmpty(): want true, got false")
				}
			} else if s.IsEmpty() {
				t.Errorf("actual.IsEmpty(): want false, got true")
			}
			if len(tt.expectHas) > 0 {
				checkExpectHas(s, tt.expectHas, t)
			}
			if len(tt.expectHasNot) > 0 {
				checkExpectHasNot(s, tt.expectHasNot, t)
			}
			if len(tt.expectSlice) > 0 {
				checkExpectedSlice(s, tt.expectSlice, t)
			}
			if len(tt.expectRanges) > 0 {
				checkExpectedRanges(s, tt.expectRanges, t)
			}
			t.Logf("before actual.Remove(%d, %d): %s", math.MinInt, math.MaxInt, s.String())
			s.RemoveInterval(math.MinInt, math.MaxInt)
			t.Logf("after actual.Remove(%d, %d): %s", math.MinInt, math.MaxInt, s.String())
			if !s.IsEmpty() {
				t.Errorf("removeAllInts actual.IsEmpty(): got false, want true")
			}
		})
	}
}

func TestIntSet_DifferenceWithEmpty(t *testing.T) {
	tests := []struct {
		name         string
		ranges       []interval
	}{
		{
			name:   "empty",
			ranges: []interval{},
		},
		{
			name: "singlerange",
			ranges: []interval{{1, 10}},
		},
		{
			name: "disjointranges",
			ranges: []interval{{2, 3}, {5, 7}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s IntervalSet
			for _, r := range tt.ranges {
				s.InsertInterval(r.start, r.end)
			}
			t.Logf("checking IntervalSet %s difference {}, len=%d", s.String(), s.Len())
			var other IntervalSet
			s.DifferenceWith(&other)
			checkExpectedRanges(s, tt.ranges, t)
		})
	}
}

func TestIntSet_UnionWith(t *testing.T) {
	tests := []struct {
		name         string
		ranges       []interval
		unionWith    []interval
		expectRanges []interval
	}{
		{
			name:   "empty",
			ranges: []interval{},
			unionWith: []interval{},
			expectRanges: []interval{},
		},
		{
			name:   "emptywithsingle",
			ranges: []interval{},
			unionWith: []interval{{2, 5}},
			expectRanges: []interval{{2, 5}},
		},
		{
			name:   "emptywithdisjoint",
			ranges: []interval{},
			unionWith: []interval{{4, 6}, {1, 2}},
			expectRanges: []interval{{1, 2}, {4, 6}},
		},
		{
			name: "singlerangewithempty",
			ranges: []interval{{1, 10}},
			unionWith: []interval{},
			expectRanges: []interval{{1, 10}},
		},
		{
			name: "singlerangecontains",
			ranges: []interval{{1, 10}},
			unionWith: []interval{{0, 100}},
			expectRanges: []interval{{0, 100}},
		},
		{
			name: "singlerangecontained",
			ranges: []interval{{1, 10}},
			unionWith: []interval{{2, 5}},
			expectRanges: []interval{{1, 10}},
		},
		{
			name: "singlerangewithsingledisjointbefore",
			ranges: []interval{{1, 10}},
			unionWith: []interval{{-5, -1}},
			expectRanges: []interval{{-5, -1}, {1, 10}},
		},
		{
			name: "singlerangewithsingledisjointafter",
			ranges: []interval{{1, 10}},
			unionWith: []interval{{15, 21}},
			expectRanges: []interval{{1, 10}, {15, 21}},
		},
		{
			name: "singlerangewithadjacent",
			ranges: []interval{{4, 7}},
			unionWith: []interval{{1, 3}, {8, 11}},
			expectRanges: []interval{{1, 11}},
		},
		{
			name: "disjointranges",
			ranges: []interval{{2, 3}, {5, 6}, {10, 15}},
			unionWith: []interval{{2, 5}, {8, 8}},
			expectRanges: []interval{{2, 6}, {8, 8}, {10, 15}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("with empty", func(t *testing.T) {
				var s IntervalSet
				for _, r := range tt.ranges {
					s.InsertInterval(r.start, r.end)
				}
				t.Logf("checking IntervalSet %s union {}, len=%d", s.String(), s.Len())
				var other IntervalSet
				s.UnionWith(&other)
				checkExpectedRanges(s, tt.ranges, t)
			})

			var s IntervalSet
			for _, r := range tt.ranges {
				s.InsertInterval(r.start, r.end)
			}
			var other IntervalSet
			for _, r := range tt.unionWith {
				other.InsertInterval(r.start, r.end)
			}
			t.Logf("checking IntervalSet %s union with %s", s.String(), other.String())
			s.UnionWith(&other)
			checkExpectedRanges(s, tt.expectRanges, t)
		})
	}
}

