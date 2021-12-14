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
	"fmt"
	"strings"
)

type block struct {
	next *block
	r interval
}

type interval struct {
	start, end int
}

func (r interval) intersect(other interval) interval {
	if r.end < other.start || r.end < other.start {
		return interval{}
	}
	result := interval{r.start, r.end}
	if result.start < other.start {
		result.start = other.start
	}
	if other.end < result.end {
		result.end = other.end
	}
	return result
}

func (r interval) union(other interval) []interval {
	if r.end < other.start {
		return []interval{r, other}
	}
	if other.end < r.start {
		return []interval{other, r}
	}
	result := interval{r.start, r.end}
	if other.start < result.start {
		result.start = other.start
	}
	if result.end < other.end {
		result.end = other.end
	}
	return []interval{result}
}

func (r interval) isEmpty() bool {
	return r.end < r.start
}

func (r interval) overlaps(other interval) bool {
	return r.start <= other.end && other.start <= r.end
}

func (r interval) contains(other interval) bool {
	return r.start <= other.start && other.end <= r.end
}

func (r interval) includes(x int) bool {
	return r.start <= x && x <= r.end
}

func (r interval) String() string {
	return fmt.Sprintf("[%d, %d]", r.start, r.end)
}

type IntervalSet struct {
	tail *block
	unused *block
}

func (s *IntervalSet) AppendTo(slice []int) []int {
	if s.IsEmpty() {
		return slice
	}
	for p := s.head(); p != nil; p = s.next(p) {
		for i := p.r.start; ; i++ {
			slice = append(slice, i)
			if i == p.r.end {
				break
			}
		}
	}
	return slice
}

func (s *IntervalSet) Clear() {
	if s.IsEmpty() {
		return
	}
	head := s.head()
	s.tail.next = s.unused
	s.unused = head
	s.tail = nil
}

func (s *IntervalSet) head() *block {
	if s.tail == nil {
		return nil
	}
	return s.tail.next
}

func (s *IntervalSet) next(p *block) *block {
	if p == s.tail {
		return nil
	}
	return p.next
}

func (s *IntervalSet) extend(p *block) *block {
	newp := s.newBlock()
	if p != nil {
		newp.next = p.next
		p.next = newp
	} else {
		newp.next = newp
	}

	return newp
}

func (s *IntervalSet) deleteNext(p *block) *block {
	if p == nil || p.next == nil || p.next == p {
		s.Clear()
		return nil
	}
	result := p.next.next
	p.next.next = s.unused
	s.unused = p.next
	p.next = result
	// handle the case where next was the last blck
	if s.tail == s.unused {
		s.tail = p
		return nil
	}
	return p
}

func (s *IntervalSet) newBlock() *block {
	if s.unused == nil {
		return &block{}
	}
	p := s.unused
	s.unused = s.unused.next
	p.next = nil
	return p
}

func (s *IntervalSet) split(sp *block, r interval) {
	if !r.overlaps(sp.r) {
		panic(fmt.Errorf("attempt to split range %s with non-intersecting range %s", sp.r, r))
	}
	if r.contains(sp.r) {
		panic(fmt.Errorf("attempt to split range %s with containing range %s", sp.r, r))
	}
	if r.start <= sp.r.start {
		sp.r.start = r.end + 1
		return
	}
	if sp.r.end <= r.end {
		sp.r.end = r.start - 1
		return
	}
	next := s.extend(sp)
	next.r.end = sp.r.end
	sp.r.end = r.start - 1
	next.r.start = r.end + 1
	if s.tail == sp {
		s.tail = next
	}
}

func (s *IntervalSet) Copy(x *IntervalSet) {
	x.Clear()
	if s.IsEmpty() {
		return
	}

	for p := s.head(); p != nil; p = s.next(p) {
		x.tail = x.extend(x.tail)
		x.tail.r = p.r
	}
}

