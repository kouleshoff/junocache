package main

import (
  "fmt"
  "net"
  "time"
)

type Command struct {
  verb    string
  key     string
  ttl     int64
  body    []byte
  conn    net.Conn
}

func (cmd Command) RequiresBody() (req bool) {
  switch cmd.verb {
  case "SET":
    req = true
  default:
    req = false
  }
  return
}

func (c *Cache) HandleCommand(command Command) {
  var result []byte
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
    c.Items[command.key] = StringItem{expiration, command.body}
    success = true
  case "DEL":
    _, ok := c.Items[command.key]
    delete(c.Items, command.key)
    success = ok
  default:
    success = false
  }
  if success {
    command.conn.Write([]byte("+OK\r\n"))
    if result != nil {
      command.conn.Write(result)
      command.conn.Write([]byte{13,10,13,10})
    }
  } else {
    info = fmt.Sprintf("-ERR %s with key %s\r\n", command.verb, command.key)
    command.conn.Write([]byte(info))
  }
}
