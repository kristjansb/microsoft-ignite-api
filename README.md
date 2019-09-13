# microsoft-ignite-api

This Go tool will pull the latest MS Ignite search results from the API then store the results in a JSON file and a number of CSV files based on session type.

## Building

This Golang command line tool requires a Go development environment to build. If necessary, install the latest Go release for your system first https://golang.org/dl/. Once installed, pull and build/install with the following command:

`go get github.com/kristjansb/microsoft-ignite-api`

This will clone the tool into `$GOPATH/src/github.com/kristjansb/microsoft-ignite-api` and install the executable in `$GOPATH/bin/microsoft-ignite-api.exe`.

To run the tool without installing it, use `go run main.go` in the src path above.

## Running the executable

Place `microsoft-ignite-api.exe` somewhere on your path or in the current directory. Running the command will query the Ignite search API, create or overwrite CSV and JSON files in the current directory, and print the facets summary to screen.
