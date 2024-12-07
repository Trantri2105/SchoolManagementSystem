package middleware

import (
	"SchoolManagement/dto/response"
	"SchoolManagement/utils"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

type AuthMiddleware interface {
	ValidateAndExtractJwt() gin.HandlerFunc
	CheckUserAuthorities(ctx context.Context, r ...string) error
}

const (
	JWTClaimsContextKey = "JWTClaimsContextKey"
)

type authMiddleware struct {
	jwtService utils.JwtUtils
}

func (a *authMiddleware) CheckUserAuthorities(ctx context.Context, r ...string) error {
	claims := ctx.Value(JWTClaimsContextKey).(jwt.MapClaims)
	var isAuthorized bool
	for _, ro := range r {
		if claims["role"].(string) == ro {
			isAuthorized = true
			break
		}
	}
	if !isAuthorized {
		return errors.New("unauthorized")
	}
	return nil
}

func (a *authMiddleware) ValidateAndExtractJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Message{Error: "require access token"})
			return
		}
		header := strings.Fields(authHeader)
		if len(header) != 2 && header[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Message{Error: "wrong access token format"})
			return
		}
		accessToken := header[1]
		claims, err := a.jwtService.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Message{Error: err.Error()})
		}

		ctx := context.WithValue(c.Request.Context(), JWTClaimsContextKey, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func NewAuthMiddleware(jwtService utils.JwtUtils) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
