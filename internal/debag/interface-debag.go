package debag

import (
	"net/http"
	"time"
)

type configDebagI interface {
	GetDebagAddr() string
	GetDebagReadTimeout() time.Duration
	GetDebagWriteTimeout() time.Duration
	GetConfigInfo(w http.ResponseWriter, r *http.Request)
}
