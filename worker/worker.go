package worker

import (
	"archer/common"
	"archer/dag"
	"github.com/coreos/etcd/contrib/recipes"
	"log"
	"strconv"
)

type Worker struct {
	taskMap map[string]interface{}
}

func (w *Worker) get(q *recipe.PriorityQueue) map[string]interface{} {
	taskStr, err := q.Dequeue()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Get " + taskStr + " from queue")
	return common.StrToMap(taskStr)
}

func (w *Worker) updateTaskStatus(status string) {
	preStatus := w.taskMap["status"]
	w.taskMap["status"] = status
	log.Println("update task status from " + preStatus.(string) + " to " + status)
}

func (w *Worker) updateTaskTimes() {
	times := int(w.taskMap["times"].(float64)) + 1
	w.taskMap["times"] = times
	log.Println("update task times from " + strconv.Itoa(times - 1) + " to " + strconv.Itoa(times))
}

func (w *Worker) updateTaskPriority() {
	var pr int
	if w.taskMap["times"].(int) < 3 {
		pr = common.QueueHighPriority
	} else {
		pr = common.QueueLowPriority
	}
	prePr := int(w.taskMap["priority"].(float64))
	w.taskMap["priority"] = pr
	log.Println("update task priority from " + strconv.Itoa(prePr) + " to " + strconv.Itoa(pr))
}

func (w *Worker) updateTaskTopo() {
	subTasks := make([]string, 0)
	for _, v := range w.taskMap["subtasks"].([]interface{}) {
		subTasks = append(subTasks, v.(string))
	}

	d := dag.New()

	for _, target := range subTasks {
		for _, root := range common.TaskDepsMap[target] {
			d.AddEdge(root, target)
		}
	}

	w.taskMap["topo"] = d.TopoSort()
}

func (w *Worker) action() {
	for _, v := range w.taskMap["topo"].([][]string) {
		n := len(v)
		common.WG.Add(n)
		chs := make([]chan int, n)
		for index, subtask := range v {
			chs[index] = make(chan int, 1)
			go common.TaskFuncMap[subtask](chs[index])
		}
		common.WG.Wait()
		for index, ch := range(chs) {
			returnCode := <- ch
			if returnCode == 1 {
				log.Println(v[index] + " succeed")
			} else {
				log.Println(v[index] + " failed")
			}
		}
	}
}

func (w *Worker) Work(q *recipe.PriorityQueue) {
	w.taskMap = w.get(q)
	w.updateTaskStatus(common.TaskProcessing)
	w.updateTaskTimes()
	w.updateTaskPriority()
	w.updateTaskTopo()
	w.action()
}

func New() *Worker {
	return &Worker{}
}



