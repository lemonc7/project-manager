package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lemonc7/project-manager/config"
)

type JWT struct {
	secret   []byte
	duration time.Duration
}

func NewJWT(cfg config.JWTConfig) *JWT {
	return &JWT{
		secret:   []byte(cfg.Secret),
		duration: cfg.Duration,
	}
}

type UserClaims struct {
	Username string   `json:"username"`
	UserID   int32    `json:"userID"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

func (j *JWT) Generate(username string, id int32, roles []string) (string, error) {
	claims := UserClaims{
		Username: username,
		UserID:   id,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

func (j *JWT) ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("不支持的签名方法")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("解码user信息: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token已失效")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("无效token信息: %w", err)
	}

	return claims, nil
}
