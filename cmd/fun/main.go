package main

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cli.Close()
	_, err = cli.KV.Put(context.Background(), "x", "13", clientv3.WithKeysOnly())
	if err != nil {
		log.Fatal(err.Error())

	}
	log.Println("ok")

	res, err := cli.KV.Get(context.Background(), "x")
	if err != nil {
		log.Fatalf("could not get value for key, %s, %s", "x", err.Error())
	}

	log.Println(res.Kvs)

}
