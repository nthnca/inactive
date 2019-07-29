# Command `inactive`

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/nthnca/inactive)
[![Go Report Card](https://goreportcard.com/badge/github.com/nthnca/inactive?style=flat-square)](https://goreportcard.com/report/github.com/nthnca/inactive)
[![Latest Version](https://img.shields.io/github/release/nthnca/inactive.svg?style=flat-square)](https://github.com/nthnca/inactive/releases)
[![Latest Version](https://img.shields.io/github/license/nthnca/inactive.svg?style=flat-square)](https://github.com/nthnca/inactive/blob/master/LICENSE)

Command line tool that can be used to determine if a computer has not been active
recently. I use it to make sure my machines don't get left on when I am not using them.

The basics of how it works:
- looks at all files named: `/tmp/stayawake.*`, if any of the files have been modified in
  the last 25 minutes the command exits with a status of 1.
- if the computer has an uptime of less than 5 minutes the command exits with a status of
  1.
- if neither of the above are true the command exits with a status of 0.

As a result in order to keep your computer awake you just need to touch a file that looks
like `/tmp/stayawake.\*`, in my case I use my [bash prompt to automatically touch a
file](https://github.com/nthnca/dotbash/blob/master/bash/stayawake.sh) of that sort.

Example usage in a crontab (I have it run every 3 minutes):
```
*/3 * * * * /root/inactive && \
    ( test -f /run/nologin || /sbin/shutdown -h +5 ) || /sbin/shutdown -c
```

To install you can simply run `go get github.com/nthnca/inactive`.

NOTES:
- This command doesn't actually shut down the computer because of possible permission
  issues and different OSes may have different ways to shut the computer down.
- The check for /run/nologin is to see if a shutdown has already been scheduled, if so
  we don't want to schedule another shutdown since it will cancel the first shutdown
  ... basically resulting in a denial of service for shutting down... oops.  :-)
- On my system to get the shutdown -h +5 command to work I had to run `sudo apt install
  --reinstall dbus`
