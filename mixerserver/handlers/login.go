/*
 * Copyright 2015, Robert Bieber
 *
 * This file is part of mixer.
 *
 * mixer is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mixer is distributed in the hope that it will be useful,
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with mixer.  If not, see <http://www.gnu.org/licenses/>.
 */

package handlers

import (
	"encoding/json"
	"errors"
	"github.com/bieber/mixer/mixerserver/context"
	"github.com/bieber/mixer/mixerserver/crypto"
	"github.com/bieber/mixer/mixerserver/spotify"
	"github.com/bieber/mixer/mixerserver/util"
	"net/http"
)

// Login handles login responses from the Spotify API
func Login(globalContext *context.GlobalContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfJSON, err := crypto.Decrypt(r.URL.Query().Get("state"))
		if err != nil {
			panic(err)
		}

		csrf := make(map[string]string)
		err = json.Unmarshal([]byte(csrfJSON), &csrf)
		if err != nil {
			panic(err)
		}

		if ip, ok := csrf["IP"]; !ok || ip != util.StripPort(r.RemoteAddr) {
			panic(errors.New("CSRF mismatch"))
		}
		if a, ok := csrf["User-Agent"]; !ok || a != r.Header.Get("User-Agent") {
			panic(errors.New("CSRF mismatch"))
		}

		data := map[string]interface{}{
			"error": r.URL.Query().Get("error"),
		}

		if r.URL.Query().Get("error") == "" {
			redirectURI, err := loginURI(globalContext, r.Host)
			if err != nil {
				panic(err)
			}

			tokens, err := spotify.GetAuthTokens(
				globalContext.Spotify.ClientID,
				globalContext.Spotify.ClientSecret,
				r.URL.Query().Get("code"),
				redirectURI,
			)

			if err != nil {
				panic(err)
			}

			data["expires_in"] = tokens.ExpiresIn

			jsonToken, err := json.Marshal(tokens)
			if err != nil {
				panic(err)
			}

			token, err := crypto.Encrypt(string(jsonToken))
			if err != nil {
				panic(err)
			}
			data["token"] = token
		}

		err = globalContext.Templates.Login.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

// Refresh fetches new auth tokens for an existing session that has
// expired.
func Refresh(globalContext *context.GlobalContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		localContext := context.Get(r)

		tokens, err := spotify.RefreshAuthTokens(
			localContext.AuthTokens,
			globalContext.Spotify.ClientID,
			globalContext.Spotify.ClientSecret,
		)
		if err != nil {
			panic(err)
		}

		data := map[string]interface{}{
			"expires_in": tokens.ExpiresIn,
		}

		jsonToken, err := json.Marshal(tokens)
		if err != nil {
			panic(err)
		}

		token, err := crypto.Encrypt(string(jsonToken))
		if err != nil {
			panic(err)
		}
		data["token"] = token

		w.Header().Set("Content-type", "application/json")
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			panic(err)
		}
	}
}
