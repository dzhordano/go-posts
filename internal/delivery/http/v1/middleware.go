package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"

	userCtx  = "userId"
	adminCtx = "adminId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, id)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(adminCtx, id)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.ValidateToken(headerParts[1])
}

func (h *Handler) getUserId(c *gin.Context) (uint, error) {
	return h.getIdFromContext(c, userCtx)
}

func (h *Handler) getIdFromContext(c *gin.Context, context string) (uint, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return 0, errors.New("user context not found")
	}

	idString, ok := idFromCtx.(string)
	if !ok {
		return 0, errors.New("user context of invalid type")
	}

	id, err := strconv.ParseInt(idString, 16, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
