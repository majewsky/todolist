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

## User management

Credentials are stored in the `todolist_passwd` file (in the working directory
of the application). To create a new user account, log on with the new
username/password and copy the displayed line into `todolist_passwd`.

To reset a user's password, log on with the username and new password, then
copy the displayed line into `todolist_passwd` while replacing the existing
line with that user's name.

The todolist data for a user is stored in the application's working directory
as `todolist-$username.txt`.
