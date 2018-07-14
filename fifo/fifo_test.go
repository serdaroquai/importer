package fifo

import (
	"testing"
)

func TestAdd(t *testing.T) {
	q := NewFifoQueue()
	if q.Get() != nil {
		t.Fail()
	}
	q.Add("1").Add("2")
	if q.Get() != "1" {
		t.Fail()
	}
	if q.Get() != "2" {
		t.Fail()
	}
	if q.Get() != nil {
		t.Fail()
	}
}