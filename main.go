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
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	role = flag.String("role", "", "roles: submitter/worker")
	task = flag.String("task", "", "task name: task0/task1/task2/task3")
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

	runtime.GOMAXPROCS(runtime.NumCPU())

	common.Init()

	if *role == "submitter" {
		common.CheckTaskExist(*task)
	}

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
		s := submitter.New()
		s.Submit(*task, q)

		for {
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(3)
			taskList := []string{"task0", "task1", "task2", "task3"}

			s.Submit(taskList[n], q)
			time.Sleep(30 * time.Second)
		}
	} else {
		w := worker.New()
		for {
			w.Work(q)
			time.Sleep(10 * time.Second)
		}
	}
}
