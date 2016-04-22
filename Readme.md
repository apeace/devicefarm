# DeviceFarm

Work in progress.

## Dev setup

 * [Install Go](https://golang.org/doc/install)
 * [Setup a Go workspace](https://golang.org/doc/code.html)

Make sure you setup `$GOPATH` and add `$GOPATH/bin` to your `$PATH`.

```
go get github.com/tools/godep
mkdir -p $GOPATH/src/github.com/ride/
cd $GOPATH/src/github.com/ride/
git clone git@github.com:ride/devicefarm.git
cd devicefarm
```

## Running

Right now the main binary just validates a config file:

```
go run main.go config/testdata/config.yml
```

## Testing

To run tests:

```
godep go test ./...
```

To get a more detailed report, including coverage info, run this script:

```
# run tests with coverage
./test.sh
# open coverage report
open coverage.html
```

## Documentation

To view docs:

```
go build ./...
godoc -http=:8080
open "http://localhost:8080/pkg/github.com/ride/devicefarm/"
```

## Adding Dependencies

Use `godep` from the project root to add new dependencies. For example,
here's how you would add package `foo/bar`:

```
cd $GOPATH/src/github.com/ride/devicefarm
go get foo/bar
# edit code to import foo/bar
godep save ./...
```

The package should appear in the `vendor/` directory, which you should
commit.
