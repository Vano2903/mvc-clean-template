package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	_DEFAULT_EXPIRATION_TIME = 15 * time.Minute
)

type JWThandler struct {
	secret         string
	issuer         string
	expirationTime time.Duration
}

type JWTClaims struct {
	jwt.StandardClaims
	UserId    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`
}

func NewJWThandler(secret, issuer string, expirationTime ...time.Duration) *JWThandler {
	expTime := _DEFAULT_EXPIRATION_TIME
	if len(expirationTime) > 0 {
		expTime = expirationTime[0]
	}
	return &JWThandler{
		secret:         secret,
		issuer:         issuer,
		expirationTime: expTime,
	}
}

func (j *JWThandler) SigningKey() []byte {
	return []byte(j.secret)
}

// GenerateToken generates a new token with the given claims
// that lasts for 15 minutes
func (j *JWThandler) GenerateToken(userId int, userEmail string, userRole string) (string, error) {
	claims := JWTClaims{
		UserId:    userId,
		UserEmail: userEmail,
		UserRole:  userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.expirationTime).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ValidateToken validates the given token and returns the claims
func (j *JWThandler) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}

func (j *JWThandler) IsTokenExpired(tokenString string) (bool, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return false, err
	}
	return claims.ExpiresAt < time.Now().Unix(), nil
}
