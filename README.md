# scm-status

**scm-status** is a tool to quickly generate a file that snapshots where your current working copy is in development. It's useful for knowing what version of your code your production site or app is running.

## Support

Right now, scm-status supports **Git** and **Mercurial (hg)** on **POSIX** systems.

## Building

* [Install go](http://golang.org/doc/install) (any version >= 1.0 should do)
* `go get github.com/jimmysawczuk/go-binary` (to embed version information)
* `git clone http://github.com/jimmysawczuk/scm-status.git && cd scm-status`
* `make && make install`

## Installing on your repository

* From your repository's working path directory, `scm-status setup`. This will install all the hooks you need to keep the snapshot file updated.
  * You can change the path of the executable used by the hook using the `-executable` flag, whether or not you want compressed output via `-pretty`, and the output file via the `-out` flag.

## Using

* Run `scm-status` to generate your snapshot. The output defaults to STDOUT, but can be redirected to a file using the `-out` flag.
* Then, parse the output or file using whatever programming language you wish as JSON, and use whatever you need!
  * If you're using PHP, feel free to use the included SDK.
  * If you're using another language and wish to contribute an SDK I'd love to see it.
* Use the `-pretty` flag to control the format of the output.

## License

    The MIT License (MIT)
    Copyright (C) 2012-2014 by Jimmy Sawczuk

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
    THE SOFTWARE.