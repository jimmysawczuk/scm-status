# scm-status

**scm-status** is a tool to quickly generate a file that snapshots where your current working copy is in development. It's useful for knowing what version of your code your production site or app is running.

## Support

Right now, scm-status supports Git and Mercurial and has been tested on Linux and Mac OS X.

## Installing on your system

* [Install go](http://golang.org/doc/install) (any version >= 1.0 should do)
* `go get github.com/jimmysawczuk/scm-status/cmd/scm-status`

## Installing on your repository

* From your repository's working path directory, `scm-status setup`. This will install all the hooks you need to keep the snapshot file updated.
  * You can change the path of the executable used by the hook using the `-executable` flag, whether or not you want compressed output via `-pretty`, and the output file via the `-out` flag.

## Using

* Run `scm-status` to generate your snapshot. The output defaults to STDOUT, but can be redirected to a file using the `-out` flag.
* Then, parse the output or file using whatever programming language you wish as JSON, and use whatever you need!
  * If you're using PHP, feel free to use the included SDK.
  * If you're using another language and wish to contribute an SDK I'd love to see it.
* Use the `-pretty` flag to control the format of the output.
