package main

import (
  "code.google.com/p/go.net/websocket"  /* websocket.Handler() */
  "net/http"                            /* http.Handle(), http.ListenAndServer() */
  "fmt"                                 /* fmt.Println() */
  //"bytes"                               /* buffer, parsing */
  "io"                                  /* io.Copy() */
)

/* echo it back */
func Echo(sock *websocket.Conn) {
  fmt.Println("echo for sock:", sock)
  io.Copy(sock, sock)
}

/* parse the command */
func parseCommand(sock *websocket.Conn) {
  fmt.Println("conn on page <" + sock.LocalAddr().String() +
              ".html> from <" +  sock.RemoteAddr().String() + ">")
  /* array of bytes for a buffer */
  buf := make([]byte, 1024)
  /* read, error check */
  r, err := sock.Read(buf)
  if err != nil {
    panic("Read: " + err.Error())
  }
  /* print size and what was read */
  fmt.Println("Read", r, "bytes")
  fmt.Println(string(buf))
  sock.Write([]byte("OK"))
}

func main() {
  http.Handle("/", http.FileServer(http.Dir(".")))

  http.Handle("/echo", websocket.Handler(Echo))
  http.Handle("/chat", websocket.Handler(parseCommand))

  err := http.ListenAndServe(":8787", nil)

  if err != nil {
      panic("ListenAndServe: " + err.Error())
  }
}
