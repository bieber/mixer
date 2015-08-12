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
	"net/http"
	"net/url"
)

// GetUserID fetches the Spotify user ID of the logged-in user.
func GetUserID(authTokens AuthTokens) (userID string, err error) {
	client := &http.Client{}
	uri, err := url.Parse("https://api.spotify.com/v1/me")
	if err != nil {
		return
	}

	request, err := NewAuthenticatedRequest(authTokens, "GET", uri, nil)
	if err != nil {
		return
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	output := struct {
		UserID string `json:"id"`
	}{}
	err = json.NewDecoder(response.Body).Decode(&output)
	if err != nil {
		return
	}

	userID = output.UserID
	return
}
