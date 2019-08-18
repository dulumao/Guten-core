package session

import (
	"github.com/gorilla/sessions"
	"github.com/dulumao/Guten-core/env"
)

var Value sessions.Store

func New() {
	if env.Value.Session.Driver == "cookie" {
		Value = NewCookieStore([]byte(env.Value.Server.HashKey))
	}

	if env.Value.Session.Driver == "file" {
		Value = NewFilesystemStore(env.Value.Session.File.Path, []byte(env.Value.Server.HashKey))
	}

	if env.Value.Session.Driver == "redis" {
		var err error

		Value, err = NewRedisStore(32, "tcp", env.Value.Session.Redis.Addr, env.Value.Session.Redis.Password, []byte(env.Value.Server.HashKey))

		if err != nil {
			panic(err)
		}
	}
}
