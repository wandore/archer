# archer
Archer is a distributed scheduler system for dag tasks based on etcd.

go run main.go --role submitter --task task3  
You can use such command to start a submitter to produce a task called task3 and send it to etcd queue.  

go run main.go --role worker  
You can use such command to start a worker to consume task from etcd queue.
