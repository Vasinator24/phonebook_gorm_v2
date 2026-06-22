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

type contextKey string

const userIDContextKey contextKey = "user_id"

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthServeConfig struct {
	DB            *gorm.DB
	JWTSecret     []byte
	TokenDuration time.Duration
	CookieName    string
	CookieSecure  bool
}

type AuthServe struct {
	config AuthServeConfig
}

func NewAuthServeConfig(dbConn *gorm.DB) AuthServeConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set")
	}

	hours, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_HOURS"))
	if err != nil || hours <= 0 {
		hours = 24
	}

	return AuthServeConfig{
		DB:            dbConn,
		JWTSecret:     []byte(secret),
		TokenDuration: time.Duration(hours) * time.Hour,
		CookieName:    "token",
		CookieSecure:  os.Getenv("COOKIE_SECURE") == "true",
	}
}

func NewAuthServe(config AuthServeConfig) *AuthServe {
	return &AuthServe{config: config}
}

func (a *AuthServe) GenerateToken(userID uint) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(a.config.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.config.JWTSecret)
}

func (a *AuthServe) ParseToken(tokenValue string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return a.config.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (a *AuthServe) TokenHash(tokenValue string) string {
	sum := sha256.Sum256([]byte(tokenValue))
	return hex.EncodeToString(sum[:])
}

func (a *AuthServe) IsTokenBlacklisted(tokenValue string) (bool, error) {
	var count int64
	err := a.config.DB.Model(&appdb.BlacklistedToken{}).
		Where("token_hash = ? AND expires_at > ?", a.TokenHash(tokenValue), time.Now().Unix()).
		Count(&count).Error

	return count > 0, err
}

func (a *AuthServe) SetTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     a.config.CookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(a.config.TokenDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   a.config.CookieSecure,
	})
}

func (a *AuthServe) ClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     a.config.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   a.config.CookieSecure,
	})
}

func (a *AuthServe) UserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uint)
	return userID, ok
}

func (a *AuthServe) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(a.config.CookieName)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := a.ParseToken(cookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		isBlacklisted, err := a.IsTokenBlacklisted(cookie.Value)
		if err != nil || isBlacklisted {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthServe) TokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(a.config.CookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
