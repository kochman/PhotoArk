# PhotoArk [![Build Status](https://travis-ci.org/kochman/PhotoArk.svg?branch=master)](https://travis-ci.org/kochman/PhotoArk)

PhotoArk is a web interface for viewing photos.

![PhotoArk screenshot](https://cloud.githubusercontent.com/assets/335234/13827464/e0c4531a-eb79-11e5-8478-055fbe855290.png)

## Running

`go get github.com/kochman/PhotoArk`

Inside the PhotoArk directory:

- `bower install` - download web components.
- `go get github.com/mjibson/esc` - install [esc](https://github.com/mjibson/esc) for embedding static assets.
- `go generate` - bundle the static folder for embedding.
- `go run main.go syncmap.go static.go -photoDir photos -cacheDir cache`

To build a binary, do the above, and then:

`go build`

The `-devel` flag will force PhotoArk to serve static files off of the disk instead of from the embedded static folder.

## Metadata

PhotoArk looks for `metadata.yaml` files inside folders in `photoDir`. These metadata files allow users to filter photos. An example metadata file:

```
event: Fashion Show
photographer: Sidney Kochman
date: 2016-2-4
location: Auditorium
```

Currently, only those four fields are exposed in the web interface. In the future, arbitrary keys will be visible.
