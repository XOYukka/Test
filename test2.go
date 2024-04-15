package main

import (
	"fmt"
	"time"
)

type Task struct {
	id          int
	createdTime time.Time
	result      string
}

func main() {
	taskCreator := func(taskChan chan Task) {
		defer close(taskChan)
		for {
			time.Sleep(time.Second)
			createdTime := time.Now()
			task := Task{id: int(createdTime.Unix()), createdTime: createdTime}
			taskChan <- task
		}
	}

	taskProcessor := func(task Task, doneTasks, undoneTasks chan Task) {
		if time.Now().Nanosecond()%2 > 0 {
			task.result = "error"
			undoneTasks <- task
		} else {
			task.result = "success"
			doneTasks <- task
		}
	}

	taskChan := make(chan Task)
	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)

	go taskCreator(taskChan)

	for i := 0; i < 10; i++ {
		go func() {
			for task := range taskChan {
				taskProcessor(task, doneTasks, undoneTasks)
			}
		}()
	}

	done := make(map[int]Task)
	undone := make(map[int]Task)

	go func() {
		for t := range doneTasks {
			done[t.id] = t
		}
	}()

	go func() {
		for t := range undoneTasks {
			undone[t.id] = t
		}
	}()

	time.Sleep(5 * time.Second)

	fmt.Println("Done tasks:")
	for _, t := range done {
		fmt.Printf("Task ID: %d, Created Time: %s, Result: %s\n", t.id, t.createdTime.Format(time.RFC3339), t.result)
	}

	fmt.Println("\nUndone tasks:")
	for _, t := range undone {
		fmt.Printf("Task ID: %d, Created Time: %s, Result: %s\n", t.id, t.createdTime.Format(time.RFC3339), t.result)
	}
}
