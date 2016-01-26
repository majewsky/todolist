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
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var globalTemplate = template.Must(template.New("global.html").Parse(`<!doctype html>
<html lang="de">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>todolist - {{.Title}}</title>
		<link rel="stylesheet" type="text/css" href="/static/style.css">
	</head>
	<body>
		<header>
			<h1>{{.Title}}</h1>
		</header>
		{{.Content}}
	</body>
</html>`))

func serveHTML(w http.ResponseWriter, title, content string) {
	//render content into globalTemplate
	type vars struct {
		Title   string
		Content template.HTML
	}
	var buf bytes.Buffer
	globalTemplate.Execute(&buf, &vars{title, template.HTML(content)})

	w.Header().Add("Content-Type", "text-html; charset=utf-8")
	n, err := w.Write(buf.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "ResponseWriter.Write failed: ", err)
	}
	if n < buf.Len() {
		fmt.Fprintf(os.Stderr, "ResponseWriter.Write aborted after %d of %d bytes", n, buf.Len())
	}
}
