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

func (s *Submitter) Submit(name string, q *recipe.PriorityQueue) {
	s.taskMap["name"] = name
	s.taskMap["uuid"] = uuid.NewV4()
	s.taskMap["start"] = time.Now().Format("2006-01-02 15:04:05")
	s.taskMap["end"] = ""
	s.taskMap["subtasks"] = common.GetSubTasks(name)
	s.taskMap["topo"] = ""
	s.taskMap["status"] = common.TaskPreparing
	s.taskMap["priority"] = common.QueueMiddlePriority
	s.taskMap["times"] = 0

	taskStr := common.MapToStr(s.taskMap)
	pr := s.taskMap["priority"].(int)
	q.Enqueue(taskStr, uint16(pr))
	log.Println(taskStr + " push into queue")
}

func New() *Submitter {
	return &Submitter{taskMap: make(map[string]interface{}, 0)}
}