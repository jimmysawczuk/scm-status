scm-status [![Go Report Card](https://goreportcard.com/badge/github.com/jimmysawczuk/scm-status)](https://goreportcard.com/report/github.com/jimmysawczuk/scm-status)
==============

**scm-status** is a tool to quickly generate a file that snapshots where your current working copy is in development. It's useful for knowing what version of your code your production site or app is running.

## Support

Right now, scm-status supports Git and Mercurial and has been tested on Linux and Mac OS X.

## Installing on your system

* [Install go](http://golang.org/doc/install) (any version >= 1.6 should do)
* `go get github.com/jimmysawczuk/scm-status`

## Using

- Run `scm-status` to generate your snapshot. The output defaults to STDOUT, but can be redirected to a file using the `-out` flag. Formatted output is on by default, but you can turn it off via `-pretty=false`.
- Then, parse the output or file using whatever programming language you wish as JSON, and use whatever you need!
- See `scm-status -help` for a complete command reference.

## Installing on your repository

scm-status can hook itself into your repository automatically to keep itself up to date.

- From your repository's working path directory, `scm-status -install-hooks -out=REVISION.json`. This will install all the hooks you need to keep the snapshot file updated.
  - You can optionally turn off formatted output using `-pretty=false`.
  - If you install the executable somewhere other than `$GOPATH/bin/scm-status` you'll need to pass that in as a flag as well.

## License

[MIT](https://github.com/jimmysawczuk/scm-status/blob/master/LICENSE)