package parser

import (
	"testing"

	"github.com/thehivecorporation/raccoon"
	"strings"
)

func TestGeneric_Build(t *testing.T) {
	generic := Generic{}

	t.Run("Testing an example job JSON file", func(t *testing.T){
		f, err := generic.Parse("../examples/exampleJob.json")
		if err != nil {
			t.Fatal(err)
		}

		var job raccoon.JobRequest
		if err := generic.Build(f, &job); err != nil {
			t.Fatal(err)
		}

		if len(*job.TaskList) == 0 {
			t.Error("Task list should not have a length of 0")
		}

		if !strings.Contains((*job.TaskList)[0].Title, "task") {
			t.Error("First task of task list must contain word 'task'")
		}
	})

	t.Run("Testing an example tasks JSON file", func(t *testing.T) {
		f, err := generic.Parse("../examples/exampleTasks.json")
		if err != nil {
			t.Fatal(err)
		}
		var taskList []raccoon.Task
		if err := generic.Build(f, &taskList); err != nil {
			t.Fatal(err)
		}

		if len(taskList) == 0 {
			t.Error("Task list should not have a length of 0")
		}
	})

	t.Run("Testing an example infrastructure JSON file", func(t *testing.T) {
		f, err := generic.Parse("../examples/exampleInfrastructure.json")
		if err != nil {
			t.Fatal(err)
		}
		var infrastructure raccoon.Infrastructure
		if err := generic.Build(f, &infrastructure); err != nil {
			t.Fatal(err)
		}

		if len(infrastructure.Infrastructure) == 0 {
			t.Error("Cluster list should not have a length of 0")
		}
	})
}
