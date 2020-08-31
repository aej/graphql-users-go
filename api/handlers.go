package api

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

func GraphqlHandler(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
