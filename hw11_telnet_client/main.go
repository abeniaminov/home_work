package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout connection")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatal("hint: go-telnet <host> <port>")
	}

	address := net.JoinHostPort(flag.Args()[0], flag.Args()[1])
	cl := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err := cl.Connect()
	if err != nil {
		log.Fatalf("halt with error: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cl.Send(); err != nil {
					return
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cl.Receive(); err != nil {
					log.Fatalln("Connection to the server is lost")
					return
				}
			}
		}
	}()

	wg.Wait()
}
