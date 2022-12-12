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

func startProject(id int) *Project {
	/**
	* Function for starting a new project
	*
	* Args:
	* 	id (int): The id of the project
	 */
	fmt.Printf("Started project %d\n", id)
	proj := Project{
		Id:       id,
		Progress: 0,
		Done:     false,
	}
	return &proj
}

func focusOn(brain *sync.Mutex, proj *Project) {
	/**
	* Function for focusing on a project
	*
	* Args:
	*	brain (*sync.Mutex): Pointer to the brain mutex
	*	proj (*Project): Pointer to the project
	 */

	// Brain is single core
	brain.Lock()
	defer brain.Unlock()

	// Get work done
	proj.Progress += 10

	// Check when project is considered finished
	projType := reflect.TypeOf((*proj))
	field, _ := projType.FieldByName("Progress")
	fieldTag := field.Tag.Get("max")
	max, _ := strconv.Atoi(fieldTag)
	fmt.Printf("Working on project %d (%d/%d %%)..", proj.Id, proj.Progress, max)

	// Check if project is done
	if proj.Progress == max {
		proj.Done = true
		fmt.Printf(" Done!")
	}
	fmt.Printf("\n")

	// Take a break
	time.Sleep(200 * time.Millisecond)
}

func idealWorld(ideas int) {
	/**
	* Function to simulate how I would work in an ideal world
	*
	* Args:
	*     ideas (int): The number of ideas
	 */

	// Initialize my brain
	var myBrain sync.Mutex

	// For each idea
	for p := 0; p < ideas; p++ {
		// Start a project
		proj := startProject(p)

		// Work on project until it's done
		for proj.Done == false {
			focusOn(&myBrain, proj)
		}
	}
}

func reality(ideas int) {
	/**
	* Function to simulate how I acutally work
	*
	* Args:
	*     ideas (int): The number of ideas
	 */

	// Initialize my brain
	var myBrain sync.Mutex

	// Make space for projects
	projects := make(chan *Project, ideas)

	// Can't think of a metaphore for this one
	var wg sync.WaitGroup

	// For each idea
	for p := 0; p < ideas; p++ {
		// Start a project
		projects <- startProject(p)
		wg.Add(1)

		// Concurrently work on projects when brain is free
		go func() {
			// Make sure to let the world know when the project is done
			defer wg.Done()
			// Get a project from the projects pile
			project := <-projects
			// Work on project until it's done
			for project.Done == false {
				focusOn(&myBrain, project)
			}
		}()
	}
	// Wait until all projects are done
	wg.Wait()
}

func main() {
	// Number of ideas
	ideas := 3

	fmt.Println("What I should be doing!")
	idealWorld(ideas)

	fmt.Println("\nWhat I acutally do..")
	reality(ideas)
}
