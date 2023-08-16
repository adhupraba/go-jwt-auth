# CompileDaemon

we use `compiledaemon` to automatically restart the server whenever any file changes (like `nodemon`)

```bash
compiledaemon --command="./<module_name>"
```

example:

```bash
compiledaemon --command="./go-jwt-auth"
```

> module_name is the name given in the `go mod init` command. it can be found in the `go.mod` file also at the end of the github url
