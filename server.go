package main

import (
  "bufio"
  "fmt"
  "net"
  "time"
  "bytes"
  "strconv"
  "strings"
  "math/rand"
)

const (
  LISTEN_PORT string = ":8080"
  EXIT_VERB string = "exit"
)

func main() {
  rand.Seed(time.Now().UnixNano())
  ln, err := net.Listen("tcp", LISTEN_PORT)
  if !checkError(err, "Server is ready.") {
    panic("Could not create listener.")
  }

  memCache := Cache{make(map[string]Item)}
  clientCommands := make(chan Command)
  defer ln.Close()
  go func(in chan Command) {
    for {
      cmd := <-in
      if cmd.verb == EXIT_VERB {
        cmd.conn.Write([]byte("bye.\r\n"))
        ln.Close()
        break
      } else {
        memCache.HandleCommand(cmd)
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
        header, err := buf.ReadString('\n')
        if err != nil {
          fmt.Printf("Client disconnected: %v\n", err)
          break
        }
        if len(strings.TrimSpace(header)) == 0 {
          fmt.Printf("Client %v requested disconnect.\n", conn)
          conn.Close()
          break
        }
        cmdSlice := strings.Split(strings.Trim(header,"\n\r"), " ")
        if len(cmdSlice) > 1 {
          cmd := Command{verb:cmdSlice[0], key:cmdSlice[1], conn:conn}
          // time to live is optionally specified in seconds
          if len(cmdSlice) > 3 && cmdSlice[len(cmdSlice)-2] == "TTL" {
            ttl, err := strconv.Atoi(cmdSlice[len(cmdSlice)-1])
            if err == nil {
              cmd.ttl = int64(ttl)
            }
          }
          if cmd.RequiresBody() {
            body, ok := readBody(buf)
            if ok {
              cmd.body = body
              clientCommands <- cmd
            } else {
              conn.Write([]byte("-ERR expected command body\r\n"))
            }
          } else {
            clientCommands <- cmd
          }
        } else {
          conn.Write([]byte("-ERR expected input \"CMD KEY[ ...][TTL n]\"\r\n"))
        }
      }
    }(conn)
  }
  fmt.Println("Server exiting now...")
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

// reads all available input line by line
// stops when an empty line is encountered
// the lines are concatenated into body variable
func readBody(reader *bufio.Reader) (body []byte, ok bool) {
  ok = true
  for ok {
    // ReadLine doesn't return the trailing \r\n symbols
    line, isPrefix, err := reader.ReadLine()
    if err != nil {
      fmt.Printf("error: %v\n", err)
      ok = false
    } else {
      if len(line) == 0 && !isPrefix {
        break
      }
      body = bytes.Join([][]byte{body,line},[]byte{13,10})
    }
  }
  return
}
