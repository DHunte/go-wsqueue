package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/fsamin/go-wsqueue"
	"github.com/gorilla/mux"

	"github.com/jmcvetta/randutil"
)

var fServer = flag.Bool("server", false, "Run server")
var fClient1 = flag.Bool("client1", false, "Run client #1")
var fClient2 = flag.Bool("client2", false, "Run client #2")

func main() {
	flag.Parse()
	forever := make(chan bool)

	if *fServer {
		server()
	}
	if *fClient1 {
		client("1")
	}
	if *fClient2 {
		client("2")
	}

	<-forever
}

func server() {
	r := mux.NewRouter()
	s := wsqueue.NewServer(r, "")
	q := s.CreateQueue("queue1", 2)

	http.Handle("/", r)
	go http.ListenAndServe("0.0.0.0:9000", r)

	//Start send message to queue
	go func() {
		for {
			time.Sleep(5 * time.Second)
			s, _ := randutil.AlphaString(10)
			fmt.Println("send")
			q.Send("> message from goroutine 1 : " + s)
		}
	}()

	go func() {
		for {
			time.Sleep(6 * time.Second)
			s, _ := randutil.AlphaString(10)
			fmt.Println("send")
			q.Send("> message from goroutine 2 : " + s)
		}
	}()
}

func client(ID string) {
	//Connect a client
	go func() {
		c := &wsqueue.Client{
			Protocol: "ws",
			Host:     "localhost:9000",
			Route:    "/",
		}
		cMessage, cError, err := c.Listen("queue1")
		if err != nil {
			panic(err)
		}
		for {
			select {
			case m := <-cMessage:
				fmt.Println("\n\n********* Client " + ID + " *********" + m.String() + "\n******************")
			case e := <-cError:
				fmt.Println("\n\n********* Client " + ID + "  *********" + e.Error() + "\n******************")
			}
		}
	}()
}
