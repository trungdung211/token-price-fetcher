package app

import (
	"net/http"

	"github.com/trungdung211/token-price-fetcher/pkg/request"

	"github.com/gin-gonic/gin"
)

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		res := request.NewResponse()
		status := http.StatusInternalServerError

		for _, err := range c.Errors {
			if err != nil {
				if appErr, ok := err.Err.(*request.AppError); ok {
					res.ErrorCode = appErr.Code
					res.ErrorMessage = appErr.Msg
					if appErr.StatusCode > 0 {
						status = appErr.StatusCode
					}

					c.JSON(status, res)
					return
				} else {
					res.ErrorMessage = http.StatusText(status)
					c.JSON(status, res)
					return
				}
			}
		}
	}
}
