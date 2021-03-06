package wsqueue

import (
	"fmt"
	"sync"
	"time"
)

//Stack is a thread-safe "First In First Out" stack
type Stack struct {
	top   *stackItem
	count int
	mutex *sync.Mutex
	max   int
}

type stackItem struct {
	data interface{}
	next *stackItem
}

//NewStack intialize a brand new Stack
func NewStack() *Stack {
	s := &Stack{}
	s.mutex = &sync.Mutex{}
	return s
}

//Open the connection to the storage driver
func (s *Stack) Open(o *Options) {
	if o != nil {
		m := o.Storage
		i, b := m["capacity"].(int)
		if !b {
			Logfunc("Error with stack capacity option : %s", i)
			return
		}
		s.max = i
	}
}

// Get peeks at the n-th item in the stack. Unlike other operations, this one costs O(n).
func (s *Stack) Get(index int) (interface{}, error) {
	if index < 0 || index >= s.count {
		return nil, fmt.Errorf("Requested index %d outside stack, length %d", index, s.count)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	n := s.top
	for i := 1; i < s.count-index; i++ {
		n = n.next
	}

	return n.data, nil
}

// Dump prints of the stack.
func (s *Stack) Dump() {
	n := s.top
	fmt.Print("[ ")
	for i := 0; i < s.count; i++ {
		fmt.Printf("%+v ", n.data)
		n = n.next
	}
	fmt.Print("]")
}

//Len returns current length of the stack
func (s *Stack) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.count
}

//Push add an item a the top of the stack
func (s *Stack) Push(item interface{}) {
	if s.max > 0 {
		f := NewFibonacci()
		for s.Len() >= s.max {
			Warnfunc("Stack overflow. Waiting...")
			f.WaitForIt(time.Second)
		}
	}

	n := &stackItem{data: item}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.top == nil {
		s.top = n
	} else {
		n.next = s.top
		s.top = n
	}

	s.count++
}

//Pop returns and removes the botteom of the stack
func (s *Stack) Pop() interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var n *stackItem
	if s.top != nil {
		n = s.top
		s.top = n.next
		s.count--
	}

	if n == nil {
		return nil
	}

	return n.data

}

//Peek returns but doesn't remove the top of the stack
func (s *Stack) Peek() interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	n := s.top
	if n == nil || n.data == nil {
		return nil
	}

	return n.data
}
