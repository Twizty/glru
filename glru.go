package glru

import (
  "container/heap"
  "time"
)

type Calculater interface {
  Calculate() interface{}
}

type LRUCache struct {
  queue      PriorityQueue
  size       int
  hash       map[interface{}]interface{} // yolo
  timeout    time.Duration
}

func NewLRUCache(size int) *LRUCache {
  return &LRUCache{
    queue: make(PriorityQueue, 0, size),
    size: size,
    hash: make(map[interface{}]interface{}),
  }
}

func NewLRUCacheWithTimeout(size int, timeout time.Duration) *LRUCache {
  return &LRUCache{
    queue: make(PriorityQueue, 0, size),
    size: size,
    hash: make(map[interface{}]interface{}),
    timeout: timeout,
  }
}

func (self *LRUCache) Count() int {
  return len(self.queue)
}

func (self *LRUCache) CalculateWithCache(key interface{}, calculater Calculater) interface{} {
  val, exists := self.hash[key]
  time := time.Now()

  if exists && !self.isExpiredBy(time, self.queue.FindByValue(key)) {
    self.queue.Set(key, time)
    return val
  }

  result := calculater.Calculate()
  if len(self.hash) == self.size {
    item := heap.Pop(&self.queue).(*Item)
    delete(self.hash, item.Value)
  }

  self.hash[key] = result
  heap.Push(&self.queue, &Item{
    Value: key,
    Priority: time,
  })
  return result
}

func (self *LRUCache) Get(key interface{}) interface{} {
  val, exists := self.hash[key]
  time := time.Now()

  if exists && !self.isExpiredBy(time, self.queue.FindByValue(key)) {
    self.queue.Set(key, time)
    return val
  }

  return nil
}

func (self *LRUCache) isExpiredBy(t time.Time, item *Item) bool {
  if self.timeout <= 0 { return false }

  return int64(self.timeout) + item.Priority.UnixNano() > t.UnixNano()
}

func (self *LRUCache) CalculateWithoutCache(key interface{}, calculater Calculater) interface{} {
  _, exists := self.hash[key]
  result := calculater.Calculate()
  time := time.Now()
  self.hash[key] = result

  if exists {
    self.queue.Set(key, time)
  } else {
    if len(self.hash) == self.size {
      item := heap.Pop(&self.queue).(*Item)
      delete(self.hash, item.Value)
    }
    heap.Push(&self.queue, &Item{
      Value: key,
      Priority: time,
    })
  }

  return result
}
