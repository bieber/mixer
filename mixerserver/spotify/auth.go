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

package spotify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// AuthTokens stores information about (and accepts JSON requests for)
// Spotify API access tokens.  The AccessToken is used to access the
// API endpoints, the RefreshToken to get a new token after ExpiresIn
// seconds have elapsed.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// GetLoginURI generates a login URI you can direct a user to to
// authenticate the Spotify API.
func GetLoginURI(
	clientID string,
	csrfToken string,
	redirectURI *url.URL,
) (*url.URL, error) {
	scopes := []string{
		"playlist-read-private",
		"playlist-read-collaborative",
		"playlist-modify-public",
		"playlist-modify-private",
	}

	loginURI, err := url.Parse("https://accounts.spotify.com/authorize/")
	if err != nil {
		return nil, err
	}
	loginURI.RawQuery = url.Values{
		"client_id":     []string{clientID},
		"response_type": []string{"code"},
		"state":         []string{csrfToken},
		"scope":         []string{strings.Join(scopes, " ")},
		"redirect_uri":  []string{redirectURI.String()},
	}.Encode()

	return loginURI, nil
}

// GetAuthTokens fetches authentication tokens from the Spotify server
// given an access code returned in the redirect from the login page
// (or a refresh token).  redirectURI is required to authenticate the
// request, and should be exactly the same as the redirect_uri that
// was initially sent to Spotify.
func GetAuthTokens(
	clientID string,
	clientSecret string,
	code string,
	redirectURI *url.URL,
) (out AuthTokens, err error) {
	response, err := http.PostForm(
		"https://accounts.spotify.com/api/token",
		url.Values{
			"grant_type":    []string{"authorization_code"},
			"code":          []string{code},
			"redirect_uri":  []string{redirectURI.String()},
			"client_id":     []string{clientID},
			"client_secret": []string{clientSecret},
		},
	)
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&out)
	response.Body.Close()
	return
}

// NewAuthenticatedRequest returns a new *http.Request with the
// authentication headers set for the Spotify API.  Aside from the
// authTokens argument, it is equivalent to http.NewRequest.
func NewAuthenticatedRequest(
	authTokens AuthTokens,
	method string,
	uri *url.URL,
	body io.Reader,
) (request *http.Request, err error) {
	request, err = http.NewRequest(method, uri.String(), body)
	if err != nil {
		return
	}

	request.Header.Set("Authorization", "Bearer "+authTokens.AccessToken)
	return
}
