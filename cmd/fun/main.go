package main

import (
	"context"
	"log"
	"time"

	v3 "github.com/coreos/etcd/clientv3"
	v3c "github.com/coreos/etcd/clientv3/concurrency"
)

func main() {
	cli, err := v3.New(v3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cli.Close()

	s, err := v3c.NewSession(cli)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("elec")
	e := v3c.NewElection(s, "e")
	ctx := context.Background()
	ctx, canc := context.WithTimeout(ctx, time.Second*3)

	if err := e.Campaign(ctx, "e"); err != nil {
		log.Fatal(err.Error())
	}

	if err := e.Resign(ctx); err != nil {
		log.Fatal(err.Error())
	}

	defer canc()
}
