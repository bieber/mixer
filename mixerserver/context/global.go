package context

import (
	"github.com/gorilla/mux"
	"html/template"
	"io"
)

// GlobalContext stores data relevant to the entire server process.
// Only a single instance need exist, and controllers should not write
// to it.
type GlobalContext struct {
	Router    *mux.Router
	Templates struct {
		Index *template.Template
	}
	LogOut  io.Writer
	Spotify struct {
		ClientID     string
		ClientSecret string
	}
}
