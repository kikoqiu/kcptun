package kcp

import "github.com/pkg/errors"

const checkSegmentQueue = false

type SegmentQueue struct {
	container []segment
	front     int
	tail      int
	capacity  int
}

func NewSegmentQueue(cap int) *SegmentQueue {
	loop := &SegmentQueue{
		container: make([]segment, cap, cap),
		front:     0,
		tail:      0,
		capacity:  cap,
	}
	return loop
}

//go:inline
func (q *SegmentQueue) Reset() {
	q.tail = 0
	q.front = 0
}

//go:inline
func (q *SegmentQueue) Len() int {
	return (q.tail + q.capacity - q.front) % q.capacity
}

//go:inline
func (q *SegmentQueue) At(pos int) *segment {
	return &q.container[(q.front+pos)%q.capacity]
}

//go:inline
func (q *SegmentQueue) Cap() int {
	return q.capacity
}

//go:inline
func (q *SegmentQueue) Empty() bool {
	return q.Len() == 0
}

//go:inline
func (q *SegmentQueue) Full() bool {
	return (q.Len() + 1) == q.capacity
}

func (q *SegmentQueue) Extend() {
	buffer := NewSegmentQueue(2 * q.capacity)
	len := q.Len()
	for i := 0; i < len; i++ {
		buffer.container[i] = q.container[(q.front+i)%q.capacity]
	}
	q.container = buffer.container
	q.front = 0
	q.tail = len
	q.capacity = buffer.capacity
}

//go:inline
func (q *SegmentQueue) Enq(elem segment) {
	if q.Full() {
		q.Extend()
	}
	q.container[q.tail] = elem
	q.tail = (q.tail + 1) % q.capacity
}

//go:inline
func (q *SegmentQueue) InsertAt(pos int, elem segment) {
	if q.Full() {
		q.Extend()
	}
	for t := q.Len(); t > pos; t-- {
		q.container[(q.front+t)%q.capacity] = q.container[(q.front+t+q.capacity-1)%q.capacity]
	}
	q.container[(q.front+pos)%q.capacity] = elem
	q.tail = (q.tail + 1) % q.capacity
}

//go:inline
func (q *SegmentQueue) Deq() (segment, error) {
	if checkSegmentQueue && q.Empty() {
		return segment{}, errors.New(
			"failed to dequeue,container is empty.")
	}
	q.Shink()

	ele := q.container[q.front]
	q.front = (q.front + 1) % q.capacity
	return ele, nil
}

func (q *SegmentQueue) Shink() {
	/*if q.capacity<256 && q.Len() <= q.capacity/4 {
	    buffer := NewLoopQueue(q.capacity/2)
	    for i := 0; i < q.Len(); i++ {
	        buffer.container[i] = q.container[(q.front + i) % q.capacity]
	    }
	    q.container = buffer.container
	    q.front = 0
	    q.tail = q.length
	    q.capacity = q.capacity / 2
	}*/
}

//go:inline
func (q *SegmentQueue) DeqN(n int) error {
	if checkSegmentQueue && q.Len() < n {
		return errors.New(
			"failed to dequeue,container is empty.")
	}
	q.front = (q.front + n) % q.capacity
	q.Shink()
	return nil
}
