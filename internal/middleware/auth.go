package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/pkg/jwt"
	"github.com/edustack/go-boilerplate/pkg/tokenstore"
	"github.com/gin-gonic/gin"
)

type AuthDeps struct {
	JWTUtils   *jwt.JWTUtils
	TokenStore *tokenstore.TokenStore
}

func AuthMiddleware(cfg *config.Config, deps *AuthDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]

		var claims map[string]interface{}
		if deps != nil && deps.JWTUtils != nil {
			jwtClaims, err := deps.JWTUtils.VerifyAccessToken(tokenString)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "invalid or expired token",
				})
				return
			}
			claims = jwtClaims
		} else {
			claims = verifyJWTDirect(tokenString, cfg.JWTSecret)
			if claims == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "invalid or expired token",
				})
				return
			}
		}

		userID, _ := claims["user_id"].(string)
		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)
		role, _ := claims["role"].(string)

		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid user ID claim",
			})
			return
		}

		if deps != nil && deps.TokenStore != nil {
			storedToken, _ := deps.TokenStore.GetAccessToken(c.Request.Context(), userID)
			if storedToken != "" && storedToken != tokenString {
				slog.Warn("token mismatch - possible reuse", "user_id", userID)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "token has been revoked",
				})
				return
			}
		}

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("name", name)
		c.Set("role", role)

		c.Next()
	}
}

func verifyJWTDirect(tokenString, secret string) map[string]interface{} {
	jwtUtils := jwt.NewJWTUtils(secret, "", 1, 1)
	parsedToken, err := jwtUtils.VerifyAccessToken(tokenString)
	if err != nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range parsedToken {
		result[k] = v
	}
	return result
}
