# Go Serverless

## How this project was initialized
```
$ mkdir golang-serverless
$ cd golang-serverless
$ go mod init github.com/danielmonteiro/go-serverless
```

## Generating the deployable package
- `go mod tidy` ensures that the go.mod file matches the source code in the module. It adds any missing module requirements necessary to build the current module’s packages and dependencies, and it removes requirements on modules that don’t provide any relevant packages. It also adds any missing entries to go.sum and removes unnecessary entries.
- From the `cmd` folder run `go build main.go` and then move the builded file inside the `build` folder
- From the root directory of the project run `zip -jrm build/main.zip build/main`
- The builded `build/main.zip` file will be uploaded to AWS Lambda

### Resources
- https://www.youtube.com/watch?v=zHcef4eHOc8