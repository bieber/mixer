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
	"fmt"
	"github.com/bieber/mixer/mixerserver/context"
	"github.com/bieber/mixer/mixerserver/crypto"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	for _, flag := range os.Args {
		if flag == "-h" || flag == "--help" {
			key, err := crypto.GenerateAESKey()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(key)
			return
		}
	}

	viper.SetDefault("port", 80)

	viper.BindEnv("port")
	viper.BindEnv("static_path")
	viper.BindEnv("spotify_client_id")
	viper.BindEnv("spotify_client_secret")
	viper.BindEnv("token_key")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Couldn't load config file: %s", err.Error())
	}

	crypto.SetAESKey(viper.GetString("token_key"))

	globalContext := &context.GlobalContext{}
	globalContext.Spotify.ClientID = viper.GetString("spotify_client_id")
	globalContext.Spotify.ClientSecret = viper.GetString(
		"spotify_client_secret",
	)

	initRoutes(globalContext, viper.GetString("static_path"))

	err = initTemplates(globalContext, viper.GetString("static_path"))
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", globalContext.Router)

	log.Printf("Starting server on port %d...", viper.GetInt("port"))
	log.Fatal(
		http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), nil),
	)
}
