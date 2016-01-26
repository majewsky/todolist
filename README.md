# todolist
Tiny single-user todolist app for my personal use

## Building

```bash
$ make
$ ./todolist
starting todolist server on port 8080
```

## Deploying

In productive deployments, setup an HTTP server to terminate TLS and
`proxy_pass` the requests to the application.
