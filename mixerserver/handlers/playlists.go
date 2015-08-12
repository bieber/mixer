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
	"github.com/bieber/mixer/mixerserver/context"
	"net/http"
)

// Playlists fetches and returns a list of the user's playlists as
// JSON.
func Playlists(w http.ResponseWriter, r *http.Request) {
	localContext := context.Get(r)
	localContext.Logger.Println(localContext.AuthTokens)
	_ = localContext
}
