package utils

import (
	"errors"
	"fmt"
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/Xacnio/aptms-backend/pkg/configs"
	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
	"strconv"
	"time"
)

func CreateAccessToken(userid uint) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 24).Unix()
	td.AccessUuid = uuid.NewV4().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS384, atClaims)
	td.AccessToken, err = at.SignedString([]byte(configs.Get("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func CreateAuthData(userid uint, td *models.TokenDetails) error {
	// Store the token data in Redis
	// TODO: Redis

	//at := time.Unix(td.AtExpires, 0)
	//rt := time.Unix(td.RtExpires, 0)
	//now := time.Now()
	//

	//errAccess := database.RedisClient.Set(td.AccessUuid, strconv.FormatUint(uint64(userid), 10), at.Sub(now)).Err()
	//if errAccess != nil {
	//	return errAccess
	//}
	//errRefresh := database.RedisClient.Set(td.RefreshUuid, strconv.FormatUint(uint64(userid), 10), rt.Sub(now)).Err()
	//if errRefresh != nil {
	//	return errRefresh
	//}
	return nil
}

func VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.Get("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenPayload(token *jwt.Token) (*models.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("Access details could not be validated")
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     uint(userId),
		}, nil
	}
	return nil, errors.New("Access details could not be validated")
}

func IsTokenValid(token *jwt.Token) bool {
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false
	}
	return true
}

func GetTokenMapClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

func RefreshToken(token *jwt.Token) (*models.TokenDetails, error) {
	claims, ok := GetTokenMapClaims(token) // the token claims should conform to MapClaims
	if ok && token.Valid {
		//refreshUuid, ok := claims["refresh_uuid"].(string) // convert the interface to string
		//if !ok {
		//	return nil, errors.New("Refresh uuid is none")
		//}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		//deleted, delErr := utils.DeleteAuthData(refreshUuid)
		//if delErr != nil || deleted == 0 {
		//	return c.Status(fiber.StatusUnauthorized).JSON("unauthorized")
		//}
		ts, createErr := CreateAccessToken(uint(userId))
		if createErr != nil {
			return nil, createErr
		}
		saveErr := CreateAuthData(uint(userId), ts)
		if saveErr != nil {
			return nil, saveErr
		}
		return ts, nil
	} else {
		return nil, errors.New("refresh token expired or revoked")
	}
}
