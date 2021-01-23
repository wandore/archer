package main

import (
	"archer/common"
	"archer/submitter"
	"archer/worker"
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/contrib/recipes"
	"log"
	"os"
)


var (
	role = flag.String("role", "", "roles: submitter or worker")
	task = flag.String("task", "", "task name: task3")
)


func main() {
	flag.Parse()

	if *role == "" {
		fmt.Println("Please input role")
		os.Exit(1)
	} else {
		if *role != "submitter" && *role != "worker" {
			fmt.Println("Roles: submitter or worker")
			os.Exit(1)
		}
	}

	common.Init()
	common.CheckTaskExist(*task)

	endpoints := []string{common.EtcdUrl}

	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connect to etcd succeed")
	defer cli.Close()

	queueName := common.ProcessingQueue

	q := recipe.NewPriorityQueue(cli, queueName)
	log.Println("connect to queue succeed")

	if *role == "submitter" {
		s := submitter.New(*task)
		s.Submit(q)
	} else {
		w := worker.New()
		w.Work(q)
	}
}