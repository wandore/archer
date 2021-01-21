package main


import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/contrib/recipes"
	"log"
	"os"
	"strings"
)


var (
	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	queueName = flag.String("name", "jobs", "queue name")
)


func main() {
	flag.Parse()

	// 解析etcd地址
	endpoints := strings.Split(*addr, ",")

	// 创建etcd的client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 创建/获取队列
	q := recipe.NewQueue(cli, *queueName)

	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		switch action {
		case "push": // 加入队列
			paramMp := make(map[string]string, 0)
			paramMp["name"] = "create_vm"
			paramMp["source"] = "disk"
			paramJson, _ := json.Marshal(paramMp)
			paramStr := string(paramJson)

			q.Enqueue(paramStr) // 入队
		case "pop": // 从队列弹出
			param, err := q.Dequeue() // 出队
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(param) // 输出出队的元素
		case "exit": //退出
			return
		default:
			fmt.Println("unknown action")
		}
	}
}