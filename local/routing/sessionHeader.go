package routing

import (
	"net/http"
	"strconv"
)

type sessionHeader struct {
	oauthPub  string
	oauthPriv string
	twitUser  int64
}

func (s *sessionHeader) startSession(header http.Header) {
	s.oauthPub = header["oath"][0]
	s.oauthPriv = header["oathSecret"][0]
	s.twitUser, _ = strconv.ParseInt(header["userID"][0], 10, 64) // may be variable in url

}

// instead of writing this check if the header preserves into a response.
// then paginate the response to the unfollowers route. for now do the first ten as a test.
