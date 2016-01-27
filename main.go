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
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/majewsky/todolist/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
)

func main() {
	//parse flags
	port := flag.Int("port", 8080, "serving port")
	flag.Parse()

	//run all requests through the authorization middleware first
	http.HandleFunc("/", authMiddlewareHandler)

	//run server
	fmt.Printf("starting local HTTP server on port %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ListenAndServe: ", err)
	}
}

var authHeaderRx = regexp.MustCompile(`^\s*Basic\s*([a-zA-Z0-9+/=]+)\s*$`)

func authMiddlewareHandler(w http.ResponseWriter, r *http.Request) {
	//check presence of Authorization header
	authData := authHeaderRx.FindStringSubmatch(r.Header.Get("Authorization"))
	if len(authData) == 0 {
		authMissingHandler(w)
		return
	}

	//decode authData into "username:password"
	authDataRaw, err := base64.StdEncoding.DecodeString(authData[1])
	if err != nil {
		authMissingHandler(w)
		return
	}
	authFields := strings.SplitN(string(authDataRaw), ":", 2)
	if len(authFields) < 2 {
		authMissingHandler(w)
		return
	}
	username := authFields[0]
	password := []byte(authFields[1])

	//check username/password
	passwdContents, err := ioutil.ReadFile("todolist-passwd")
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "ReadFile(todolist-passwd): ", err)
		serveError(w, 500, "Failed to read user database.")
		return
	}
	passwdLines := strings.Split(string(passwdContents), "\n")
	for _, line := range passwdLines {
		if !strings.Contains(line, ":") {
			continue
		}
		fields := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if username == fields[0] && bcrypt.CompareHashAndPassword([]byte(fields[1]), password) == nil {
			//auth successful! pass request to the actual handlers, including the username
			r.Header.Set("X-Todo-UserName", username)
			Router.ServeHTTP(w, r)
			return
		}
	}

	//username or password wrong - offer to create this user/password
	hash, _ := bcrypt.GenerateFromPassword(password, 10)
	html := `<section class="wide"><p>To create this user account (or reset its password to the one that you entered, contact the site administrator on a <b style="color:red">secure</b> channel and send the following token to him:</p><pre>`
	html += username + ":" + string(hash)
	html += `</pre></section>`
	serveCommon(w, "Unknown credentials", html, 403)
}

func authMissingHandler(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", `Basic realm="todolist"`)
	w.WriteHeader(401)
}
