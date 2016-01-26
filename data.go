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
	"errors"
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

//ReadData reads Data from the todolist.txt file. If the file is broken, an
//error is logged and nil is returned.
func ReadData() *Data {
	data, err := readData()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadData: ", err)
		return nil
	}
	return data
}

func readData() (*Data, error) {
	contents, err := ioutil.ReadFile("todolist.txt")
	if err != nil {
		//missing data file is a valid initial state
		if os.IsNotExist(err) {
			return &Data{}, nil
		}
		return nil, err
	}

	headerRx := regexp.MustCompile(`^>\s*(\S.*)$`)
	var data Data
	var milestone *Milestone

	lines := strings.Split(string(contents), "\n")
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
				return nil, errors.New("found a task that is not within a milestone")
			}
			task := &Task{
				Done: strings.HasPrefix(line, "OK "),
				Text: line,
			}
			if task.Done {
				task.Text = strings.TrimPrefix(task.Text, "OK ")
			}
			milestone.Tasks = append(milestone.Tasks, task)
		}
	}

	return &data, nil
}