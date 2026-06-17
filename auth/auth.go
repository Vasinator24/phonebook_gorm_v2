package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	appdb "phonebook_gorm/db"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const CookieName = "token"

type contextKey string

const userIDContextKey contextKey = "user_id"

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set")
	}

	return []byte(secret)
}

func tokenDuration() time.Duration {
	hours, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_HOURS"))
	if err != nil || hours <= 0 {
		hours = 24
	}

	return time.Duration(hours) * time.Hour
}

func cookieSecure() bool {
	return os.Getenv("COOKIE_SECURE") == "true"
}

func GenerateToken(userID uint) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration())),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func ParseToken(tokenValue string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return jwtSecret(), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func TokenHash(tokenValue string) string {
	sum := sha256.Sum256([]byte(tokenValue))
	return hex.EncodeToString(sum[:])
}

func IsTokenBlacklisted(dbConn *gorm.DB, tokenValue string) (bool, error) {
	var count int64
	err := dbConn.Model(&appdb.BlacklistedToken{}).
		Where("token_hash = ? AND expires_at > ?", TokenHash(tokenValue), time.Now().Unix()).
		Count(&count).Error

	return count > 0, err
}

func SetTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(tokenDuration().Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   cookieSecure(),
	})
}

func ClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   cookieSecure(),
	})
}

func UserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uint)
	return userID, ok
}

func AuthMiddleware(dbConn *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(CookieName)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := ParseToken(cookie.Value)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			isBlacklisted, err := IsTokenBlacklisted(dbConn, cookie.Value)
			if err != nil || isBlacklisted {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
