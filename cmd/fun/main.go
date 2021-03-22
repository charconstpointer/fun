package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	v3 "github.com/coreos/etcd/clientv3"
	v3c "github.com/coreos/etcd/clientv3/concurrency"
)

var (
	name = flag.String("name", "unnamed", "name")
)

type Role int

const (
	Leader Role = iota
	Follower
)

type Node struct {
	r Role
	c *v3.Client

	res chan struct{}
}

func NewNode(endpoints ...string) *Node {
	cli, err := v3.New(v3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return &Node{
		c:   cli,
		r:   Follower,
		res: make(chan struct{}),
	}
}

func (n *Node) Accept(ctx context.Context) {
	log.Println("creating new session")
	s, err := v3c.NewSession(n.c)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("creating new election")
	e := v3c.NewElection(s, "dddddd")

	for {
		if err := e.Campaign(ctx, "dddd"); err != nil {
			log.Println(err.Error())
		}
		log.Println("im th leader")
		<-n.res
		if err := e.Resign(ctx); err != nil {
			log.Println(err.Error())
		}
		log.Println("resigned")

	}
}

func (n *Node) resign() {
	n.res <- struct{}{}
}

func main() {
	flag.Parse()
	n := NewNode("127.0.0.1:2379")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	go n.Accept(ctx)

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(99999) * int(time.Millisecond)))
			log.Println("resigning")
			n.resign()
		}
	}()

	<-ctx.Done()

}
