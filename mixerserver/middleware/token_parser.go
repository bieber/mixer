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

package middleware

import (
	"encoding/json"
	"errors"
	"github.com/bieber/mixer/mixerserver/context"
	"github.com/bieber/mixer/mixerserver/crypto"
	"net/http"
)

// TokenParser looks for a "token" GET parameter, decrypts and parses
// it, and kills the request if anything fails along the way.
func TokenParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localContext := context.Get(r)

		token := r.URL.Query().Get("token")
		decryptedToken, err := crypto.Decrypt(token)
		if err != nil {
			panic(err)
		}

		tokenData := struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{}
		err = json.Unmarshal([]byte(decryptedToken), &tokenData)
		if err != nil {
			panic(err)
		}

		if tokenData.AccessToken == "" || tokenData.RefreshToken == "" {
			panic(errors.New("Missing access or refresh token"))
		}

		localContext.Spotify.AccessToken = tokenData.AccessToken
		localContext.Spotify.RefreshToken = tokenData.RefreshToken

		next.ServeHTTP(w, r)
	})
}