func (s *IntervalSet) DifferenceWith(x *IntervalSet) {
	if s.IsEmpty() || x.IsEmpty() {
		return
	}

	p := s.tail
	sp := s.head()
	xp := x.head()
	for sp != nil && xp != nil {
		if xp.r.end < sp.r.start {
			xp = x.next(xp)
			continue
		}
		if sp.r.end < xp.r.start {
			p = sp
			sp = s.next(sp)
			continue
		}
		if xp.r.contains(sp.r) {
			p = s.deleteNext(p)
			if p == nil {
				return
			}
			sp = p.next
			continue
		}
		s.split(sp, xp.r)
		p = sp
		sp = s.next(sp)
	}
}

func (s *IntervalSet) Has(x int) bool {
	for p := s.head(); p != nil && p.r.start <= x; p = s.next(p) {
		if p.r.includes(x) {
			return true
		}
	}
	return false
}

func (s *IntervalSet) Insert(x int) bool {
	return s.InsertInterval(x, x)
}

func (s *IntervalSet) mergeNext(p *block) {
	if p == nil || p.next == nil {
		panic(fmt.Errorf("attempt to merge with nil range"))
	}
	if p.next == p {
		panic(fmt.Errorf("attempt to merge with self"))
	}
	if p.next.r.start <= p.r.end || p == s.tail {
		panic(fmt.Errorf("attempt to merge tail with head"))
	}
	p.r.end = p.next.r.end
	s.deleteNext(p)
}

func (s *IntervalSet) InsertInterval(start, end int) bool {
	r := interval{start, end}
	if r.isEmpty() {
		return false
	}
	if s.IsEmpty() {
		s.tail = s.extend(nil)
		s.tail.r = r
		return true
	}

	// special-case tail extensions to avoid going around the loop
	p := s.tail
	if p.r.end < r.start {
		if p.r.end + 1 == r.start {
			p.r.end = r.end
			return true
		}
		s.tail = s.extend(s.tail)
		s.tail.r = r
		return true
	}
	if p.r.includes(r.start) {
		if r.end <= p.r.end {
			return false
		}
		p.r.end = r.end
		return true
	}

	grew := false
	for !r.isEmpty() {
		if p.next == s.tail {
			if p.next.r.end < r.start {
				if p.next.r.end + 1 == r.start {
					p.next.r.end = r.end
					return true
				}
				s.tail = s.extend(s.tail)
				s.tail.r = r
				return true
			}
			if r.end < p.next.r.start {
				if r.end + 1 == p.next.r.start {
					p.next.r.start = r.start
					return true
				}
				s.extend(p)
				p.next.r = r
				return true
			}
			// r.overlaps(p.next.r.end)
			if r.start < p.next.r.start {
				p.next.r.start = r.start
				grew = true
			}
			if p.next.r.end < r.end {
				p.next.r.end = r.end
				grew = true
			}
			return grew
		}
		// p.r.next.end < p.r.next.next.start
		if r.includes(p.next.r.start) {
			p.next.r.start = r.start
			r.start = p.next.r.end + 1
			if r.isEmpty() {
				return true
			}
			grew = true
		}
		if p.next.r.end < r.start {
			if p.next.r.end + 1 != r.start {
				p = p.next
				continue
			}
		}
		if r.end < p.next.r.start {
			if r.end + 1 == p.next.r.start {
				p.next.r.start = r.start
				return true
			}
			s.extend(p)
			p.next.r = r
			return true
		}
		if p.next.r.includes(r.start) {
			if r.end <= p.next.r.end {
				return grew
			}
			// p.next.r.end < r.end !
			r.start = p.next.r.end + 1
			if r.isEmpty() {
				return grew
			}
		}
		// r.start == p.next.r.end + 1 !
		if r.end < p.next.next.r.start {
			if r.end + 1 != p.next.next.r.start {
				p.next.r.end = r.end
				return true
			}
		}
		// r.start == p.next.r.end + 1 && p.next.next.r.start <= r.end + 1
		r.start = p.next.next.r.end
		s.mergeNext(p.next)
		p = s.tail
		grew = true
	}
	return grew
}

