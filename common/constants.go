package common

const (
	// etcd url
	EtcdUrl = "http://127.0.0.1:2379"

	// queue name
	ProcessingQueue = "ProcessingQueue"
	SucceedQueue = "SucceedQueue"

	// queue priority
	QueueHighPriority = 0
	QueueMiddlePriority = 1
	QueueLowPriority = 2

	// task status
	TaskPreparing = "Preparing"
	TaskProcessing = "Processing"
	TaskSucceed = "Succeed"
	TaskFailed = "Failed"

	// all tasks
	AllTasks = "task0, task1, task2, task3"

	// task deps
	TaskDeps = "task3: task2 #" +
		       "task2: task1, task0"
)






