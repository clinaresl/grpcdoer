// -*- coding: utf-8 -*-
// server.go
// -----------------------------------------------------------------------------
//
// Started on <mié 14-06-2023 18:34:34.477121253 (1686760474)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// Server side of the gRPC doer service (task warrior TODO list)
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/clinaresl/grpcdoer/pkg/todo"
	pb "github.com/clinaresl/grpcdoer/pkg/todoproto"
)

// constants
// ----------------------------------------------------------------------------
const LISTEN_PORT = ":50051"

// type definitions
// ----------------------------------------------------------------------------
// A task as handled by the server consists almost of the same information given
// in the add task command which is implemented in the todo package. The only
// differences are whether each task has been already completed or not and also
// its id
type Task struct {
	todo.AddCommand
	id     int  // task id
	status bool // whether the task has been completed or not
}

// Define a type for storing information about all tasks. Note it has to
// implement pb.TaskServiceServer as defined in the code generated from the
// protobuf, i.e., the methods:
//
//	AddTask(context.Context, *AddTaskRequest) (*AddTaskResponse, error)
//	DoneTask(context.Context, *DoneTaskRequest) (*DoneTaskResponse, error)
//	ListTasks(context.Context, *ListTasksRequest) (*ListTasksResponse, error)
//
// In addition, we make it an UnimplementedTaskServiceServer
type TaskAgendaService struct {
	tasks []Task // list of tasks
	pb.UnimplementedTaskServiceServer
}

// Methods
// ----------------------------------------------------------------------------

// Add a new task to the todo list or agenda
func (t *TaskAgendaService) AddTask(ctx context.Context, r *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {

	// First, process the due date which is given as a string
	due, err := time.Parse("2006-01-02", r.GetDue())

	// Well, this should never happen as dates should be always automatically
	// verified by the client
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid due date: %v", err)
	}

	// Secondly, compute a new id for the new task
	newid := len(t.tasks)

	// add task to the list of tasks to be done (i.e., status is false)
	t.tasks = append(t.tasks, Task{
		todo.AddCommand{
			Desc:    r.GetDesc(),
			Project: r.GetProject(),
			Due:     due,
		},
		newid,
		false,
	})

	// success - simply return the id of the new task
	return &pb.AddTaskResponse{
		Id: fmt.Sprintf("%d", newid),
	}, nil
}

// mark a task as been completed in the todo list or agenda
func (t *TaskAgendaService) DoneTask(ctx context.Context, r *pb.DoneTaskRequest) (*pb.DoneTaskResponse, error) {

	// get the id of the task to be marked as done
	if id, err := strconv.Atoi(r.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid task id: %v", err)
	} else {

		// check whether the id is valid
		if id < 0 || id >= len(t.tasks) {
			return nil, status.Errorf(codes.InvalidArgument,
				"task id out of bounds: %d", id)
		} else {

			// mark this task as been completed
			t.tasks[id].status = true
		}
	}

	// success
	return &pb.DoneTaskResponse{}, nil
}

func (t *TaskAgendaService) ListTasks(ctx context.Context, r *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {

	// build slice with all the available tasks
	var tasks []*pb.TaskInfo
	for _, v := range t.tasks {

		// build a pb task info
		ti := &pb.TaskInfo{
			Desc:    v.Desc,
			Project: v.Project,
			Due:     v.Due.Format("2006-01-02"),
		}

		// add it to the slice
		tasks = append(tasks, ti)
	}

	// success
	return &pb.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

func main() {

	// Listen
	lis, err := net.Listen("tcp", LISTEN_PORT)
	if err != nil {
		log.Fatalf(" Failed to listen: %v", err)
	}
	log.Println(" Servicing requests on ", LISTEN_PORT)

	// Get a new server
	server := grpc.NewServer()

	// Register the server
	pb.RegisterTaskServiceServer(server, &TaskAgendaService{})

	// Serve
	if err := server.Serve(lis); err != nil {
		log.Fatalf(" Failed to serve: %v", err)
	}
}

// Local Variables:
// mode:go
// fill-column:80
// End:
