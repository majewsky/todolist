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
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

func main() {
	//parse flags
	port := flag.Int("port", 8080, "serving port")
	flag.Parse()

	//setup static file serving (this is only used in development contexts; in
	//production, these should be served by a dedicated HTTP server that also
	//terminates TLS for this application)
	http.HandleFunc("/static/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/style.css")
	})

	//setup the remaining routes with gorilla/mux
	collectRoutes(router)
	http.Handle("/", router)

	//run server
	fmt.Printf("starting local HTTP server on port %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ListenAndServe: ", err)
	}
}
