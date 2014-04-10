package main

import (
  "code.google.com/p/go.net/websocket"  /* websocket.Handler() */
  "net/http"                            /* http.Handle(), http.ListenAndServer() */
  "fmt"                                 /* fmt.Println() */
  "strings"                             /* trimming and removing space for reqs*/
  "bufio"                               /* bufio.Reader, bufio.Writer funcs */
  "strconv"                             /* convert strings to ints */
  "io"                                  /* io.Copy() */
)

/* a list of users by map[string]=io.Writer */
var users = make(map[string]*bufio.Writer)

func closed_sock(err error) bool {
  if err == io.EOF {
    return true
  } else if err != nil {
    panic("Read: " + err.Error())
    return true
  }
  return false
}

/* parse the command */
func writeToBuf(message string, w *bufio.Writer) {
  w.WriteString(message)
  w.Flush()
}

func parseCommand(sock *websocket.Conn) {
  fmt.Println("user connected on page <" + sock.LocalAddr().String() +
              ".html> from <" +  sock.Request().RemoteAddr + ">")

  reader := bufio.NewReader(sock)
  writer := bufio.NewWriter(sock)

  /* varaibles for username and logged_in status */
  var username string
  var logged_in = false

  /* loop until login */
  for {
    /* block until a line is read */
    line, err := reader.ReadString('\n')
    /* if socket closed return from thread */
    if closed_sock(err) {
      return
    }
    /* trim off the command then trim out the spaces, username is left */
    if !strings.HasPrefix(line, "ME IS") {
      writeToBuf("ERROR invalid: must login first\n", writer)
    }
    line = strings.TrimPrefix(line, "ME IS")
    username = strings.TrimSpace(line)

    /* username must not contain spaces
     * username must not be in use
     * login user if valid              */
    if strings.Contains(username, " ") {
      writeToBuf("ERROR username cannot contain spaces\n", writer)
    } else if _, exists := users[username]; exists {
      writeToBuf("ERROR username in use already\n", writer)
    } else {
      users[username] = writer
      logged_in = true
      writeToBuf("OK user logged in\n", writer)
      break /* end the loop on succesful login */
    }
  }
  /* loop until disconnect/logout */
  for logged_in {
    /* read the next line */
    line, err := reader.ReadString('\n')
    if closed_sock(err) {
      break
    }

    /* parse through options;
       trim prefix, get username, process command */
    /* list all users connected */
    if strings.HasPrefix(line, "WHO HERE") {
      line = strings.TrimPrefix(line, "WHO HERE")
      in_name := strings.TrimSpace(line)
      if in_name != username {
        writeToBuf("ERROR name with command mismatch\n", writer)
        continue
      }
      /* add all users to the bufferWriter then flush to send */
      for name, _ := range users {
        writer.WriteString(name + " ")
      }
      writer.Flush()

    /* send to all users */
    } else if strings.HasPrefix(line, "BROADCAST") {
      line = strings.TrimPrefix(line, "BROADCAST")
      in_name := strings.TrimSpace(line)
      if in_name != username {
        writeToBuf("ERROR name with command mismatch\n", writer)
        continue
      }

      /* get size of bytes */
      line, err = reader.ReadString('\n')
      if closed_sock(err) {
        break
      }
      size, err := strconv.Atoi(strings.TrimSpace(line))
      if size == 0 || err != nil {
        writeToBuf("ERROR invalid message size\n", writer)
        continue
      }

      buf := make([]byte, size)
      for i := 0; i < size; i++ {
        buf[i], err = reader.ReadByte()
        if err != nil {
          panic("ReadByte(broadcast)" + err.Error())
        }
      }

      /* write to all users */
      for _, b_out := range users {
        writeToBuf(string(buf), b_out)
      }

    /* logout user */
    } else if strings.HasPrefix(line, "LOGOUT") {
      line = strings.TrimPrefix(line, "LOGOUT")
      in_name := strings.TrimSpace(line)
      if in_name != username {
        writeToBuf("ERROR invalid logout request\n", writer)
        continue
      }
      /* remove from map and close socket */
      delete(users, username)
      writeToBuf("OK user logged out\n", writer)
      logged_in = false

    /* cant log in again */
    } else if strings.HasPrefix(line, "ME IS") {
      writeToBuf("ERROR already logged in on this connection\n", writer)
    /* not a valid command */
    } else {
      writeToBuf("ERROR invalid command\n", writer)
    }
  }
  sock.Close()
  fmt.Println("Closed socket " + sock.Request().RemoteAddr)
}

/* server main func */
func main() {
  /* add files to be handled */
  http.Handle("/", http.FileServer(http.Dir("./html/")))

  /* pass functions which handle the different websockets
    /echo referse to   ws://localhost:8787/echo   not http://localhost:8787/echo */
  http.Handle("/echo", websocket.Handler(Echo))
  http.Handle("/ping", websocket.Handler(ping))
  http.Handle("/chat", websocket.Handler(parseCommand))

  /* listen on port 8787 */
  err := http.ListenAndServe(":8787", nil)

  if err != nil {
      panic("ListenAndServe: " + err.Error())
  }
}

/* echo it back */
func Echo(sock *websocket.Conn) {
  fmt.Println("echo conn on page <" + sock.LocalAddr().String() +
              ".html> from <" +  sock.Request().RemoteAddr + ">")
  io.Copy(sock, sock)
}

/* pings the server with a message; closes socket */
func ping(sock *websocket.Conn) {
  fmt.Println("user connected on page <" + sock.LocalAddr().String() +
              ".html> from <" +  sock.Request().RemoteAddr + ">")
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
  _, err = sock.Write([]byte("OK"))
  if err != nil {
    panic("Read: " + err.Error())
  }
  sock.Close()
}


