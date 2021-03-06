/*******************************************************************************
*
* todolist - tiny single-user todolist app
* Copyright 2016 Stefan Majewsky <majewsky@gmx.net>
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU Affero General Public License as published
* by the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU Affero General Public License for more details.
*
* You should have received a copy of the GNU Affero General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*
*******************************************************************************/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Data struct {
	Milestones []*Milestone
}

type Milestone struct {
	Name  string
	Tasks []*Task
}

type Task struct {
	Done bool
	Text string
}

func dataFileName(userName string) string {
	return "todolist-" + userName + ".txt"
}

//ReadData reads Data from the todolist.txt file. If the file is broken, an
//error is logged and nil is returned.
func ReadData(userName string) *Data {
	contents, err := ioutil.ReadFile(dataFileName(userName))
	if err != nil {
		//missing data file is a valid initial state
		if os.IsNotExist(err) {
			return &Data{}
		}
		fmt.Fprintln(os.Stderr, "ReadData: ", err)
		return nil
	}
	return ParseData(string(contents))
}

func ParseData(text string) *Data {
	headerRx := regexp.MustCompile(`^>\s*(\S.*)$`)
	doneRx := regexp.MustCompile(`^\s*OK\s*(\S.*)$`)
	var data Data
	var milestone *Milestone

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// skip empty lines
		if line == "" {
			continue
		}

		if match := headerRx.FindStringSubmatch(line); len(match) > 0 {
			//line is a milestone header
			milestone = &Milestone{Name: match[1]}
			data.Milestones = append(data.Milestones, milestone)
		} else {
			//otherwise, it's a task
			if milestone == nil {
				//empty milestone
				milestone = &Milestone{Name: ""}
				data.Milestones = append(data.Milestones, milestone)
			}
			doneMatch := doneRx.FindStringSubmatch(line)
			task := &Task{
				Done: len(doneMatch) > 0,
				Text: line,
			}
			if task.Done {
				task.Text = doneMatch[1]
			}
			milestone.Tasks = append(milestone.Tasks, task)
		}
	}

	return &data
}

func (d *Data) String() string {
	var lines []string
	for _, milestone := range d.Milestones {
		if milestone.Name != "" {
			lines = append(lines, fmt.Sprintf("> %s", milestone.Name))
		}
		for _, task := range milestone.Tasks {
			line := task.Text
			if task.Done {
				line = "OK " + line
			}
			lines = append(lines, line)
		}
		//one empty line between milestones for readability
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func (d *Data) WriteData(userName string) (success bool) {
	err := ioutil.WriteFile(dataFileName(userName), []byte(d.String()), 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "WriteData: ", err)
		return false
	}
	return true
}
