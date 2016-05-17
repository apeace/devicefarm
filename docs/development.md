# Development

This doc is for **developers making changes to this repo**. If you are an end-user (e.g. a mobile developer) you want the [Readme](../Readme.md).

## Dev setup

 * [Install Go](https://golang.org/doc/install)
 * [Setup a Go workspace](https://golang.org/doc/code.html)

Make sure you setup `$GOPATH` and add `$GOPATH/bin` to your `$PATH`.

```
go get -u github.com/kardianos/govendor
mkdir -p $GOPATH/src/github.com/ride/
cd $GOPATH/src/github.com/ride/
git clone git@github.com:ride/devicefarm.git
cd devicefarm
```

## Testing

To run tests:

```
govendor test +local
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

```
# go to the project root
cd $GOPATH/src/github.com/ride/devicefarm

# add dependency foo/bar (usually will be something like github.com/foo/bar)
go get foo/bar

# now edit code to import foo/bar and use it...

# now save the dependency
govendor add +external

# and commit it
git add --all vendor/
```

## Releasing

Just push a tag. Circle will create the release on Github automatically.