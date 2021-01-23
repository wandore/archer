package submitter

import (
	"archer/common"
	"github.com/coreos/etcd/contrib/recipes"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type Submitter struct {
	taskMap map[string]interface{}
}

func (s *Submitter) Submit(q *recipe.PriorityQueue) {
	taskStr := common.MapToStr(s.taskMap)
	pr := s.taskMap["priority"].(int)
	q.Enqueue(taskStr, uint16(pr))
	log.Println(taskStr + " push into queue")
}

func New(name string) *Submitter {
	submitter := Submitter{taskMap: make(map[string]interface{}, 0)}
	submitter.taskMap["name"] = name
	submitter.taskMap["uuid"] = uuid.NewV4()
	submitter.taskMap["start"] = time.Now().Format("2006-01-02 15:04:05")
	submitter.taskMap["end"] = ""
	submitter.taskMap["subtasks"] = common.GetSubTasks(name)
	submitter.taskMap["topo"] = ""
	submitter.taskMap["status"] = common.TaskPreparing
	submitter.taskMap["priority"] = common.QueueMiddlePriority
	submitter.taskMap["times"] = 0
	return &submitter
}