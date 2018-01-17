# fileserve

fileserve acts as a HTTP file server supporting GET, HEAD, PUT, DELETE in a specified directory.
Useful for CockroachDB IMPORT/BACKUP/RESTORE testing with local `http` storage.
No security measures are taken whatsoever.
Only for use in local testing.

# options

- `-addr` serve address, defaults to `:1123`
- `-base` file directory, defaults to `.`
