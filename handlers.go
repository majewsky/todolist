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
	"net/http"

	"github.com/gorilla/mux"
)

func collectRoutes(router *mux.Router) {
	router.HandleFunc("/", indexHandler).Methods("GET").Name("index")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: login
	data := ReadData()
	if data == nil {
		serveError(w, "Cannot read data. Check the server log for details.")
	}

	var html string
	if len(data.Milestones) == 0 {
		html = "<p>No tasks defined yet.</p>"
	}

	for _, milestone := range data.Milestones {
		html += "<h2>" + milestone.Name + "</h2>"
		if len(milestone.Tasks) == 0 {
			html += "<p>No tasks in this group.</p>"
		} else {
			html += "<ul>"

			for _, task := range milestone.Tasks {
				html += "<li>" + task.Text + "</li>"
			}
			html += "</ul>"
		}

	}
	serveHTML(w, "Tasks", html)
}
