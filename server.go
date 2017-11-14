package main

import (
  "bufio"
  "fmt"
  "net"
  "strings"
)

const (
  LISTEN_PORT string = ":8080"
  EXIT_VERB string = "exit"
)

type CacheCommand struct {
  verb    string
  key     string
  body    interface{}
  conn    net.Conn
}

func handleCommand(command CacheCommand) {
  info := fmt.Sprintf("Performing command %s with key %s / body %v\n", command.verb, command.key, command.body)
  command.conn.Write([]byte(info))
}

func checkError(e error, message string) bool {
  if e != nil {
    fmt.Println("Error", e)
    return false
  } else {
    fmt.Println(message)
    return true
  }
}

func main() {
  ln, err := net.Listen("tcp", LISTEN_PORT)
  if !checkError(err, "Server is ready.") {
    panic("Could not create listener.")
  }

  clientCommands := make(chan CacheCommand)
  defer ln.Close()
  go func(in chan CacheCommand) {
    for {
      cmd := <-in
      if cmd.verb == EXIT_VERB {
        cmd.conn.Write([]byte("bye\n"))
        ln.Close()
        break
      } else {
        handleCommand(cmd)
      }
    }
  }(clientCommands)

  for {
    conn, err := ln.Accept()
    if !checkError(err, "Accepted new connection.") {
      break
    }
    go func(conn net.Conn) {
      buf := bufio.NewReader(conn)
      for {
        verb, err := buf.ReadString('\n')
        if err != nil {
          fmt.Printf("Client disconnected: %v\n", err)
          break
        }
        if len(strings.TrimSpace(verb)) == 0 {
          fmt.Printf("Client requested disconnect.\n")
          conn.Close()
          break
        }
        clientCommands <- CacheCommand{strings.Trim(verb,"\n\r"),"key","body",conn}
      }
    }(conn)
  }
  fmt.Println("Server exiting now...")
}
