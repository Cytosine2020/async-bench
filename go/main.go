package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
    "sync"
)

const CLIENT_COUNT = 1024
const MSG_COUNT = 1024
const BUFFER_SIZE = 16

func main() {
    port := 54321

    var server net.Listener

    for {
        socket, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
        if err != nil {
            if port >= 54421 {
                log.Fatalln(err)
                return
            }

            port += 1
            continue
        }
        server = socket
        break
    }

    defer func() {
        server.Close()        
    }()

    var wg sync.WaitGroup

    for i := 0; i < CLIENT_COUNT; i += 1 {
        wg.Add(1)
        go func(port int) {
            defer func() {
                wg.Done()
            }()

            socket, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
            if err != nil {
                log.Fatalln(err)
                return
            }

            defer func() {
                socket.Close()
            }()

            send := make([]byte, BUFFER_SIZE)
            recv := make([]byte, BUFFER_SIZE)

            for i := 0; i < MSG_COUNT; i += 1 {
                n, err := rand.Read(send)
                if n != BUFFER_SIZE || err != nil {
                    log.Fatalln(err)
                    return
                }

                n, err = socket.Write(send)
                if n != BUFFER_SIZE || err != nil {
                    log.Fatalln(err)
                    return
                }
                n, err = socket.Read(recv)
                if n != BUFFER_SIZE || err != nil {
                    log.Fatalln(err)
                    return
                }

                if !reflect.DeepEqual(send, recv) {
                    log.Fatalln(send, "!=", recv)
                    return
                }
            }
        }(port)
    }

    for i := 0; i < CLIENT_COUNT; i += 1 {
        conn, err := server.Accept()
        if err != nil {
            log.Println("Failed to accept conn.", err)
            return
        }

        wg.Add(1)
        go func(conn net.Conn) {
            defer func() {
                conn.Close()
                wg.Done()
            }()
        
            buffer := make([]byte, BUFFER_SIZE)


            for {
                n, err := conn.Read(buffer)
                if n != BUFFER_SIZE || err != nil {
                    if err == io.EOF {
                        break
                    }

                    log.Fatalln(err)
                    return
                }

                n, err = conn.Write(buffer)
                if n != BUFFER_SIZE || err != nil {
                    log.Fatalln(err)
                    return
                }
            }
        }(conn)
    }

    wg.Wait()
}
