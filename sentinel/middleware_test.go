package sentinel

import (
	"net/http"

	. "gopkg.in/check.v1"
)

type MiddlewareSuite struct{}

var _ = Suite(&MiddlewareSuite{})

func noop(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "foobar")
}
