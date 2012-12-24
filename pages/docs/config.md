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
* If `lockview` is "true", then the user must be authenticated to view the
wiki.
* `keyfile` and `certfile`, if set, will cause the wiki to serve over SSL (TLS).

### server

* `address` specifies an address other than 127.0.0.1 to listen on
* `port` specifies a port to listen on

## Some example config files

### Basic authentication
This config file implements basic authentication with a username and password; anyone can 
view the wiki, but only authenticated users can edit it or delete pages.

```
[security]
username = user
password = secret
```

### Listen on all addresses
By default, `gowik` will only listen on the localhost. We can tell it to listen on all available addresses by using an empty address:

```
[server]
address = 
```