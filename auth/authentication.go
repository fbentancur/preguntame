package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/preguntame/preguntame-backend/models"
)

const JWT_SECRET = "Zoy Un Zecreto Muy Zecretozo"

type JwtHeaders struct {
	JWT string `header:"authorization"`
}

type UserClaims struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	jwt.StandardClaims
}

func NewTokenForUser(user models.User) (string, error) {
	claims := UserClaims {
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,

		StandardClaims: jwt.StandardClaims {
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(JWT_SECRET))
}

func DecodeUserToken(e echo.Context) (models.User, error) {
	headerString := e.Request().Header.Get("Authorization")

	if !strings.HasPrefix(headerString, "Bearer ") {
		return models.User{}, fmt.Errorf("Wrong prefix")
	}

	tokenString := strings.TrimPrefix(headerString, "Bearer ")

	claims := UserClaims {}
	_, err := jwt.ParseWithClaims(tokenString, &claims, getSecretForToken)
	if err != nil {
		return models.User{}, err
	}

	err = claims.Valid()
	if err != nil {
		return models.User{}, err
	}

	user := models.User {
		Id:    claims.Id,
		Name:  claims.Name,
		Email: claims.Email,
	}

	return user, nil
}

//TODO(fepalacios): Validate that the agl of the token is what we expected
func getSecretForToken(token *jwt.Token) (interface{}, error) {
	return []byte(JWT_SECRET), nil
}