package gqueue

import (
	"container/list"
	"reflect"
	"sync"
)

// Queue main struct
type Queue struct {
	sem  sync.Mutex
	list *list.List
}

var tFunc func(val interface{}) bool

// NewQueue new creation func
func NewQueue() *Queue {
	l := list.New()
	return &Queue{list: l}
}

// Size show the queue size
func (q *Queue) Size() int {
	return q.list.Len()
}

// PushQueue insert item to the queue head
func (q *Queue) PushQueue(val interface{}) *list.Element {
	q.sem.Lock()
	defer q.sem.Unlock()
	e := q.list.PushFront(val)
	return e
}

// DeQueue get and remove the latest item of the queue
func (q *Queue) DeQueue() *list.Element {
	q.sem.Lock()
	defer q.sem.Unlock()
	e := q.list.Back()
	q.list.Remove(e)
	return e
}

// DelElement delate value from queue
func (q *Queue) DelElement(val interface{}) interface{} {
	q.sem.Lock()
	defer q.sem.Unlock()
	e := q.list.Front()
	for e != nil {
		if e.Value == val {
			removed := q.list.Remove(e)
			return removed
		}
		e = e.Next()
	}

	return nil
}

// Query query item from queue
func (q *Queue) Query(queryFunc interface{}) *list.Element {
	q.sem.Lock()
	e := q.list.Front()
	defer q.sem.Unlock()
	for e != nil {
		if reflect.TypeOf(queryFunc) == reflect.TypeOf(tFunc) {
			if queryFunc.(func(val interface{}) bool)(e.Value) {
				return e
			}
		} else {
			return nil
		}
		e = e.Next()
	}
	return nil
}

// DelConditionAll delete all the elements which macthc the confition-query func
// The elements deleted whill be pushed into the result list []*list.Element
func (q *Queue) DelConditionAll(queryFunc interface{}) []interface{} {
	var delElements []interface{}
	if reflect.TypeOf(queryFunc) != reflect.TypeOf(tFunc) {
		return delElements
	}

	q.sem.Lock()
	defer q.sem.Unlock()

	e := q.list.Front()
	for e != nil {
		if queryFunc.(func(val interface{}) bool)(e.Value) {
			removedValue := q.list.Remove(e)
			delElements = append(delElements, removedValue)
		}
		e = e.Next()
	}
	return delElements
}

// DelConditionSingle delete an target value from the queue
// which condition is matched the arg queryFunc
func (q *Queue) DelConditionSingle(queryFunc interface{}) interface{} {
	if reflect.TypeOf(queryFunc) != reflect.TypeOf(tFunc) {
		return nil
	}

	q.sem.Lock()
	defer q.sem.Unlock()

	e := q.list.Front()
	for e != nil {
		if queryFunc.(func(val interface{}) bool)(e.Value) {
			removedValue := q.list.Remove(e)
			return removedValue
		}
		e = e.Next()
	}

	return nil
}

// Contain check item exit or not
func (q *Queue) Contain(val interface{}) bool {
	q.sem.Lock()
	defer q.sem.Unlock()
	e := q.list.Front()
	for e != nil {
		if e.Value == val {
			return true
		}
		e = e.Next()
	}
	return false
}
