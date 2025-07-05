package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtUtils JWTUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := extractBearerToken(c.GetHeader("Authorization"))
        if tokenStr == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }

        userID, accessJTI, err := jwtUtils.ValidateJWT(tokenStr)
	
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

        c.Set("userID", *userID)
        c.Set("accessJTI", accessJTI)

        c.Next()
    }
}

func extractBearerToken(authHeader string) string {
    const prefix = "Bearer "

    if strings.HasPrefix(authHeader, prefix) {
        return strings.TrimPrefix(authHeader, prefix)
    }

    return ""
}
