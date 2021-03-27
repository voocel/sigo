package util

import "sync"

var once sync.Once

// 任务列表
var taskChan chan *TaskExecutor

type TaskFunc func(params ...interface{})

type TaskExecutor struct {
	f TaskFunc
	p []interface{}
}

// 执行
func init() {
	chlist := getTaskChan()
	go func() {
		for t := range chlist {
			t.Exec()
		}
	}()
}

// 单例 初始化任务列表
func getTaskChan() chan *TaskExecutor {
	once.Do(func() {
		taskChan = make(chan *TaskExecutor)
	})
	return taskChan
}

func NewTaskExecutor(f TaskFunc, p []interface{}) *TaskExecutor {
	return &TaskExecutor{f, p}
}

func (t *TaskExecutor) Exec() {
	t.f(t.p)
}

// 加入任务队列
func Task(f TaskFunc, p ...interface{}) {
	getTaskChan() <- NewTaskExecutor(f, p)
}
