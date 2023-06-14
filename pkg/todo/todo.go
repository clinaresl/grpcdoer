// -*- coding: utf-8 -*-
// todo.go
// -----------------------------------------------------------------------------
//
// Started on <dom 28-05-2023 19:50:56.039505947 (1685296256)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// Definition of commands and parsing functions for handling (task warrior) todo
// commands
package todo

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// constants

// the following regexp identify the format of the list and done commands
const reListCommand = `^\s*list\s*$`
const reDoneCommand = `^\s*done\s+(\d+)\s*$`
const reHelpCommand = `^\s*help(\s+(add|bye|done|help|list|version))?\s*$`
const reByeCommand = `^\s*(bye|quit|exit)\s*$`
const reVersionCommand = `^\s*version\s*$`

// The add command matches the command name and the args are processed
// separately
const reAddCommand = `^\s*add\s+(.+)\s*$`

// The following regexps are used to extract the project and due date of an add
// command
const reAddCommandProject = `project:([^\s]+)`
const reAddCommandDue = `due:(\d{4}-\d{2}-\d{2})`

// types
// ----------------------------------------------------------------------------

// Commands are identified as different types

// The list command has no args and thus it is identified just as a number
type ListCommand int

// the done command has only one argument, the id of the task to remove
type DoneCommand struct {
	Id int
}

// Finally, the add command has three different arguments: the description of
// the task, the project it is attached to and also the due date
type AddCommand struct {
	Desc    string
	Project string
	Due     time.Time
}

// The help command can either have no args or the name of the command to get
// information about
type HelpCommand struct {
	Query string
}

// The bye command has no args, and therefore it is defined as an integer
type ByeCommand int

// The version command has no args, and therefore it is defined also as an
// integer
type VersionCommand int

// functions
// ----------------------------------------------------------------------------

// The following function processes the given command line and returns an
// instance of the appropriate command in case any is well recognized.
// Otherwise, an error is returned
func ParseLine(line string) (interface{}, error) {

	// Compile all regexps
	listRegexp := regexp.MustCompile(reListCommand)
	doneRegexp := regexp.MustCompile(reDoneCommand)
	helpRegexp := regexp.MustCompile(reHelpCommand)
	byeRegexp := regexp.MustCompile(reByeCommand)
	versionRegexp := regexp.MustCompile(reVersionCommand)
	addRegexp := regexp.MustCompile(reAddCommand)

	// check whether this is the version command
	if ok := versionRegexp.MatchString(line); ok == true {
		return VersionCommand(0), nil
	}

	// check whether this is the bye command
	if ok := byeRegexp.MatchString(line); ok == true {
		return ByeCommand(0), nil
	}

	// check whether this is the help command
	if ok := helpRegexp.MatchString(line); ok == true {

		// get the command to get help about in case any was given, and return
		// it
		query := helpRegexp.FindStringSubmatch(line)[2]
		return HelpCommand{Query: query}, nil
	}

	// check whether this is the list command
	if ok := listRegexp.MatchString(line); ok == true {
		return ListCommand(0), nil
	}

	// check whether this is the done command
	if ok := doneRegexp.MatchString(line); ok == true {

		// extract the id of the task to remove and return a correct instance of
		// the command "done"
		id, _ := strconv.Atoi(doneRegexp.FindStringSubmatch(line)[1])
		return DoneCommand{
			Id: id,
		}, nil
	}

	// check whether this is the add command
	if ok := addRegexp.MatchString(line); ok == true {

		// search for the project name
		match := regexp.MustCompile(reAddCommandProject).FindStringSubmatchIndex(line)
		if len(match) == 0 {
			return nil, errors.New("Error: missing project")
		}
		projectname := line[match[2]:match[3]]

		// search for the due date
		match = regexp.MustCompile(reAddCommandDue).FindStringSubmatchIndex(line)
		if len(match) == 0 {
			return nil, errors.New("Error: missing due date")
		}
		duedate, _ := time.Parse("2006-01-02", line[match[2]:match[3]])

		// next, get the deescription. First remove all attributes along with
		// their value
		desc := regexp.MustCompile(reAddCommandProject).ReplaceAllString(line, "")
		desc = regexp.MustCompile(reAddCommandDue).ReplaceAllString(desc, "")

		// now, just simply extract the description which is given in the only
		// group remaining in the match of the resulting add command
		txt := strings.Trim(addRegexp.FindStringSubmatch(desc)[1], " ")
		if len(txt) == 0 {
			return nil, errors.New("Error: missing description")
		}

		// and return an instance of the command "add"
		return AddCommand{
			Desc:    txt,
			Project: projectname,
			Due:     duedate,
		}, nil
	}
	return nil, errors.New("Fatal Error: Unknown command. Type 'help' for more information")
}

// Local Variables:
// mode:go
// fill-column:80
// End:
