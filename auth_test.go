package main

import (
	`testing`
	`github.com/jscherff/gotest`
)

func TestAuth(t *testing.T) {

	var err error

	goodUsername := conf.Server.Auth.Username
	goodPassword := conf.Server.Auth.Password

	t.Run(`Success with Good Credentials`, func(t *testing.T) {

		err = auth()
		gotest.Ok(t, err)
	})

	t.Run(`Failure with Bad Username`, func(t *testing.T) {

		mux.Lock()
		defer mux.Unlock()

		authenticated = false
		conf.Server.Auth.Username = `baduser`
		conf.Server.Auth.Password = goodPassword

		err = auth()
		conf.Server.Auth.Username = goodUsername
		gotest.Assert(t, err != nil, `authentication with bad username should fail`)
	})

	t.Run(`Failure with Bad Password`, func(t *testing.T) {

		mux.Lock()
		defer mux.Unlock()

		authenticated = false
		conf.Server.Auth.Username = goodUsername
		conf.Server.Auth.Password = `badpass`

		err = auth()
		conf.Server.Auth.Password = goodPassword
		gotest.Assert(t, err != nil, `authentication with bad password should fail`)
	})
}
