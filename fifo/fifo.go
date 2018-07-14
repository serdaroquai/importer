package fifo

// FifoQueue is a simple first in first out queue with two methods. Not thread safe
type FifoQueue struct {
	first *element
	last  *element // keep a pointer to last not to iterate from the start each insertion
}

type element struct {
	next  *element
	value interface{}
}

// NewFifoQueue creates a new queue instance
func NewFifoQueue() *FifoQueue {
	return &FifoQueue{}
}

// Get returns the first inserted element, null if empty
func (q *FifoQueue) Get() (interface{}, bool) {
	if q.first == nil { //empty queue
		return nil, false
	} else {
		value := q.first.value
		q.first = q.first.next
		return value, true
	}
}

// Add adds the element to the end of queue (last to be retrieved unless the queue is empty)
func (q *FifoQueue) Add(thing interface{}) *FifoQueue {
	newElement := &element{value: thing}
	if q.first == nil { //empty queue
		q.first = newElement
	} else {
		q.last.next = newElement // update last element to point to new one
	}
	q.last = newElement
	return q
}