func (s *IntervalSet) IsEmpty() bool {
	return s.tail == nil || (s.tail.next == s.tail && s.tail.r.isEmpty())
}

func (s *IntervalSet) IsEqual(other *IntervalSet) bool {
	if s.tail == other.tail {
		return true
	}
	sp := s.head()
	op := other.head()
	for sp != nil && op != nil && sp.r == op.r {
		sp = s.next(sp)
		op = other.next(op)
	}
	return sp == nil && op == nil
}

func (s *IntervalSet) Len() int {
	size := 0
	for p := s.head(); p != nil; p = s.next(p) {
		if p.r.start <= p.r.end {
			size += p.r.end - p.r.start + 1
		}
	}
	return size
}

func (s *IntervalSet) Remove(x int) bool {
	return s.RemoveInterval(x, x)
}

func (s *IntervalSet) RemoveInterval(start, end int) bool {
        if s.IsEmpty() {
		return false
	}

	r := interval{start, end}
	p := s.tail
	shrunk := false
	for p != nil {
		if r.end < p.next.r.start {
			return shrunk
		}
		if p.next.r.end < r.start {
			if p.next == s.tail {
				return shrunk
			}
			p = p.next
			continue
		}
		shrunk = true
		// r.overlaps(p.next.r) !
		if r.start < p.next.r.start {
			r.start = p.next.r.start
		}
		// p.next.r.start <= r.start <= p.next.r.end
		if p.next.r.start == r.start && p.next.r.end <= r.end {
			r.start = p.next.r.end + 1 // can only wrap-around when deleteNext returns nil
			p = s.deleteNext(p)
			if p == nil {
				return true
			}
			continue
		}
		s.split(p.next, r)
	}
	return shrunk
}

func (s *IntervalSet) TakeMin(x *int) bool {
	if s.IsEmpty() {
		return false
	}
	*x = s.tail.next.r.start
	s.tail.next.r.start++
	if s.tail.next.r.isEmpty() {
		s.deleteNext(s.tail)
	}
	return true
}

func (s *IntervalSet) UnionWith(x *IntervalSet) bool {
	if x.IsEmpty() {
		return false
	}
	grew := false
	for xp := x.head(); xp != nil; xp = x.next(xp) {
		if s.InsertInterval(xp.r.start, xp.r.end) {
			grew = true
		}
	}
	return grew
}

func (s *IntervalSet) String() string {
	if s.IsEmpty() {
		return "{}"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "{")
	sep := ""
	for p := s.head(); p != nil; p = s.next(p) {
		fmt.Fprintf(&sb, "%s%s", sep, p.r.String())
		sep = ", "
	}
	fmt.Fprintf(&sb, "}")
	if s.unused != nil {
		fmt.Fprintf(&sb, " (")
		sep = ""
		for p := s.unused; p != nil; p = p.next {
			fmt.Fprintf(&sb, "%s%v", sep, *p)
			sep = ", "
		}
		fmt.Fprintf(&sb, ")")
	}
	return sb.String()
}

func (s *IntervalSet) FindFirst(index *int, matches func(index int) bool) bool {
	for p := s.head(); p != nil; p = s.next(p) {
		i := p.r.start
		for {
			if matches(i) {
				*index = i
				return true
			}
			if i == p.r.end {
				break
			}
			i++
		}
	}
	return false
}

func (s *IntervalSet) VisitAll(visit func(index int)) {
	for p := s.head(); p != nil; p = s.next(p) {
		i := p.r.start
		for {
			visit(i)
			if i == p.r.end {
				break
			}
			i++
		}
	}
}

func (s *IntervalSet) VisitIntervals(visit func(start, end int)) {
	for p := s.head(); p != nil; p = s.next(p) {
		visit(p.r.start, p.r.end)
	}
}
