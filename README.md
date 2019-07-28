# inactive


[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/nthnca/inactive)
[![Go Report Card](https://goreportcard.com/badge/github.com/nthnca/inactive?style=flat-square)](https://goreportcard.com/report/github.com/nthnca/inactive)
[![Latest Version](https://img.shields.io/github/release/nthnca/inactive.svg?style=flat-square)](https://github.com/nthnca/inactive/releases)
[![Latest Version](https://img.shields.io/github/license/nthnca/inactive.svg?style=flat-square)](https://github.com/nthnca/inactive/blob/master/LICENSE)

Command line tool that can be used to determine if a computer has not been active recently. I use it to make sure my machines don't get left on when I am not using them.

The basics of how it works:
- looks at all files named: `/tmp/stayawake.*`, if any of the files have been modified in the last
  25 minutes the command exits with a status of 1.
- if the computer has an uptime of less than 5 minutes the command exits with a status of 1.
- if neither of the above are true the command exits with a status of 0.

As a result in order to keep your computer awake you just need to touch a file that looks like /tmp/stayawake.\*, in
my case I use my bash prompt to automatically touch a file of that sort.

Example usage in a crontab (I have it run every 2 minutes):
```
inactive | logger; test ${PIPESTATUS[0]} -eq 0 && /sbin/shutdown -h +5 || /sbin/shutdown -c "Shutdown cancelled"
```

NOTES:
- We don't use this command to actually shutdown the computer because of possible permission issues and different
  OSes may have different ways to shut the computer down.
