package v1

import (
	"github.com/gin-gonic/gin"
)

type dataResponse struct {
	Data interface{} `json:"data"`
}

type response struct {
	Message string `json:"message"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, msg)
}
