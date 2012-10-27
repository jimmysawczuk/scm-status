# PHP SDK

The PHP SDK is pretty simple, with just two static methods on the `ScmStatus` class.

## `setFilepath($filepath)`

Pass this method a string (`$filepath`) at which the output of `scm-status` is stored.

## `format($formatting_string, array $options = array())`

Works a little like `printf` or `date`. Pass in a formatting string, get the string you want returned.

Here are the formatting options:

* `%T`: the repository type (e.g. `git`)
* `%n`: the decimal version of the revision (or short hex, if decimal isn't available)
* `%r`: the short version of the hex
* `%R`: the long version of the hex
* `%d`: the date of the commit in the format given by the SCM (e.g. `Sat Oct 27 13:41:06 EDT 2012`
* `%U`: the date of the commit in a UNIX timestamp format (e.g. `1351359666`)
* `%b`: the branch of the commit (e.g. `master`)
* `%t`: the tags of the commit separated by a delimiter
* `%a`: the committer's name
* `%e`: the committer's e-mail
* `%m`: the commit message
* `%s`: the commit subject

You can also pass in some options in the second parameter:

* `delimiter` (default: `,`): delimiter for any array of values
* `format_date` (default: none): a PHP `date`-formatted date string that you can use to create your own flag (`%F`)
* `return_if_fail` (default: nothing): what to return if this call fails for any reason