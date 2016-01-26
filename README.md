# todolist
Tiny single-user todolist app for my personal use

## Building

```bash
$ make
$ ./todolist
starting todolist server on port 8080
```

In a development setup, the application will serve its static files by itself
from the `./static` directory. Therefore, you have to start the application in
the root directory (or in a directory that contains a symlink to the `static`
directory).

## Deploying

In productive deployments, setup an HTTP server to terminate TLS, serve the
static files below `/static`, and proxy all remaining requests to the
application's internal HTTP port.
