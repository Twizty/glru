package glru

import (
  "time"
)

const (
  CALCULATE_WITH_CACHE_ACTION = 1
  CALCULATE_WITHOUT_CACHE_ACTION = 2
  GET_COUNT = 3
  EXIT = 999
)

type ThreadsafeLRUCache struct {
  c chan *request
}

type request struct {
  action       int
  responsePipe chan<- interface{}
  body         interface{}
}

func handleRequests(cache *LRUCache, c <-chan *request) {
  for request := range c {
    switch request.action {
      case CALCULATE_WITH_CACHE_ACTION:
        request.responsePipe <- cache.CalculateWithCache(request.body)
      case CALCULATE_WITHOUT_CACHE_ACTION:
        request.responsePipe <- cache.CalculateWithoutCache(request.body)
      case GET_COUNT:
        request.responsePipe <- cache.Count()
      case EXIT:
        return
    }
  }
}

func NewThreadsafeLRUCache(size int, chanSize int, calcFunc CalcFunc) *ThreadsafeLRUCache {
  c := make(chan *request, chanSize)

  go func() {
    cache := NewLRUCache(size, calcFunc)

    handleRequests(cache, c)
  }()

  return &ThreadsafeLRUCache{c}
}

func NewThreadsafeLRUCacheWithTimeout(size int, chanSize int, timeout time.Duration, calcFunc CalcFunc) *ThreadsafeLRUCache {
  c := make(chan *request, chanSize)

  go func() {
    cache := NewLRUCacheWithTimeout(size, timeout, calcFunc)

    handleRequests(cache, c)
  }()

  return &ThreadsafeLRUCache{c}
}

func (self *ThreadsafeLRUCache) CalculateWithCache(key interface{}) interface{} {
  resultPipe := make(chan interface{})
  self.c <- &request{CALCULATE_WITH_CACHE_ACTION, resultPipe, key}
  return <-resultPipe
}

func (self *ThreadsafeLRUCache) CalculateWithoutCache(key interface{}) interface{} {
  resultPipe := make(chan interface{})
  self.c <- &request{CALCULATE_WITHOUT_CACHE_ACTION, resultPipe, key}
  return <-resultPipe
}

func (self *ThreadsafeLRUCache) Count() int {
  resultPipe := make(chan interface{})
  self.c <- &request{GET_COUNT, resultPipe, nil}
  return (<-resultPipe).(int)
}

func (self *ThreadsafeLRUCache) Close() error {
  self.c <- &request{EXIT, nil, nil}
  close(self.c)
  return nil
}
