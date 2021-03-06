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
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/majewsky/todolist/Godeps/_workspace/src/github.com/gorilla/mux"
)

var Router = mux.NewRouter()

func init() {
	Router.HandleFunc("/", indexHandler).Name("index")
	Router.HandleFunc("/toggle/{milestone:[0-9]+}/{task:[0-9]+}", toggleHandler).Name("toggle")
	Router.HandleFunc("/edit", editHandler).Methods("GET").Name("edit")
	Router.HandleFunc("/edit", saveHandler).Methods("POST").Name("save")
	Router.HandleFunc("/prune", pruneHandler).Name("prune")
	Router.HandleFunc("/backup", backupHandler).Name("backup")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := ReadData(r.Header.Get("X-Todo-UserName"))
	if data == nil {
		serveError(w, 500, "Cannot read data. Check the server log for details.")
		return
	}

	html := `<section>`
	hasDone := false

	for mIdx, milestone := range data.Milestones {
		if milestone.Name != "" {
			html += `<h2>` + milestone.Name + `</h2>`
		}
		if len(milestone.Tasks) == 0 {
			if milestone.Name != "" {
				html += `<p>No tasks in this group.</p>`
			}
		} else {
			html += `<div class="group">`

			for tIdx, task := range milestone.Tasks {
				class := "open"
				if task.Done {
					class = "done"
					hasDone = true
				}
				path, _ := Router.Get("toggle").URLPath(
					"milestone", strconv.Itoa(mIdx),
					"task", strconv.Itoa(tIdx),
				)
				html += fmt.Sprintf(`<a href="%s" class="%s">%s</a>`,
					path, class,
					HTMLEscapeString(task.Text),
				)
			}
			html += `</div>`
		}
	}

	html += `<div class="table"><div class="row"><a href="/edit" class="action">Edit</a>`
	if hasDone {
		html += `<a href="/prune" class="action">Prune</a>`
	}
	html += `<a href="/backup" class="action">Backup</a></div></div></section>`

	serveHTML(w, "Tasks", html)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	data := ReadData(r.Header.Get("X-Todo-UserName"))
	if data == nil {
		serveError(w, 500, "Cannot read data. Check the server log for details.")
		return
	}

	//retrieve parameters
	vars := mux.Vars(r)
	mIdx, _ := strconv.Atoi(vars["milestone"])
	tIdx, _ := strconv.Atoi(vars["task"])

	if mIdx < 0 || mIdx >= len(data.Milestones) {
		serveError(w, 400, "Milestone index out of range.")
		return
	}
	milestone := data.Milestones[mIdx]
	if tIdx < 0 || tIdx >= len(milestone.Tasks) {
		serveError(w, 400, "Task index out of range.")
		return
	}
	task := milestone.Tasks[tIdx]
	task.Done = !task.Done
	if !data.WriteData(r.Header.Get("X-Todo-UserName")) {
		serveError(w, 500, "Cannot write data. Check the server log for details.")
		return
	}

	w.Header().Add("Location", "/")
	w.WriteHeader(302)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	data := ReadData(r.Header.Get("X-Todo-UserName"))
	if data == nil {
		serveError(w, 500, "Cannot read data. Check the server log for details.")
		return
	}

	html := `<section><p>Write one task per line. If the line starts with the word &quot;OK&quot;, the task is done. To group tasks in a milestone, prefix them with a line starting with &quot;&gt;&quot;.</p><form action="/edit" method="POST"><textarea name="text">`
	html += HTMLEscapeString(data.String())
	html += `</textarea><p><button type="submit">Save</button></p></form></section>`

	serveHTML(w, "Edit Tasks", html)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ParseForm: ", err)
		serveError(w, 400, "Malformed HTTP request.")
		return
	}

	data := ParseData(r.PostFormValue("text"))
	if !data.WriteData(r.Header.Get("X-Todo-UserName")) {
		serveError(w, 500, "Cannot write data. Check the server log for details.")
		return
	}

	w.Header().Add("Location", "/")
	w.WriteHeader(302)
}

func pruneHandler(w http.ResponseWriter, r *http.Request) {
	data := ReadData(r.Header.Get("X-Todo-UserName"))
	if data == nil {
		serveError(w, 500, "Cannot read data. Check the server log for details.")
		return
	}

	//filter tasks that are done
	var openMilestones []*Milestone
	hasDone := false

	for _, milestone := range data.Milestones {
		var openTasks []*Task
		for _, task := range milestone.Tasks {
			if task.Done {
				hasDone = true
			} else {
				openTasks = append(openTasks, task)
			}
		}

		milestone.Tasks = openTasks
		if len(openTasks) > 0 {
			openMilestones = append(openMilestones, milestone)
		}
	}
	data.Milestones = openMilestones

	//write data only if changed
	if hasDone {
		if !data.WriteData(r.Header.Get("X-Todo-UserName")) {
			serveError(w, 500, "Cannot write data. Check the server log for details.")
			return
		}
	}

	w.Header().Add("Location", "/")
	w.WriteHeader(302)
}

func backupHandler(w http.ResponseWriter, r *http.Request) {
	data := ReadData(r.Header.Get("X-Todo-UserName"))
	if data == nil {
		serveError(w, 500, "Cannot read data. Check the server log for details.")
		return
	}

	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	dateStr := time.Now().Format("2006-01-02")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=todolist-%s.txt", dateStr))

	writeWithLogging(w, data.String())
}
