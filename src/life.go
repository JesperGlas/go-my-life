package main

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

type Project struct {
	Id       int
	Progress int `max:"100"`
	Done     bool
}

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
	progress_field, _ := project_type.FieldByName("Progress")
	progress_tag := progress_field.Tag.Get("max")
	max, _ := strconv.Atoi(progress_tag)
	fmt.Printf("Working on project %d (%d/%d %%)..", project.Id, project.Progress, max)
	if project.Progress == max {
		project.Done = true
		fmt.Printf(" Done!")
	}
	fmt.Printf("\n")
}

func idealWorld(ideas int) {
	my_focus := sync.Mutex{}
	for p := 0; p < ideas; p++ {
		project := startProject(p)
		for project.Done == false {
			put(&my_focus, project)
		}
	}
}

func main() {
	idealWorld(3)
}
