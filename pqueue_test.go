package glru

import (
  "testing"
  "time"
  "container/heap"
)

func TestPushPop(t *testing.T) {
  pq := make(PriorityQueue, 0, 3)
  i := &Item{Value: 5, Priority: time.Now()}
  heap.Push(&pq, i)
  heap.Push(&pq, &Item{Value: 3, Priority: time.Now()})
  heap.Push(&pq, &Item{Value: 9, Priority: time.Now()})
  pq.Update(i, 5, time.Now())

  heap.Push(&pq, &Item{Value: 2, Priority: time.Now()})

  item := heap.Pop(&pq).(*Item)
  value := item.Value.(int)
  if value != 3 {
    t.Error("Value must equal 3, but got", value)
  }

  item = heap.Pop(&pq).(*Item)
  value = item.Value.(int)
  if value != 9 {
    t.Error("Value must equal 9, but got", value)
  }

  item = heap.Pop(&pq).(*Item)
  value = item.Value.(int)
  if value != 5 {
    t.Error("Value must equal 5, but got", value)
  }

  item = heap.Pop(&pq).(*Item)
  value = item.Value.(int)
  if value != 2 {
    t.Error("Value must equal 2, but got", value)
  }
}

func TestSet(t *testing.T) {
  pq := make(PriorityQueue, 0, 3)
  heap.Push(&pq, &Item{Value: 5, Priority: time.Now()})
  heap.Push(&pq, &Item{Value: 3, Priority: time.Now()})
  heap.Push(&pq, &Item{Value: 9, Priority: time.Now()})
  pq.Set(5, time.Now())

  item := heap.Pop(&pq).(*Item)
  value := item.Value.(int)
  if value != 3 {
    t.Error("Value must equal 3, but got", value)
  }
}
