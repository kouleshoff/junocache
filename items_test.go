package main

import (
  "testing"
)

func TestGetFirst(t *testing.T) {
  item := StringItem{0, "one"}
  data, ok := item.GetFirst()
  if !ok {
    t.Error("GetFirst operation failed on initialized StringItem:", item)
  }
  if data != "one" {
    t.Error("Getting string value which should not exist:", data)
  }
}

func TestSize(t *testing.T) {
  item := StringItem{0, ""}
  s := item.Size()
  if s != 1 {
    t.Error("Size operation failed on initialized StringItem:", item)
  }
}
