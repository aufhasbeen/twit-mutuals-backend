package routing

import (
	"net/http"
)

type sessionHeader struct {
	oauthPub  string
	oauthPriv string
	twitUser  int64
}

func (s *sessionHeader) startSession(header http.Header) {
	s.oauthPub = ""  // header["oath"][0]
	s.oauthPriv = "" // header["oathSecret"][0]
	s.twitUser = 0   //, _ =  strconv.ParseInt(header["userID"][0], 10, 64) // may be variable in url
}

// TODO: use sessions for oauth messages not config. config is universal and persists across user sessions!!!!!
// sidenote: they may also be changed async under the current user request leading to issues with twitter api requests
// func (s *sessionHeader) startSessionTesting(header http.Header) {
// 	config := conf.GetConfig().Authentication.Developer
// 	s.oauthPub = config.Oauth
// 	s.oauthPriv = config.OauthSecret
// 	s.twitUser, _ = strconv.ParseInt(header["User-Id"][0], 10, 64) // may be variable in url
// }

// instead of writing this check if the header preserves into a response.
// then paginate the response to the unfollowers route. for now do the first ten as a test.
