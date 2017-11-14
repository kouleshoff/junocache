package main

import (
  "fmt"
  "net"
  "time"
  "math/rand"
)

type Command struct {
  verb    string
  key     string
  ttl     int64
  conn    net.Conn
}

func (c *Cache) HandleCommand(command Command) {
  var result string
  var success bool
  var info string
  var expiration int64
  if command.ttl != 0 {
    // convert ttl to nanoseconds
    expiration = time.Now().UnixNano() + command.ttl * 1000000000
  }
  switch command.verb {
  case "GET":
    item, ok := c.FindByKey(command.key)
    if ok {
      result, success = item.GetFirst()
    }
  case "SET":
    // TODO get the body from connection Reader
    body := RandStringRunes(16)
    c.Items[command.key] = StringItem{expiration, body}
    success = true
  case "DEL":
    _, ok := c.Items[command.key]
    delete(c.Items, command.key)
    success = ok
  default:
    success = false
  }
  if success {
    info = fmt.Sprintf("+OK %s with key %s %s\n\n", command.verb, command.key, result)
  } else {
    info = fmt.Sprintf("-ERR %s with key %s\n\n", command.verb, command.key)
  }
  command.conn.Write([]byte(info))
}

var LETTERS = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    r := make([]rune, n)
    for i := range r {
        r[i] = LETTERS[rand.Intn(len(LETTERS))]
    }
    return string(r)
}
