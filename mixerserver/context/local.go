package context

import (
	"github.com/bieber/logger"
	"github.com/bieber/mixer/mixerserver/spotify"
	"net/http"
	"sync"
)

// LocalContext stores context relevant to a single request.  It
// should be both written to and read from by middleware, and read
// from by controllers.
type LocalContext struct {
	Logger     *logger.Logger
	AuthTokens spotify.AuthTokens
}

var localMutex = sync.Mutex{}
var localContexts = make(map[*http.Request]*LocalContext)

// Get returns the LocalContext for a given request, or creates one if
// it doesn't already exist.
func Get(request *http.Request) *LocalContext {
	localMutex.Lock()
	defer localMutex.Unlock()

	if c, ok := localContexts[request]; ok {
		return c
	}
	c := &LocalContext{}
	localContexts[request] = c
	return c
}

// Clear removes the LocalContext entry for a request after it's
// finished.
func Clear(request *http.Request) {
	localMutex.Lock()
	defer localMutex.Unlock()
	delete(localContexts, request)
}
