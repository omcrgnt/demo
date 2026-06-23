package http

import (
	srvhttp "github.com/omcrgnt/srv-http"

	"github.com/omcrgnt/builder"
)

// Server is the items HTTP server resource; config shape is [srvhttp.Config].
type Server srvhttp.Config[*API]

func (s *Server) BuildConfig() (builder.Builder, error) {
	cfg := srvhttp.Config[*API]{}
	return &cfg, nil
}
