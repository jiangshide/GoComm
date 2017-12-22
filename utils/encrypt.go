package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"crypto/rand"
	"encoding/base64"
	r "math/rand"
	"time"
	"github.com/jiangshide/jwt-go"
	"errors"
)

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetGuid() (string, error) {
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	return Md5(base64.URLEncoding.EncodeToString(b)), err
}

var alphaNum = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)

// RandomCreateBytes generate random []byte by specify chars.
func RandomCreateBytes(n int, alphabets ...byte) []byte {
	if len(alphabets) == 0 {
		alphabets = alphaNum
	}
	var bytes = make([]byte, n)
	var randBy bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randBy = true
	}
	for i, b := range bytes {
		if randBy {
			bytes[i] = alphabets[r.Intn(len(alphabets))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return bytes
}

func Token(m map[string]interface{}, name string, time int64) (token string, err error) {
	claims := make(jwt.MapClaims)
	for k, v := range m {
		claims[k] = v
	}
	claims["exp"] = time
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tk.SignedString([]byte(name))
	return
}

func UnToken(tokenStr, name string) (m map[string]interface{}, errs error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(name), nil
	})
	errs = err
	if errs == nil {
		if token.Valid {
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				m = claims
			} else {
				errs = errors.New("no Permission")
			}
		} else {
			errs = errors.New("Token invalid:" + tokenStr)
		}
	}
	return
}
