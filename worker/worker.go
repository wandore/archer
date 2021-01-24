package worker

import (
	"archer/common"
	"archer/dag"
	"github.com/coreos/etcd/contrib/recipes"
	"log"
	"strconv"
	"time"
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
	name := w.taskMap["name"].(string)
	log.Println(name + " update task status from " + preStatus.(string) + " to " + status)
}

func (w *Worker) updateTaskTimes() {
	times := int(w.taskMap["times"].(float64)) + 1
	w.taskMap["times"] = times
	name := w.taskMap["name"].(string)
	log.Println(name + " update task times from " + strconv.Itoa(times-1) + " to " + strconv.Itoa(times))
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
	name := w.taskMap["name"].(string)
	log.Println(name + " update task priority from " + strconv.Itoa(prePr) + " to " + strconv.Itoa(pr))
}

func (w *Worker) updateTaskTopo() {
	subTasks := make([]string, 0)
	for _, v := range w.taskMap["subtasks"].([]interface{}) {
		subTasks = append(subTasks, v.(string))
	}

	if len(subTasks) == 1 {
		topo := make([][]string, 0)
		topo = append(topo, subTasks)
		w.taskMap["topo"] = topo
	} else {
		d := dag.New()

		for _, target := range subTasks {
			for _, root := range common.TaskDepsMap[target] {
				d.AddEdge(root, target)
			}
		}

		w.taskMap["topo"] = d.TopoSort()
	}

	name := w.taskMap["name"].(string)
	log.Println(name + " update topo succeed")
}

func (w *Worker) updateSubTasks(overSubTasks []string) {
	subTasks := make([]string, 0)
	for _, v := range w.taskMap["subtasks"].([]interface{}) {
		subTasks = append(subTasks, v.(string))
	}

	for _, t := range overSubTasks {
		for k, v := range subTasks {
			if v == t {
				w.taskMap["subtasks"] = append(subTasks[:k], subTasks[k+1:]...)
				break
			}
		}
	}

	name := w.taskMap["name"].(string)
	log.Println(name + " update subtasks succeed")
}

func (w *Worker) updateEndTime() {
	w.taskMap["end"] = time.Now().Format("2006-01-02 15:04:05")

	name := w.taskMap["name"].(string)
	log.Println(name + " update endtime succeed")
}

func (w *Worker) submit(q *recipe.PriorityQueue) {
	taskStr := common.MapToStr(w.taskMap)
	pr := w.taskMap["priority"].(int)
	q.Enqueue(taskStr, uint16(pr))
	log.Println(taskStr + " push into queue")
}

func (w *Worker) action() ([]string, bool) {
	overSubTasks := make([]string, 0)
	flag := true
	for _, v := range w.taskMap["topo"].([][]string) {
		n := len(v)
		common.WG.Add(n)
		chs := make([]chan int, n)
		for index, subtask := range v {
			chs[index] = make(chan int, 1)
			go common.TaskFuncMap[subtask](chs[index])
		}
		common.WG.Wait()
		for index, ch := range chs {
			returnCode := <-ch
			if returnCode == 0 {
				name := w.taskMap["name"].(string)
				log.Println("subtask of " + name + ": " + v[index] + " succeed")
				overSubTasks = append(overSubTasks, v[index])
			} else {
				name := w.taskMap["name"].(string)
				log.Println("subtask of " + name + ": " + v[index] + " failed")
				flag = false
			}
		}
		if !flag {
			log.Println(w.taskMap["name"].(string) + " failed because subtasks failed")
			return overSubTasks, flag
		}
	}
	log.Println(w.taskMap["name"].(string) + " succeed")
	return overSubTasks, flag
}

func (w *Worker) Work(q *recipe.PriorityQueue) {
	w.taskMap = w.get(q)
	w.updateTaskStatus(common.TaskProcessing)
	w.updateTaskTimes()
	w.updateTaskPriority()
	w.updateTaskTopo()
	overSubTasks, flag := w.action()
	w.updateSubTasks(overSubTasks)
	if flag {
		w.updateTaskStatus(common.TaskSucceed)
		w.updateEndTime()
		taskStr := common.MapToStr(w.taskMap)
		log.Println(taskStr + " succeed")
	} else {
		w.updateTaskStatus(common.TaskFailed)
		w.submit(q)
	}
}

func New() *Worker {
	return &Worker{}
}
