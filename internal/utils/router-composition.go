package utils

import "net/http"

type Hanlde struct {
	Pattern string
	Handler http.Handler
}

func RouterComposition(hanldes ...Hanlde) *http.ServeMux {
	multiplexer := http.NewServeMux()
	for _, v := range hanldes {
		multiplexer.Handle(v.Pattern, v.Handler)
	}

	return multiplexer
}
