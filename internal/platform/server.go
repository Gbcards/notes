package platform

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle, router *gin.Engine) *http.Server {
	server := &http.Server{Addr: ":8080", Handler: router}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}
