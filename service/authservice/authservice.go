package authservice

import (
	"strings"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
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
	const op = "authservice.createToken"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&Claims{
			Subject: subject,
			UserID:  userID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: &jwt.NumericDate{time.Now().Add(ttl)},
			},
		})

	tokenString, signErr := token.SignedString(signKey)
	if signErr != nil {

		return "", richerr.New(op).
			WithError(signErr).
			WithCode(richerr.ErrServer).
			WithMessage("Couldn't sign JWT Token")
	}

	return tokenString, nil
}

func (s *Service) VerifyToken(bearerToken string) (*Claims, error) {
	const op = "authservice.VerifyToken"

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {

		return nil, richerr.New(op).
			WithCode(richerr.ErrInvalidToken).
			WithMessage("Authorization header must be in the format: Bearer {token}").
			WithMetadata(map[string]interface{}{"token": bearerToken})
	}

	token, signErr := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, richerr.New(op).
				WithCode(richerr.ErrInvalidToken).
				WithMessage("unexpected signing method").
				WithMetadata(map[string]interface{}{"sign-method": token.Header["alg"]})
		}

		return []byte(s.config.SignKey), nil
	})

	if signErr != nil {

		return nil, signErr
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {

		return claims, nil
	}

	return nil, richerr.New(op).
		WithCode(richerr.ErrInvalidToken).
		WithMessage("Invalid Claims").
		WithMetadata(map[string]interface{}{"claims": token.Claims})
}
