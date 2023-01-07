package rest

import (
	"net/http"
)

func ChainMiddlewares(hdl http.Handler, mds ...func(http.Handler) http.Handler) http.Handler {
	for _, md := range mds {
		hdl = md(hdl)
	}

	return hdl
}
