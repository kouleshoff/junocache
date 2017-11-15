package main

type StringItem struct {
  expiration int64
  Value []byte
}

func (v StringItem) Expiration() int64 {
  return v.expiration
}

func (v StringItem) Size() int {
  return 1
}

func (v StringItem) GetFirst() ([]byte,bool) {
  return v.Value, true
}

func (v StringItem) String() string {
	return string(v.Value)
}
