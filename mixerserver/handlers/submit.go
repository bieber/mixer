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
	"math/rand"
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

		go mixPlaylists(globalContext, localContext, userID, data)
	}
}

// mixPlaylists performs the mix operation triggered by the Submit
// handler.
func mixPlaylists(
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

	combinedTrackIDs := combineSourceTracks(sourceTrackIDs, data.Options)

	err := spotify.WritePlaylist(
		localContext.AuthTokens,
		data.DestList.OwnerID,
		data.DestList.ID,
		combinedTrackIDs,
	)
	if err != nil {
		panic(err)
	}

	log.Printf("FINISHED IN %v", time.Now().Sub(t0))
	loggerMutex.Lock()
	_, err = log.WriteTo(globalContext.LogOut)
	loggerMutex.Unlock()

	if err != nil {
		panic(err)
	}
}

func combineSourceTracks(
	sourceTrackIDs [][]string,
	options submissionOptions,
) []string {
	if options.Dedup {
		sourceTrackIDs = dedupSourceTracks(sourceTrackIDs)
	}
	if options.Shuffle {
		sourceTrackIDs = shuffleSourceTracks(sourceTrackIDs, options.Pad)
	}
	// If both Pad and Shuffle were set, the tracks have already been
	// shuffled and padded
	if options.Pad && !options.Shuffle {
		sourceTrackIDs = padSourceTracks(sourceTrackIDs)
	}

	totalLength := 0
	for _, list := range sourceTrackIDs {
		totalLength += len(list)
	}
	destList := make([]string, totalLength)

	srcList := 0
	srcPositions := make([]int, len(sourceTrackIDs))

	for i := range destList {
		destList[i] = sourceTrackIDs[srcList][srcPositions[srcList]]
		if i == len(destList)-1 {
			break
		}

		srcPositions[srcList]++
		if options.RoundRobin {
			srcList = (srcList + 1) % len(sourceTrackIDs)
		}
		for srcPositions[srcList] >= len(sourceTrackIDs[srcList]) {
			srcList = (srcList + 1) % len(sourceTrackIDs)
		}
	}

	return destList
}

func dedupSourceTracks(sourceTrackIDs [][]string) [][]string {
	sort.Sort(trackLists(sourceTrackIDs))

	seenIDs := map[string]bool{}
	deduped := [][]string{}

	for _, list := range sourceTrackIDs {
		newList := []string{}

		for _, track := range list {
			if _, ok := seenIDs[track]; ok {
				continue
			}

			seenIDs[track] = true
			newList = append(newList, track)
		}

		deduped = append(deduped, newList)
	}

	return deduped
}

// If a list is being both padded and shuffled, the padding needs to
// happen at the same time as the shuffling so we can make sure not to
// include duplicates before the entire list has been exhausted.
func shuffleSourceTracks(sourceTrackIDs [][]string, pad bool) [][]string {
	maxLength := 0
	for _, list := range sourceTrackIDs {
		if len(list) > maxLength {
			maxLength = len(list)
		}
	}

	shuffled := [][]string{}
	for _, sourceList := range sourceTrackIDs {
		targetLength := len(sourceList)
		if pad {
			targetLength = maxLength
		}
		destList := make([]string, targetLength, targetLength)

		for i := range destList {
			modLen := i % len(sourceList)
			baseChars := (i / len(sourceList)) * len(sourceList)
			srcPos := i % len(sourceList)

			j := baseChars
			if modLen != 0 {
				j += rand.Intn(modLen)
			}

			if j == i {
				destList[i] = sourceList[srcPos]
			} else {
				destList[i] = destList[j]
				destList[j] = sourceList[srcPos]
			}
		}

		shuffled = append(shuffled, destList)
	}

	return shuffled
}

func padSourceTracks(sourceTrackIDs [][]string) [][]string {
	maxLength := 0
	for _, list := range sourceTrackIDs {
		if len(list) > maxLength {
			maxLength = len(list)
		}
	}

	padded := [][]string{}
	for _, sourceList := range sourceTrackIDs {
		newList := make([]string, maxLength)
		for i := range newList {
			newList[i] = sourceList[i%len(sourceList)]
		}
		padded = append(padded, newList)
	}

	return padded
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
