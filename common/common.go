package common

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var AllTasksList []string
var TaskDepsMap map[string][]string
var SubTasksMap map[string][]string
var TaskFuncMap map[string]func(chan int)

var WG sync.WaitGroup

func MapToStr(taskMap map[string]interface{}) string {
	taskJson, _ := json.Marshal(taskMap)
	taskStr := string(taskJson)

	name := taskMap["name"].(string)
	log.Println(name + " taskMap to taskStr succeed")

	return taskStr
}

func StrToMap(taskStr string) map[string]interface{} {
	var taskMap map[string]interface{}

	err := json.Unmarshal([]byte(taskStr), &taskMap)
	if err != nil {
		log.Fatal(err)
	}

	name := taskMap["name"].(string)
	log.Println(name + " taskStr to taskMap succeed")

	return taskMap
}

func Init()  {
	PreAllTasksList()
	PreTaskDepsMap()
	PreSubTasksMap()
	PreTaskFuncMap()
}

func PreAllTasksList() {
	allTasksTmp := strings.Replace(AllTasks, " ", "", -1)
	AllTasksList = strings.Split(allTasksTmp, ",")
}

func PreTaskDepsMap() {
	TaskDepsMap = make(map[string][]string, 0)
	taskDepsArr := strings.Split(TaskDeps, "#")
	for _, v := range taskDepsArr {
		v = strings.Replace(v, " ", "", -1)
		taskDepsPair := strings.Split(v, ":")
		key := taskDepsPair[0]
		value := strings.Split(taskDepsPair[1], ",")
		TaskDepsMap[key] = value
	}

	log.Println("prepare TaskDepsMap succeed")
}

func PreSubTasksMap() {
	SubTasksMap = make(map[string][]string, 0)
	for _, v := range AllTasksList {
		subTasks := make([]string, 0)
		subTasks = append(subTasks, v)

		i := 0

		for i < len(subTasks)  {
			task := subTasks[i]
			for _, v := range TaskDepsMap[task] {
				exist := false
				for _, t := range subTasks {
					if t == v {
						exist = true
						break
					}
				}
				if !exist {
					subTasks = append(subTasks, v)
				}
			}
			i++
		}

		SubTasksMap[v] = subTasks
	}

	log.Println("prepare SubTasksMap succeed")
}

func PreTaskFuncMap() {
	TaskFuncMap = make(map[string]func(chan int), 0)

	TaskFuncMap["task0"] = Task0
	TaskFuncMap["task1"] = Task1
	TaskFuncMap["task2"] = Task2
	TaskFuncMap["task3"] = Task3

	log.Println("prepare TaskFuncMap succeed")
}

func CheckTaskExist(name string) {
	for _, v := range AllTasksList {
		if v == name {
			log.Println("check task exist succeed")
			return
		}
	}
	log.Fatal(name + " not in AllTasksList")
}

func GetSubTasks(name string) []string {
	_, exist := SubTasksMap[name]
	if !exist {
		log.Fatal(name + " not in SubTasksMap")
	}

	return SubTasksMap[name]
}

func GetTaskFunc(name string) func(chan int) {
	_, exist := TaskFuncMap[name]
	if !exist {
		log.Fatal(name + " not in TaskFuncMap")
	}

	return TaskFuncMap[name]
}

func Task0(ch chan int) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(2)
	log.Println("return_code from Task0 is " + strconv.Itoa(n%2))
	ch <- n%2
	WG.Done()
}

func Task1(ch chan int) {
	log.Println("return_code from Task1 is " + strconv.Itoa(0))
	ch <- 0
	WG.Done()
}

func Task2(ch chan int) {
	log.Println("return_code from Task2 is " + strconv.Itoa(0))
	ch <- 0
	WG.Done()
}

func Task3(ch chan int) {
	log.Println("return_code from Task3 is " + strconv.Itoa(0))
	ch <- 0
	WG.Done()
}




