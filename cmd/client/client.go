// -*- coding: utf-8 -*-
// client.go
// -----------------------------------------------------------------------------
//
// Started on <mié 14-06-2023 19:52:44.058030790 (1686765164)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// Description
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"google.golang.org/grpc"

	"github.com/clinaresl/grpcdoer/pkg/todo"
	pb "github.com/clinaresl/grpcdoer/pkg/todoproto"
)

// constants
// ----------------------------------------------------------------------------
const DIAL_PORT = ":50051"
const VERSION string = "0.1.0" // current version
const EXIT_SUCCESS int = 0     // exit with success
const PROMPT string = " > "    // prompt character

// options
var verbose bool // has verbose output been requested?
var version bool // has version info been requested?

// functions
// ----------------------------------------------------------------------------

// initialize the command line parser
func init() {

	// Only the following optional parameters are provided through the
	// command-line interface
	flag.BoolVar(&verbose, "verbose", false, "provides verbose output")
	flag.BoolVar(&version, "version", false, "shows version info and exists")

}

// show current version
func showVersion() {

	// show version
	fmt.Printf("grpcdoer version %s\n", VERSION)
}

// show the help on the interpreter
func showHelp(query string) {

	help := map[string]string{
		"add":     "<task desc> project:<project name> due:<YYYY-MM-DD>: adds a new task",
		"done":    "<task id>: remove the given task",
		"list":    ": show all pending tasks",
		"version": ": show current version",
		"help":    ": shows this help banner",
		"bye":     ": exits to the OS",
	}

	fmt.Println()
	if query == "" {

		// retrieve all messages to shown in the help banner
		lines := make([]string, 0)
		for cmd, cmdhelp := range help {
			lines = append(lines, fmt.Sprintf("\t%v: %v", cmd, cmdhelp))
		}

		// and show them sorted in ascending order
		sort.Strings(lines)
		for _, line := range lines {
			fmt.Println(line)
		}
	} else {

		// show the help line corresponding to the given query
		fmt.Printf("\t%v: %v", query, help[query])
	}
	fmt.Println()
}

func addTask(task todo.AddCommand) error {

	// Connect to the server
	conn, err := grpc.Dial(DIAL_PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(" Connection failure: %v\n", err)
		return err
	}
	defer conn.Close()

	// Create a client
	client := pb.NewTaskServiceClient(conn)

	// Create a context with a timeout equal to one second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a task
	response, err := client.AddTask(ctx,
		&pb.AddTaskRequest{
			Desc:    task.Desc,
			Project: task.Project,
			Due:     task.Due.Format("2006-01-02"),
		})
	if err != nil {
		log.Fatalf(" Server error: %v\n", err)
		return err
	} else {
		log.Printf(" Server response: '%v'\n", response)
	}

	// return success
	return nil
}

func doneTask(id int) error {

	// Connect to the server
	conn, err := grpc.Dial(DIAL_PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(" Connection failure: %v\n", err)
		return err
	}
	defer conn.Close()

	// Create a client
	client := pb.NewTaskServiceClient(conn)

	// Create a context with a timeout equal to one second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a task
	response, err := client.DoneTask(ctx,
		&pb.DoneTaskRequest{
			Id: fmt.Sprintf("%d", id),
		})
	if err != nil {
		log.Fatalf(" Server error: %v\n", err)
		return err
	} else {
		log.Printf(" Server response: '%v'\n", response)
	}

	// return success
	return nil
}

func listTasks() error {

	// Connect to the server
	conn, err := grpc.Dial(DIAL_PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(" Connection failure: %v\n", err)
		return err
	}
	defer conn.Close()

	// Create a client
	client := pb.NewTaskServiceClient(conn)

	// Create a context with a timeout equal to one second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a task
	response, err := client.ListTasks(ctx,
		&pb.ListTasksRequest{})
	if err != nil {
		log.Fatalf(" Server error: %v\n", err)
		return err
	}

	// Get all tasks
	tasks := response.GetTasks()

	// and show them on the standard output
	for _, task := range tasks {
		fmt.Printf(" desc: <%s> project:<%s> due:<%s>\n", task.GetDesc(), task.GetProject(), task.GetDue())
	}

	// return success
	return nil
}

func main() {
	// process the command line interface
	flag.Parse()

	// in case version information was required show it and exit
	if version {
		showVersion()
		os.Exit(EXIT_SUCCESS)
	}

	// run the interpreter!
	for {

		// show the prompt and get the user input
		fmt.Printf(PROMPT)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		response := scanner.Text()

		// process the command in all its args
		if cmd, err := todo.ParseLine(response); err != nil {
			fmt.Println(err)
		} else {

			switch cmd.(type) {

			case todo.ListCommand:
				listTasks()
			case todo.DoneCommand:
				doneTask(cmd.(todo.DoneCommand).Id)
			case todo.AddCommand:
				addTask(todo.AddCommand{
					Desc:    cmd.(todo.AddCommand).Desc,
					Project: cmd.(todo.AddCommand).Project,
					Due:     cmd.(todo.AddCommand).Due,
				})
			case todo.HelpCommand:
				showHelp(cmd.(todo.HelpCommand).Query)
			case todo.ByeCommand:
				os.Exit(EXIT_SUCCESS)
			case todo.VersionCommand:
				showVersion()
			}
		}
	}
}

// Local Variables:
// mode:go
// fill-column:80
// End:
