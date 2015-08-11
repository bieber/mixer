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

package main

import (
	"github.com/bieber/mixer/mixerserver/context"
	"github.com/bieber/mixer/mixerserver/handlers"
	"github.com/bieber/mixer/mixerserver/middleware"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
	"path"
	"path/filepath"
)

func initRoutes(
	globalContext *context.GlobalContext,
	staticResourcesPath string,
) {
	r := mux.NewRouter().StrictSlash(true)
	globalContext.Router = r

	basicStack := alice.New(
		// This bottom instance of ErrorCatcher will catch any
		// failures in the logging or cleanup code, as a last resort.
		middleware.ErrorCatcher,
		middleware.ContextCleaner,
		middleware.Logger(globalContext),
		middleware.ErrorCatcher,
	)

	r.NotFoundHandler = basicStack.ThenFunc(handlers.FourOhFour)

	r.Handle("/", basicStack.Then(handlers.Index(globalContext)))
	r.Handle("/login/", basicStack.Then(handlers.Login(globalContext))).
		Name("login")

	staticHandler := func(subpath string) http.Handler {
		return basicStack.Then(
			http.StripPrefix(
				path.Join("/static/", subpath),
				http.FileServer(
					http.Dir(
						filepath.Join(staticResourcesPath, subpath),
					),
				),
			),
		)
	}
	s := r.PathPrefix("/static").Subrouter()

	s.Handle("/js/{rest:.*}", staticHandler("/js"))
	s.Handle("/css/{rest:.*}", staticHandler("/css"))
	s.Handle("/img/{rest:.*}", staticHandler("/img"))
}
