package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/youth-service/auth/pkg/setting"
	"github.com/youth-service/auth/pkg/store"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Passport struct {
	User          *User
	AccessToken   string
	AccessUuid    string
	AccessExpire  int64
	RefreshToken  string
	RefreshUuid   string
	RefreshExpire int64
}

func NewPassport(u *User) *Passport {
	return &Passport{
		User: u,
	}
}

// GenerateToken 토큰 생성
func (p *Passport) GenerateToken() error {
	accessUuid := uuid.New().String()
	accessExpire := time.Now().Add(setting.JWTSetting.AccessExpireTime).Unix()
	refreshUuid := accessUuid + "++" + strconv.Itoa(int(p.User.UID))
	refreshExpire := time.Now().Add(setting.JWTSetting.RefreshExpireTime).Unix()

	accessClaims := jwt.MapClaims{}
	accessClaims["authorized"] = true
	accessClaims["access_uuid"] = accessUuid
	accessClaims["uid"] = p.User.UID
	accessClaims["exp"] = accessExpire
	os.Setenv("ACCESS_SECRET", setting.JWTSetting.AccessSecret)

	var err error
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	p.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return err
	}
	p.AccessUuid = accessUuid
	p.AccessExpire = accessExpire

	refreshClaims := jwt.MapClaims{}
	refreshClaims["refresh_uuid"] = refreshUuid
	refreshClaims["uid"] = p.User.UID
	refreshClaims["exp"] = refreshExpire
	os.Setenv("REFRESH_SECRET", setting.JWTSetting.RefreshSecret)

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	p.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return err
	}
	p.RefreshUuid = refreshUuid
	p.RefreshExpire = refreshExpire

	return nil
}

// AddPassport save token data at redis
func (p *Passport) AddPassport() error {
	err := store.SaveMetadata(context.TODO(), p.AccessUuid, strconv.Itoa(int(p.User.UID)), p.AccessExpire)
	if err != nil {
		return err
	}
	err = store.SaveMetadata(context.TODO(), p.RefreshUuid, strconv.Itoa(int(p.User.UID)), p.RefreshExpire)
	if err != nil {
		return err
	}

	return nil
}

// FindPassportData find token data at redis
func (p *Passport) FindPassportData() (uint64, error) {
	uid, err := store.GetPassportData(context.TODO(), p.AccessUuid)
	if err != nil {
		return 0, err
	}

	if p.User.UID != uid {
		return 0, errors.New("unauthorized")
	}

	return uid, nil
}

// DeletePassportData delete token data at redis
func (p *Passport) DeletePassportData() error {

	deletedAt, err := store.DeleteMetadata(context.TODO(), p.AccessUuid)
	if err != nil {
		return err
	}

	deletedRt, err := store.DeleteMetadata(context.TODO(), p.RefreshUuid)
	if err != nil {
		return err
	}

	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}

	return nil
}

func (p *Passport) DeleteRefreshToken() error {

	_, err := store.DeleteMetadata(context.TODO(), p.AccessUuid)
	if err != nil {
		return err
	}

	deleted, err := store.DeleteMetadata(context.TODO(), p.RefreshUuid)
	if err != nil {
		return err
	}

	if deleted != 1 {
		return errors.New("something went wrong")
	}

	return nil
}

// MakeParamByHeader header에서 token 조회
func MakeParamByHeader(r *http.Request) (map[string]string, error) {
	param := map[string]string{}

	token := r.Header.Get("Authorization")
	strArr := strings.Split(token, " ")
	if len(strArr) == 2 {
		param["access_token"] = strArr[1]
		return param, nil
	}
	return nil, errors.New("cannot found access token")
}

// MakeParamByBody body에서 token 조회
func MakeParamByBody(c *gin.Context) (map[string]string, error) {
	param := map[string]string{}
	if err := c.ShouldBindJSON(&param); err != nil {
		return nil, err
	}
	return param, nil
}

// ParseAccessToken token 검증
func ParseAccessToken(access string) (*jwt.Token, error) {
	os.Setenv("ACCESS_SECRET", setting.JWTSetting.AccessSecret)
	token, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func ParseRefreshToken(refresh string) (*jwt.Token, error) {
	os.Setenv("REFRESH_SECRET", setting.JWTSetting.RefreshSecret)
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, err
}

func MakePassportByToken(token *jwt.Token) (*Passport, error) {

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("")
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, errors.New("")
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		return nil, err
	}

	return &Passport{
		User: &User{
			UID: userId,
		},
		AccessUuid: accessUuid,
	}, nil

}

func VerifyToken(token *jwt.Token) error {
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errors.New("not")
	}

	return nil
}

func MakePassport(claims jwt.MapClaims) (*Passport, error) {
	var accessUuid string
	var refreshUuid string
	var ok bool

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["uid"]), 10, 64)
	if err != nil {
		return nil, err
	}

	accessUuid, ok = claims["access_uuid"].(string)
	if ok {
		refreshUuid = fmt.Sprintf("%s++%d", accessUuid, userId)
		return &Passport{
			User: &User{
				UID: userId,
			},
			AccessUuid:  accessUuid,
			RefreshUuid: refreshUuid,
		}, nil
	}

	refreshUuid, ok = claims["refresh_uuid"].(string)
	if ok {
		accessUuid = strings.Split(refreshUuid, "++")[0]
		return &Passport{
			User: &User{
				UID: userId,
			},
			AccessUuid:  accessUuid,
			RefreshUuid: refreshUuid,
		}, nil
	}

	return nil, errors.New("cannot found tokens")
}
