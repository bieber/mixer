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
	"github.com/bieber/logger"
	"github.com/bieber/mixer/mixerserver/context"
	"github.com/bieber/mixer/mixerserver/spotify"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

var loggerMutex = sync.Mutex{}

type submissionList struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
}

type submissionOptions struct {
	RoundRobin bool `json:"round_robin"`
	Shuffle    bool `json:"shuffle"`
	Dedup      bool `json:"dedup"`
	Pad        bool `json:"pad"`
}

type submissionData struct {
	SourceLists []submissionList  `json:"source_lists"`
	DestList    submissionList    `json:"dest_list"`
	Options     submissionOptions `json:"options"`
}

type trackLists [][]string

// Submit fires off a goroutine to actually mix the selected playlists
// into the destination list with the specified options.
func Submit(globalContext *context.GlobalContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		localContext := context.Get(r)

		userID, err := spotify.GetUserID(localContext.AuthTokens)
		if err != nil {
			panic(err)
		}

		data := submissionData{}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			panic(err)
		}

		go MixPlaylists(globalContext, localContext, userID, data)
	}
}

// MixPlaylists performs the mix operation triggered by the Submit
// handler.
func MixPlaylists(
	globalContext *context.GlobalContext,
	localContext *context.LocalContext,
	userID string,
	data submissionData,
) {
	t0 := time.Now()
	log := logger.New()

	defer func(log *logger.Logger) {
		if err := recover(); err != nil {
			log.Printf("====\nPANIC WHILE MIXING: %V", err)
			log.WriteTo(globalContext.LogOut)
		}
	}(log)

	sourceListIDs := []string{}
	for _, list := range data.SourceLists {
		sourceListIDs = append(sourceListIDs, list.ID)
	}

	log.WriteString("====\n")
	log.Printf(
		"MIXING [%s] INTO %s",
		strings.Join(sourceListIDs, ", "),
		data.DestList.ID,
	)
	log.Printf("ROUND ROBIN: %t", data.Options.RoundRobin)
	log.Printf("SHUFFLE:     %t", data.Options.Shuffle)
	log.Printf("DEDUP:       %t", data.Options.Dedup)
	log.Printf("PAD:         %t", data.Options.Pad)

	sourceTrackIDs := [][]string{}
	for _, list := range data.SourceLists {
		trackIDs, err := spotify.GetPlaylistTrackIDs(
			localContext.AuthTokens,
			list.OwnerID,
			list.ID,
		)
		if err != nil {
			panic(err)
		}

		sourceTrackIDs = append(sourceTrackIDs, trackIDs)
	}

	destTrackIDs, err := spotify.GetPlaylistTrackIDs(
		localContext.AuthTokens,
		data.DestList.OwnerID,
		data.DestList.ID,
	)
	if err != nil {
		panic(err)
	}

	combinedTrackIDs, err := CombineSourceTracks(sourceTrackIDs, data.Options)
	if err != nil {
		panic(err)
	}

	_, _ = destTrackIDs, combinedTrackIDs

	log.Printf("FINISHED IN %v", time.Now().Sub(t0))
	loggerMutex.Lock()
	_, err = log.WriteTo(globalContext.LogOut)
	loggerMutex.Unlock()

	if err != nil {
		panic(err)
	}
}

// CombineSourceTracks arranges the tracks from multiple source
// playlists into a single combined result playlist.
func CombineSourceTracks(
	sourceTrackIDs [][]string,
	options submissionOptions,
) (trackIDs []string, err error) {
	sort.Sort(trackLists(sourceTrackIDs))

	return
}

func (ls trackLists) Len() int {
	return len(ls)
}

func (ls trackLists) Less(i, j int) bool {
	return len(ls[i]) < len(ls[j])
}

func (ls trackLists) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
