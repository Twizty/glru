package glru

import (
  "time"
)

const (
  CALCULATE_WITH_CACHE_ACTION = 1
  CALCULATE_WITHOUT_CACHE_ACTION = 2
  GET_COUNT = 3
  GET = 4
  EXIT = 999
)

type ThreadsafeLRUCache struct {
  c chan *request
}

type request struct {
  action       int
  responsePipe chan<- interface{}
  key          interface{}
  calculater   Calculater
}

func handleRequests(cache *LRUCache, c <-chan *request) {
  for request := range c {
    switch request.action {
      case CALCULATE_WITH_CACHE_ACTION:
        request.responsePipe <- cache.CalculateWithCache(request.key, request.calculater)
      case CALCULATE_WITHOUT_CACHE_ACTION:
        request.responsePipe <- cache.CalculateWithoutCache(request.key, request.calculater)
      case GET_COUNT:
        request.responsePipe <- cache.Count()
      case GET:
        request.responsePipe <- cache.Get(request.key)
      case EXIT:
        return
    }
  }
}

func NewThreadsafeLRUCache(size int, chanSize int) *ThreadsafeLRUCache {
  c := make(chan *request, chanSize)

  go func() {
    cache := NewLRUCache(size)

    handleRequests(cache, c)
  }()

  return &ThreadsafeLRUCache{c}
}

func NewThreadsafeLRUCacheWithTimeout(size int, chanSize int, timeout time.Duration) *ThreadsafeLRUCache {
  c := make(chan *request, chanSize)

  go func() {
    cache := NewLRUCacheWithTimeout(size, timeout)

    handleRequests(cache, c)
  }()

  return &ThreadsafeLRUCache{c}
}

func (self *ThreadsafeLRUCache) CalculateWithCache(key interface{}, calculater Calculater) interface{} {
  resultPipe := make(chan interface{})
  self.c <- &request{CALCULATE_WITH_CACHE_ACTION, resultPipe, key, calculater}
  return <-resultPipe
}

func (self *ThreadsafeLRUCache) CalculateWithoutCache(key interface{}, calculater Calculater) interface{} {
  resultPipe := make(chan interface{})
  self.c <- &request{CALCULATE_WITHOUT_CACHE_ACTION, resultPipe, key, calculater}
  return <-resultPipe
}

func (self *ThreadsafeLRUCache) Count() int {
  resultPipe := make(chan interface{})
  self.c <- &request{GET_COUNT, resultPipe, nil, nil}
  return (<-resultPipe).(int)
}

func (self *ThreadsafeLRUCache) Close() error {
  self.c <- &request{EXIT, nil, nil, nil}
  close(self.c)
  return nil
}

func (self *ThreadsafeLRUCache) Get(key interface{}) interface{} {
  resultPipe := make(chan interface{})
  self.c <- &request{GET, resultPipe, key, nil}
  return <-resultPipe
}
