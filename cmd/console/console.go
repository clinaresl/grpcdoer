// -*- coding: utf-8 -*-
// console.go
// -----------------------------------------------------------------------------
//
// Started on <jue 18-05-2023 22:47:37.023869314 (1684442857)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// Description
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/clinaresl/grpcdoer/pkg/todo"
)

// global variables
// ----------------------------------------------------------------------------
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
	fmt.Printf("ex1 version %s\n", VERSION)
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

// main function
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
				fmt.Println("list")
			case todo.DoneCommand:
				fmt.Printf("done <%d>\n", cmd.(todo.DoneCommand).Id)
			case todo.AddCommand:
				fmt.Printf("add desc:<%s> project:<%s> due:<%s>\n", cmd.(todo.AddCommand).Desc, cmd.(todo.AddCommand).Project, cmd.(todo.AddCommand).Due)
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
