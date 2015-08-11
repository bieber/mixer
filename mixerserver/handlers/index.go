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
	"github.com/bieber/mixer/mixerserver/context"
	"github.com/bieber/mixer/mixerserver/crypto"
	"github.com/bieber/mixer/mixerserver/util"
	"net/http"
	"net/url"
	"strings"
)

// Index renders the homepage.
func Index(globalContext *context.GlobalContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfData := map[string]string{
			"User-Agent": r.Header.Get("User-Agent"),
			"IP":         util.StripPort(r.RemoteAddr),
		}
		csrfPlaintext, err := json.Marshal(csrfData)
		if err != nil {
			panic(err)
		}
		csrfToken, err := crypto.Encrypt(string(csrfPlaintext))

		scopes := []string{
			"playlist-read-private",
			"playlist-read-collaborative",
			"playlist-modify-public",
			"playlist-modify-private",
		}

		loginCompletionURI, err := loginURI(globalContext, r.Host)
		if err != nil {
			panic(err)
		}

		loginURI, err := url.Parse("https://accounts.spotify.com/authorize/")
		if err != nil {
			panic(err)
		}
		loginURI.RawQuery = (url.Values{
			"client_id":     []string{globalContext.Spotify.ClientID},
			"response_type": []string{"code"},
			"state":         []string{csrfToken},
			"scope":         []string{strings.Join(scopes, " ")},
			"redirect_uri":  []string{loginCompletionURI.String()},
		}).Encode()

		playlistURI, err := globalContext.Router.Get("playlists").URL()
		if err != nil {
			panic(err)
		}

		err = globalContext.Templates.Index.Execute(
			w,
			map[string]interface{}{
				"loginURI":     loginURI.String(),
				"playlistsURI": playlistURI.String(),
			},
		)
		if err != nil {
			panic(err)
		}
	}
}

// loginURI assembles the login URI to redirect to from the Spotify
// login API.
func loginURI(
	globalContext *context.GlobalContext,
	host string,
) (*url.URL, error) {
	loginCompletionURI, err := globalContext.Router.Get("login").URL()
	if err != nil {
		return nil, err
	}
	loginCompletionURI.Scheme = "http"
	loginCompletionURI.Host = host
	return loginCompletionURI, nil
}
