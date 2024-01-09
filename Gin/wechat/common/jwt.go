package common

//进行token发放与验证
import (
	"ginEssential/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(30000 * 24 * 60 * time.Minute)
	claims := &Claims{
		UserId: user.ID, //token所属用户ID
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //token结束时间
			IssuedAt:  time.Now().Unix(),     //token发放的时间
			Issuer:    "oceanlearn.tech",     //发放事件
			Subject:   "user token",          //主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey) //用密钥制作token
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	//根据token填充Claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
