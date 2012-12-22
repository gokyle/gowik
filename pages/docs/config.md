## configuration options

### wiki

* `pages` tells `gowik` where the wiki tree is; this should be where all the
wiki pages are stored.
* `extension` tells the wiki which extension should be used to indicate a
page; the recommended (and default) value is ".md".

### security

* `username` should hold the user's name; can be blank
* `password` is the user's password; if blank no authentication will be
used.
* If 'LockView' is "true", then the user must be authenticated to view the
wiki.
* `keyfile` and `certfile`, if set, will cause the wiki to serve over SSL (TLS).

### server

* `address` specifies an address other than 127.0.0.1 to listen on
* `port` specifies a port to listen on

