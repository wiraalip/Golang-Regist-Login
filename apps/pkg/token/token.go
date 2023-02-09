package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type PayloadToken struct {
	AuthId  int
	Expired time.Time
}

const SecretKey = "HyVQNmB3SMjwYvL4Tqh90N7tD6ccoF8t"

func GenerateToken(tok *PayloadToken) (string, error) {
	tok.Expired = time.Now().Add(10 * 60 * time.Second)
	claims := jwt.MapClaims{
		"payload": tok,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ValidateToken(tokString string) (*PayloadToken, error) {
	tok, err := jwt.Parse(tokString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("unauthorized1")
	}
	payload := claims["payload"]
	var payloadToken = PayloadToken{}
	payloadByte, _ := json.Marshal(payload)
	err = json.Unmarshal(payloadByte, &payloadToken)
	if err != nil {
		return nil, err
	}

	// payloadToken := payload.(PayloadToken)
	// if !ok {
	// 	return nil, errors.New("Invalid Payload Type")
	// }

	return &payloadToken, nil
}
