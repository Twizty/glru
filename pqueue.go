package glru

import (
  "time"
  "container/heap"
)

type Item struct {
  index    int
  Priority time.Time
  Value    interface{}
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
  return pq[i].Priority.Sub(pq[j].Priority).Nanoseconds() < 0
}

func (pq PriorityQueue) Swap(i, j int) {
  pq[i], pq[j] = pq[j], pq[i]
  pq[i].index = i
  pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
  n := len(*pq)
  item := x.(*Item)
  item.index = n
  *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
  old := *pq
  n := len(old)
  item := old[n-1]
  item.index = -1
  *pq = old[0 : n-1]
  return item
}

func (pq *PriorityQueue) Set(val interface{}, priority time.Time) {
  var item *Item
  for _, e := range (*pq) {
    if e.Value == val {
      item = e
    } 
  }

  if item == nil {
    heap.Push(pq, &Item{
      Value: val,
      Priority: time.Now(),
    })
  } else {
    pq.Update(item, item.Value, priority)
  }
}

func (pq *PriorityQueue) Update(item *Item, value interface{}, priority time.Time) {
  item.Value = value
  item.Priority = priority
  heap.Fix(pq, item.index)
}
