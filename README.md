# PhotoArk

PhotoArk is a web interface for viewing photos.

## Running

`go get github.com/kochman/PhotoArk`

Inside the PhotoArk directory, run `go generate` to bundle the static folder for embedding. Then, to run PhotoArk: `go run main.go static.go -photoDir photos -cacheDir cache`. Replace `photos` and `cache` with paths to photo and cache directories, respectively.

## Metadata

PhotoArk looks for `metadata.yaml` files inside folders in `photoDir`. These metadata files allow users to filter photos. An example metadata file:

```
event: Fashion Show
photographer: Sidney Kochman
date: 2016-2-4
location: Auditorium
```

Currently, only those four fields are exposed in the web interface. In the future, arbitrary keys will be visible.
