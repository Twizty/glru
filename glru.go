package glru

import (
  "container/heap"
  "time"
)

type LRUCache struct {
  queue     PriorityQueue
  size      int
  hash      map[interface{}]interface{} // yolo
  calculate func(interface{}) interface{}
  timeout   time.Duration
}

type CalcFunc func(interface{}) interface{}

func NewLRUCache(size int, calcFunc CalcFunc) *LRUCache {
  return &LRUCache{
    queue: make(PriorityQueue, 0, size),
    size: size,
    hash: make(map[interface{}]interface{}),
    calculate: calcFunc,
  }
}

func NewLRUCacheWithTimeout(size int, timeout time.Duration, calcFunc CalcFunc) *LRUCache {
  return &LRUCache{
    queue: make(PriorityQueue, 0, size),
    size: size,
    hash: make(map[interface{}]interface{}),
    calculate: calcFunc,
    timeout: timeout,
  }
}

func (self *LRUCache) Count() int {
  return len(self.queue)
}

func (self *LRUCache) CalculateWithCache(key interface{}) interface{} {
  val, exists := self.hash[key]
  time := time.Now()

  if exists && !self.isExpiredBy(time, self.queue.FindByValue(key)) {
    self.queue.Set(key, time)
    return val
  }

  result := self.calculate(key)
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

func (self *LRUCache) isExpiredBy(t time.Time, item *Item) bool {
  if self.timeout <= 0 { return false }

  return int64(self.timeout) + item.Priority.UnixNano() > t.UnixNano()
}

func (self *LRUCache) CalculateWithoutCache(key interface{}) interface{} {
  _, exists := self.hash[key]
  result := self.calculate(key)
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
