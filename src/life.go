package main

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// A stuct for project data
type Project struct {
	Id       int
	Progress int `max:"100"`
	Done     bool
}

// Function for starting a new project
func startProject(id int) *Project {
	fmt.Printf("Started project %d\n", id)
	project := Project{
		Id:       id,
		Progress: 0,
		Done:     false,
	}
	return &project
}

func put(focus *sync.Mutex, project *Project) {
	focus.Lock()
	defer focus.Unlock()

	project.Progress += 10
	project_type := reflect.TypeOf((*project))
	field, _ := project_type.FieldByName("Progress")
	field_tag := field.Tag.Get("max")
	max, _ := strconv.Atoi(field_tag)
	fmt.Printf("Working on project %d (%d/%d %%)..", project.Id, project.Progress, max)
	if project.Progress == max {
		project.Done = true
		fmt.Printf(" Done!")
	}
	fmt.Printf("\n")
	time.Sleep(200 * time.Millisecond)
}

func idealWorld(ideas int) {
	var my_focus sync.Mutex
	for p := 0; p < ideas; p++ {
		project := startProject(p)
		for project.Done == false {
			put(&my_focus, project)
		}
	}
}

func reality(ideas int) {
	var my_focus sync.Mutex
	var wg sync.WaitGroup
	backlog := make(chan *Project, ideas)
	for p := 0; p < ideas; p++ {
		backlog <- startProject(p)
		wg.Add(1)
		go func() {
			defer wg.Done()
			project := <-backlog
			for project.Done == false {
				put(&my_focus, project)
			}
		}()
	}
	wg.Wait()
}

func main() {
	ideas := 3
	fmt.Println("In an ideal world!")
	idealWorld(ideas)
	fmt.Println("\nIn Reality..")
	reality(ideas)
}
