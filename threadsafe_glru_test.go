package glru

import (
  "testing"
  "time"
)

func TestThreadsafeCalculationWithCache(t *testing.T) {
  c := NewThreadsafeLRUCache(10, 2)
  defer c.Close()

  go c.CalculateWithCache(5, baseTestCalculater{5})
  go c.CalculateWithCache(10, baseTestCalculater{10})

  time.Sleep(10 * time.Millisecond)
  if c.Count() != 2 {
    t.Error("Expected action to populate the queue twice, but got", c.Count())
  }
}

func TestThreadsafeGetExistingValue(t *testing.T) {
  c := NewThreadsafeLRUCache(10, 2)
  defer c.Close()

  go c.CalculateWithCache(5, baseTestCalculater{5})

  time.Sleep(5 * time.Millisecond)
  res := c.Get(5)
  if res != 5 {
    t.Error("Expected value to be equal 5")
  }
}

func TestThreadsafeGetNilIfValueIsMissing(t *testing.T) {
  c := NewThreadsafeLRUCache(10, 2)
  defer c.Close()

  res := c.Get(5)
  if res != nil {
    t.Error("Expected value to be nil")
  }
}

func TestThreadsafeCalculationWithCacheOverflow(t *testing.T) {
  c := NewThreadsafeLRUCache(3, 2)
  defer c.Close()

  go c.CalculateWithCache(5, baseTestCalculater{5})
  go c.CalculateWithCache(10, baseTestCalculater{10})
  go c.CalculateWithCache(15, baseTestCalculater{15})
  go c.CalculateWithCache(20, baseTestCalculater{20})

  time.Sleep(10 * time.Millisecond)
  if c.Count() != 3 {
    t.Error("Expected action to populate the queue twice, but got", c.Count())
  }
}
