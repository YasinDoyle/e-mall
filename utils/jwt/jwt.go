package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/YasinDoyle/e-mall/consts"
)

var jwtSecret = []byte("YasinDoyle")

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT
func GenerateToken(id uint, username string) (accessToken, refreshToken string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(consts.AccessTokenExpireDuration)
	rtExpireTime := nowTime.Add(consts.RefreshTokenExpireDuration)

	claims := Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExpireTime.Unix(),
			Issuer:    "mall",
		},
	}).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err

}

func ParseRefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	accessClaim, err := ParseToken(aToken)
	if err != nil {
		return
	}

	refreshClaim, err := ParseToken(rToken)
	if err != nil {
		return
	}

	if accessClaim.ExpiresAt > time.Now().Unix() {
		// 如果 access_token 没过期,每一次请求都刷新 refresh_token 和 access_token
		return GenerateToken(accessClaim.ID, accessClaim.Username)
	}

	if refreshClaim.ExpiresAt > time.Now().Unix() {
		// 如果 access_token 过期了,但是 refresh_token 没过期, 刷新 refresh_token 和 access_token
		return GenerateToken(accessClaim.ID, accessClaim.Username)
	}

	// 如果两者都过期了,重新登陆
	return "", "", errors.New("身份过期，重新登陆")
}

type EmailClaims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.StandardClaims
}

func GenerateEmailToken(userID, Operation uint, email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(15 * time.Minute)
	claims := EmailClaims{
		UserID:        userID,
		Email:         email,
		Password:      password,
		OperationType: Operation,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cmall",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
