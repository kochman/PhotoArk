# PhotoArk

PhotoArk is a web interface for viewing photos.

## Running

`go get github.com/kochman/PhotoArk`

Inside the PhotoArk directory:

- `bower install` - download web components.
- `go get` - download Go dependencies.
- `go generate` - bundle the static folder for embedding.
- `go run main.go static.go -photoDir photos -cacheDir cache`

To build a binary, do the above, and then:

`go build`

## Metadata

PhotoArk looks for `metadata.yaml` files inside folders in `photoDir`. These metadata files allow users to filter photos. An example metadata file:

```
event: Fashion Show
photographer: Sidney Kochman
date: 2016-2-4
location: Auditorium
```

Currently, only those four fields are exposed in the web interface. In the future, arbitrary keys will be visible.
