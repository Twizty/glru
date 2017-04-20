package glru

import (
  "testing"
  "math/rand"
  "time"
)

type testCalculater struct {
  values []interface{}
  arg    interface{}
}

func (tc testCalculater) Calculate() interface{} {
  v := getUniqueVal(tc.values)
  tc.values = append(tc.values, v)
  return v
}

type baseTestCalculater struct {
  arg interface{}
}

func (btc baseTestCalculater) Calculate() interface{} {
  return btc.arg
}

func TestCalculationWithCache(t *testing.T) {
  c := NewLRUCache(10)

  c.CalculateWithCache(5, baseTestCalculater{5})

  if c.Count() != 1 {
    t.Error("Expected action to populate the queue")
  }

  c.CalculateWithCache(5, baseTestCalculater{5})

  if c.Count() != 1 {
    t.Error("Expected action not to populate the queue")
  }
}

func TestGetsValueByExistingKey(t *testing.T) {
  c := NewLRUCache(10)

  c.CalculateWithCache(5, baseTestCalculater{5})

  res := c.Get(5)
  if res != 5 {
    t.Error("Expectd value must be equal 5")
  }
}

func TestReturnsNilIfKeyIsEmpty(t *testing.T) {
  c := NewLRUCache(10)

  res := c.Get(5)
  if res != nil {
    t.Error("Expectd value must be nil")
  }
}

func TestCalculateWithTimeout(t *testing.T) {
  rand.New(rand.NewSource(time.Now().UnixNano()))
  values := make([]interface{}, 0, 0)
  c := NewLRUCacheWithTimeout(10, 100 * time.Nanosecond)

  firstResult := c.CalculateWithCache(5, testCalculater{values, 5})
  secondResult := c.CalculateWithCache(5, testCalculater{values, 5})

  if firstResult != secondResult {
    t.Error("Expected action to cache result")
  }

  time.Sleep(100 * time.Nanosecond)
  thirdResult := c.CalculateWithoutCache(5, testCalculater{values, 5})
  if firstResult == thirdResult {
    t.Error("Expected cache to be expired")
  }
}

func TestCalculationWithoutCache(t *testing.T) {
  rand.New(rand.NewSource(time.Now().UnixNano()))
  values := make([]interface{}, 0, 0)
  c := NewLRUCache(10)

  firstResult := c.CalculateWithoutCache(5, testCalculater{values, 5})
  secondResult := c.CalculateWithoutCache(5, testCalculater{values, 5})

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
