package authservice

import (
	"fmt"
	"strings"
	"time"

	utils "github.com/aghaghiamh/gocast/QAGame/pkg"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	config AuthConfig
}

type AuthConfig struct {
	SignKey              string
	AccessSubject        string
	RefreshSubject       string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func New(authConf AuthConfig) Service {
	return Service{
		config: authConf,
	}
}

type Claims struct {
	Subject string `json:"subject"`
	UserID  uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func (c *Claims) Valid() {
	return
}

func (s *Service) CreateAccessToken(userID uint) (string, error) {
	return createToken(userID, s.config.AccessSubject, []byte(s.config.SignKey), s.config.AccessTokenDuration)
}

func (s *Service) CreateRefreshToken(userID uint) (string, error) {
	return createToken(userID, s.config.RefreshSubject, []byte(s.config.SignKey), s.config.RefreshTokenDuration)
}

func createToken(userID uint, subject string, signKey []byte, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&Claims{
			Subject: subject,
			UserID:  userID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: &jwt.NumericDate{time.Now().Add(ttl)},
			},
		})

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", utils.RichErr{
			Code:    utils.GeneralServerErr,
			Message: fmt.Sprintf("Couldn't sign JWT token: %s", err),
		}
	}

	return tokenString, nil
}

func (s *Service) VerifyToken(bearerToken string) (*Claims, error) {

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("Authorization header must be in the format: Bearer {token}")
	}

	token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.SignKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, utils.RichErr{
		Code: utils.TokenInvalidErr,
	}
}
