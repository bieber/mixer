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
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

const playlistBatchSize = 30
const trackFetchBatchSize = 100
const trackWriteBatchSize = 100

// Playlist lists all the vital information for a Spotify playlist.
type Playlist struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Collaborative bool   `json:"collaborative"`
	Public        bool   `json:"public"`
	Owner         struct {
		ID string `json:"id"`
	} `json:"owner"`
}

// GetPlaylists fetches all the playlists of the given user.
func GetPlaylists(
	authTokens AuthTokens,
	userID string,
) (playlists []Playlist, err error) {
	playlists = []Playlist{}

	fetchURI, err := url.Parse(
		"https://api.spotify.com/v1/users/" + userID + "/playlists",
	)
	if err != nil {
		return
	}

	client := &http.Client{}
	var request *http.Request
	var response *http.Response
	for batch := 0; true; batch++ {
		fetchURI.RawQuery = url.Values{
			"offset": []string{strconv.Itoa(batch * playlistBatchSize)},
			"limit":  []string{strconv.Itoa(playlistBatchSize)},
		}.Encode()

		request, err = NewAuthenticatedRequest(authTokens, "GET", fetchURI, nil)
		if err != nil {
			return
		}

		response, err = client.Do(request)
		if err != nil {
			return
		}
		defer response.Body.Close()

		result := struct {
			Playlists []Playlist `json:"items"`
			Next      string     `json:"next"`
		}{}
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return
		}

		playlists = append(playlists, result.Playlists...)

		if result.Next == "" {
			break
		}
	}
	return
}

// GetPlaylistTrackIDs returns the IDs of all the tracks in the given
// playlist.  Note that some inconsistency could result here if
// someone adds or removes tracks in between batches, but that's not a
// serious enough issue to bother with for now.
func GetPlaylistTrackIDs(
	authTokens AuthTokens,
	userID string,
	playlistID string,
) (trackIDs []string, err error) {
	trackIDs = []string{}

	fetchURI, err := url.Parse("" +
		"https://api.spotify.com/v1/users/" +
		userID +
		"/playlists/" +
		playlistID +
		"/tracks",
	)
	if err != nil {
		return
	}

	client := &http.Client{}
	var request *http.Request
	var response *http.Response
	for batch := 0; true; batch++ {
		fetchURI.RawQuery = url.Values{
			"offset": []string{strconv.Itoa(batch * trackFetchBatchSize)},
			"limit":  []string{strconv.Itoa(trackFetchBatchSize)},
			"fields": []string{"items(track(id)),next"},
		}.Encode()

		request, err = NewAuthenticatedRequest(authTokens, "GET", fetchURI, nil)
		if err != nil {
			return
		}

		response, err = client.Do(request)
		if err != nil {
			return
		}
		defer response.Body.Close()

		result := struct {
			Tracks []struct {
				Track struct {
					ID string `json:"id"`
				} `json:"track"`
			} `json:"items"`
			Next string `json:"next"`
		}{}
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return
		}

		for _, track := range result.Tracks {
			trackIDs = append(trackIDs, track.Track.ID)
		}

		if result.Next == "" {
			break
		}
	}

	return
}

// WritePlaylist deletes all the existing tracks in the given playlist
// and replaces them with the specified contents.
func WritePlaylist(
	authTokens AuthTokens,
	destListOwnerID string,
	destListID string,
	trackIDs []string,
) error {
	destListTrackIDs, err := GetPlaylistTrackIDs(
		authTokens,
		destListOwnerID,
		destListID,
	)
	if err != nil {
		return err
	}

	seenIDs := map[string]bool{}
	toDelete := []string{}
	for _, track := range destListTrackIDs {
		if _, ok := seenIDs[track]; ok {
			continue
		}

		seenIDs[track] = true
		toDelete = append(toDelete, "spotify:track:"+track)
	}

	deleteBatches := len(toDelete) / trackWriteBatchSize
	if len(toDelete)%trackWriteBatchSize != 0 {
		deleteBatches++
	}

	deleteURI, err := url.Parse("" +
		"https://api.spotify.com/v1/users/" +
		destListOwnerID +
		"/playlists/" +
		destListID +
		"/tracks",
	)
	if err != nil {
		return err
	}

	client := &http.Client{}
	for batch := 0; batch < deleteBatches; batch++ {
		data := map[string][]map[string]string{
			"tracks": []map[string]string{},
		}

		batchStart := batch * trackWriteBatchSize
		batchEnd := batchStart + trackWriteBatchSize
		for i := batchStart; i < batchEnd; i++ {
			if i >= len(toDelete) {
				break
			}
			data["tracks"] = append(
				data["tracks"],
				map[string]string{"uri": toDelete[i]},
			)
		}

		body := bytes.NewBuffer([]byte{})
		err := json.NewEncoder(body).Encode(data)
		if err != nil {
			return err
		}

		request, err := NewAuthenticatedRequest(
			authTokens,
			"DELETE",
			deleteURI,
			body,
		)
		if err != nil {
			return err
		}

		response, err := client.Do(request)
		if err != nil {
			return err
		}
		if response.StatusCode != http.StatusOK {
			return errors.New(response.Status)
		}
		response.Body.Close()
	}

	writeBatches := len(trackIDs) / trackWriteBatchSize
	if len(trackIDs)%trackWriteBatchSize != 0 {
		writeBatches++
	}

	writeURI, err := url.Parse("" +
		"https://api.spotify.com/v1/users/" +
		destListOwnerID +
		"/playlists/" +
		destListID +
		"/tracks",
	)
	if err != nil {
		return err
	}

	for batch := 0; batch < writeBatches; batch++ {
		data := map[string][]string{"uris": []string{}}

		batchStart := batch * trackWriteBatchSize
		batchEnd := batchStart + trackWriteBatchSize
		for i := batchStart; i < batchEnd; i++ {
			if i >= len(trackIDs) {
				break
			}
			data["uris"] = append(data["uris"], "spotify:track:"+trackIDs[i])
		}

		body := bytes.NewBuffer([]byte{})
		err := json.NewEncoder(body).Encode(data)
		if err != nil {
			return err
		}

		request, err := NewAuthenticatedRequest(
			authTokens,
			"POST",
			writeURI,
			body,
		)
		if err != nil {
			return err
		}

		response, err := client.Do(request)
		if err != nil {
			return err
		}
		if response.StatusCode != http.StatusCreated {
			return errors.New(response.Status)
		}
		response.Body.Close()
	}

	return nil
}
