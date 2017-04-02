package glru

import (
  "testing"
  "time"
)

func TestThreadsafeCalculationWithCache(t *testing.T) {
  c := NewThreadsafeLRUCache(10, 2, func (arg interface{}) interface{} {
    return arg
  })
  defer c.Close()

  go c.CalculateWithCache(5)
  go c.CalculateWithCache(10)

  time.Sleep(10 * time.Millisecond)
  if c.Count() != 2 {
    t.Error("Expected action to populate the queue twice, but got", c.Count())
  }
}

func TestThreadsafeCalculationWithCacheOverflow(t *testing.T) {
  c := NewThreadsafeLRUCache(3, 2, func (arg interface{}) interface{} {
    return arg
  })
  defer c.Close()

  go c.CalculateWithCache(5)
  go c.CalculateWithCache(10)
  go c.CalculateWithCache(15)
  go c.CalculateWithCache(20)

  time.Sleep(10 * time.Millisecond)
  if c.Count() != 3 {
    t.Error("Expected action to populate the queue twice, but got", c.Count())
  }
}
