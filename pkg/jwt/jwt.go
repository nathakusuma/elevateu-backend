package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
)

type IJwt interface {
	Create(userID uuid.UUID, role enum.UserRole, isSubscribedBoost, isSubscribedChallenge bool) (string, error)
	Decode(tokenString string, claims *Claims) error
	Validate(token string) (ValidateJWTResponse, error)
}

type Claims struct {
	jwt.RegisteredClaims
	Role                  enum.UserRole `json:"role"`
	IsSubscribedBoost     bool          `json:"is_subscribed_boost"`
	IsSubscribedChallenge bool          `json:"is_subscribed_challenge"`
}

type ValidateJWTResponse struct {
	UserID                uuid.UUID
	Role                  enum.UserRole
	IsSubscribedBoost     bool
	IsSubscribedChallenge bool
}

type jwtStruct struct {
	exp    time.Duration
	secret []byte
}

func NewJwt(exp time.Duration, secret []byte) IJwt {
	return &jwtStruct{
		exp:    exp,
		secret: secret,
	}
}

func (j *jwtStruct) Create(userID uuid.UUID, role enum.UserRole, isSubscribedBoost,
	isSubscribedChallenge bool) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "elevateu-backend",
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role:                  role,
		IsSubscribedBoost:     isSubscribedBoost,
		IsSubscribedChallenge: isSubscribedChallenge,
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedJWT, err := unsignedJWT.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func (j *jwtStruct) Decode(tokenString string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

func (j *jwtStruct) Validate(token string) (ValidateJWTResponse, error) {
	var claims Claims
	err := j.Decode(token, &claims)
	if err != nil {
		return ValidateJWTResponse{}, errorpkg.ErrInvalidBearerToken()
	}

	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return ValidateJWTResponse{}, errorpkg.ErrInvalidBearerToken()
	}

	if expirationTime.Before(time.Now()) {
		return ValidateJWTResponse{}, errorpkg.ErrInvalidBearerToken()
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ValidateJWTResponse{}, errorpkg.ErrInvalidBearerToken()
	}

	return ValidateJWTResponse{
		UserID:                userID,
		Role:                  claims.Role,
		IsSubscribedBoost:     claims.IsSubscribedBoost,
		IsSubscribedChallenge: claims.IsSubscribedChallenge,
	}, nil
}
