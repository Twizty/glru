package glru

import (
  "testing"
  "math/rand"
  "time"
)

func TestCalculationWithCache(t *testing.T) {
  c := NewLRUCache(10, func (arg interface{}) interface{} {
    return arg
  })

  c.CalculateWithCache(5)

  if c.Count() != 1 {
    t.Error("Expected action to populate the queue")
  }

  c.CalculateWithCache(5)

  if c.Count() != 1 {
    t.Error("Expected action not to populate the queue")
  }
}

func TestCalculateWithTimeout(t *testing.T) {
  rand.New(rand.NewSource(time.Now().UnixNano()))
  values := []interface{}{}
  c := NewLRUCacheWithTimeout(10, 100 * time.Nanosecond, func (arg interface{}) interface{} {
    v := getUniqueVal(values)
    values = append(values, arg)
    return v
  })

  firstResult := c.CalculateWithCache(5)
  secondResult := c.CalculateWithCache(5)

  if firstResult != secondResult {
    t.Error("Expected action to cache result")
  }

  time.Sleep(100 * time.Nanosecond)
  thirdResult := c.CalculateWithoutCache(5)
  if firstResult == thirdResult {
    t.Error("Expected cache to be expired")
  }
}

func TestCalculationWithoutCache(t *testing.T) {
  rand.New(rand.NewSource(time.Now().UnixNano()))
  values := []interface{}{}
  c := NewLRUCache(10, func (arg interface{}) interface{} {
    v := getUniqueVal(values)
    values = append(values, arg)
    return v
  })

  firstResult := c.CalculateWithoutCache(5)
  secondResult := c.CalculateWithoutCache(5)

  if firstResult == secondResult {
    t.Error("Expected action not to populate the queue")
  }
}

func getUniqueVal(arr []interface{}) int {
  v := rand.Int()
  for contains(arr, v) {
    v = rand.Int()
  }
  return v
}

func contains(arr []interface{}, val interface{}) bool {
  for _, v := range arr {
    if v == val {
      return true
    }
  }

  return false
}
