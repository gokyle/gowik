package main

import (
	"github.com/gokyle/webshell/auth"
	"net/http"
	"strings"
)

var Security struct {
	Enabled  bool // Is authentication enabled?
	AuthView bool // Is authentication required to view pages?
	TLS      struct {
		Enabled bool // Create TLS web app?
		Key     string
		Cert    string
	}
	User struct {
		Name string
		Salt []byte
		Hash []byte
	}
	SessionStore *auth.SessionStore
}

// Initialise security options
func initSecurity(cfgSec map[string]string) {
	if cfgSec == nil {
		return
	}
	for key, val := range cfgSec {
		key = strings.ToLower(key)
		switch key {
		case "username":
			Security.User.Name = val
		case "password":
			salt, hash := auth.HashPass(val)
			Security.User.Salt = salt
			Security.User.Hash = hash
		case "lockview":
			Security.AuthView = true
		case "certfile":
			Security.TLS.Cert = val
		case "keyfile":
			Security.TLS.Key = val
		}
	}

	if Security.TLS.Key != "" && Security.TLS.Cert != "" {
		Security.TLS.Enabled = true
	}

	if len(Security.User.Hash) != 0 && len(Security.User.Salt) != 0 {
		Security.Enabled = true
		Security.SessionStore = auth.CreateSessionStore(
			"gowik_as",
			Security.TLS.Enabled,
			nil,
		)
		auth.LookupCredentials = authenticate
	}
}

func authorised(update bool, r *http.Request) bool {
	// if security isn't enabled, all users are authorised to do anything
	if !Security.Enabled {
		return true
	}

	// if authentication isn't required to view pages and the user is
	// viewing a page, let them.
	if !Security.AuthView && !update {
		return true
	}

	if !Security.SessionStore.CheckSession(r) {
		return false
	}
	return true
}

func authenticate(user interface{}) (salt, hash []byte) {
	if !Security.Enabled {
		return
	}
	if user.(string) == Security.User.Name {
		return Security.User.Salt, Security.User.Hash
	}
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// show login form
	}
	r.ParseForm()
	user := r.FormValue("user")
	pass := r.FormValue("pass")

	if !authorised(true, r) {
		cookie, err := Security.SessionStore.AuthSession(user, pass, false, "")
		if err != nil || cookie == nil {
			LoginFailed(w, r)
			return
		}

	}
}

func LoginFailed(w http.ResponseWriter, r *http.Request) {

}
