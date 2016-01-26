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
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func collectRoutes(router *mux.Router) {
	router.HandleFunc("/", indexHandler).Name("index")
	router.HandleFunc("/toggle/{milestone:[0-9]+}/{task:[0-9]+}", toggleHandler).Name("toggle")
	router.HandleFunc("/edit", editHandler).Methods("GET").Name("edit")
	router.HandleFunc("/edit", saveHandler).Methods("POST").Name("save")
	router.HandleFunc("/prune", pruneHandler).Name("prune")
	router.HandleFunc("/backup", pruneHandler).Name("backup")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: login
	data := ReadData()
	if data == nil {
		serveError(w, "Cannot read data. Check the server log for details.")
	}

	html := `<section>`
	hasDone := false

	for mIdx, milestone := range data.Milestones {
		html += `<h2>` + milestone.Name + `</h2>`
		if len(milestone.Tasks) == 0 {
			html += `<p>No tasks in this group.</p>`
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
					template.HTMLEscapeString(task.Text),
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
	serveError(w, "Not implemented")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	serveError(w, "Not implemented")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	serveError(w, "Not implemented")
}

func pruneHandler(w http.ResponseWriter, r *http.Request) {
	serveError(w, "Not implemented")
}

func backupHandler(w http.ResponseWriter, r *http.Request) {
	serveError(w, "Not implemented")
}
