package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mubarik/EVENT_MANBAGEMENT_SYSTEM/infra"
)

// ✅ Correct Claims for jwt/v5
type Claims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

//
// =========================
// 🔐 ACCESS TOKEN MIDDLEWARE
// =========================
//

func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		// 1. Check header exists
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Missing Authorization header",
				"is_success": false,
			})
			return
		}

		// 2. Check Bearer format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Invalid Authorization header format",
				"is_success": false,
			})
			return
		}

		// 3. Extract token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		secret := []byte(infra.Configuration.Access_jwt_Token)

		claims := &Claims{}

		// 4. Parse token with claims
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {

			// ✅ Security: check signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return secret, nil
		})

		// 5. Validate token
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Unauthenticated",
				"is_success": false,
			})
			return
		}

		// 6. Check expiry (v5 style)
		if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Access token expired",
				"is_success": false,
			})
			return
		}

		// 7. Validate payload
		if claims.Sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Invalid token payload",
				"is_success": false,
			})
			return
		}

		// 8. Set context
		c.Set("email", claims.Sub)
		c.Set("role", claims.Role)

		// 9. Log user
		slog.Info("Authenticated user", "email", claims.Sub, "role", claims.Role)

		c.Next()
	}
}

//
// =========================
// 🔁 REFRESH TOKEN MIDDLEWARE
// =========================
//

func RefreshAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		// 1. Check header
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Missing Authorization header",
				"is_success": false,
			})
			return
		}

		// 2. Check Bearer format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Invalid Authorization header format",
				"is_success": false,
			})
			return
		}

		// 3. Extract token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}

		// 4. Parse refresh token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(infra.Configuration.Refresh_jwt_token), nil
		})

		// 5. Validate token
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Unauthorized access",
				"is_success": false,
			})
			return
		}

		// 6. Check expiry
		if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Refresh token expired",
				"is_success": false,
			})
			return
		}

		// 7. Validate payload
		if claims.Sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":    "Invalid token payload",
				"is_success": false,
			})
			return
		}

		// 8. Set context
		c.Set("user_email", claims.Sub)

		slog.Info("Refresh token valid", "email", claims.Sub)

		c.Next()
	}
}
