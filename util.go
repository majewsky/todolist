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
	"io"
	"net/http"
	"os"
	"strings"
)

var globalTemplate = `<!doctype html>
<html lang="de">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>todolist - <template:title></title>
		<link rel="stylesheet" type="text/css" href="/static/style.css">
	</head>
	<body>
		<header>
			<h1><template:title></h1>
		</header>
		<template:content>
	</body>
</html>`

func serveHTML(w http.ResponseWriter, title, content string) {
	serveCommon(w, title, content, http.StatusOK)
}

func serveError(w http.ResponseWriter, status int, errorMsg string) {
	serveCommon(w, "Error", errorMsg, status)
}

func serveCommon(w http.ResponseWriter, title, content string, status int) {
	//place content into globalTemplate
	text := strings.Replace(globalTemplate, "<template:title>", HTMLEscapeString(title), -1)
	text = strings.Replace(text, "<template:content>", content, -1)

	//write response header
	w.Header().Add("Content-Type", "text-html; charset=utf-8")
	w.WriteHeader(status)

	//write contents
	writeWithLogging(w, text)
}

func writeWithLogging(w io.Writer, text string) {
	bytes := []byte(text)
	n, err := w.Write(bytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ResponseWriter.Write failed: ", err)
	}
	if n < len(bytes) {
		fmt.Fprintf(os.Stderr, "ResponseWriter.Write aborted after %d of %d bytes", n, len(bytes))
	}
}

func HTMLEscapeString(text string) string {
	text = strings.Replace(text, "&", "&amp;", -1)
	text = strings.Replace(text, "<", "&lt;", -1)
	return strings.Replace(text, ">", "&gt;", -1)
}
