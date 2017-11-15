package main

import (
  "testing"
)

func TestGetFirst(t *testing.T) {
  item := StringItem{0, []byte("one")}
  data, ok := item.GetFirst()
  if !ok {
    t.Error("GetFirst operation failed on initialized StringItem:", item)
  }
  if string(data) != "one" {
    t.Error("Getting string value which should not exist:", data)
  }
}

func TestSize(t *testing.T) {
  item := StringItem{0, []byte{32}}
  s := item.Size()
  if s != 1 {
    t.Error("Size operation failed on initialized StringItem:", item)
  }
}
