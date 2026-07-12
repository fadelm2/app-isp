package util

import (
	"context"
	"errors"
	"fmt"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenUtil struct {
	SecretKey []byte
}

func NewTokenUtil(secretKey string) *TokenUtil {
	return &TokenUtil{
		SecretKey: []byte(secretKey),
	}
}

func (t TokenUtil) GenerateJWT(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        user.ID,
		"role":      user.RoleID,
		"name":      user.Name,
		"region_id": user.RegionId,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(30 * time.Minute).Unix(),
	})
	return token.SignedString(t.SecretKey)
}

func (t TokenUtil) getToken(ctx *fiber.Ctx) (*jwt.Token, error) {
	tokenString, err := t.getTokenFromRequest(ctx)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.SecretKey, nil
	})
	return token, err
}

func (t TokenUtil) getTokenFromRequest(ctx *fiber.Ctx) (string, error) {
	var tokenString string
	authorization := ctx.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if ctx.Cookies("token") != "" {
		tokenString = ctx.Cookies("token")
	}

	if tokenString == "" {
		return "", fiber.ErrUnauthorized
	}
	return tokenString, nil
}

func (t TokenUtil) ValidateJWT(context *fiber.Ctx) error {
	token, err := t.getToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func (t TokenUtil) ValidateAdminRoleJWT(context *fiber.Ctx) error {
	token, err := t.getToken(context)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	roleVal, ok := claims["role"].(string)
	if !ok {
		return fiber.ErrUnauthorized
	}
	//userRole := int(roleVal)
	if ok && token.Valid && roleVal == "1" {
		return nil
	}
	return errors.New("invalid admin token provided")
}

func (t TokenUtil) ValidateCustomerRoleJWT(context *fiber.Ctx) error {
	token, err := t.getToken(context)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)
	if ok && token.Valid && userRole == "2" || userRole == "1" {
		return nil
	}

	return errors.New("invalid customer or admin token provided")
}

func (t TokenUtil) ValidateSuperAdminRoleJWT(context *fiber.Ctx) error {
	token, err := t.getToken(context)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)
	if ok && token.Valid && userRole == "99" {
		return nil
	}

	return errors.New("invalid superadmin token provided")
}

func (t TokenUtil) ValidateDriverRoleJWT(context *fiber.Ctx) error {
	token, err := t.getToken(context)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)
	if ok && token.Valid && userRole == "4" {
		return nil
	}

	return errors.New("invalid driver token provided")
}

func (t TokenUtil) ParseToken(ctx context.Context, jwtToken string) (*model.Auth, error) {
	// 🔹 Buang "Bearer "
	if strings.HasPrefix(jwtToken, "Bearer ") {
		jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")
	}
	//
	//token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
	//	return []byte(t.SecretKey), nil
	//})
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return t.SecretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}

	expire := claims["exp"].(float64)
	if int64(expire) < time.Now().Unix() {
		return nil, fiber.ErrUnauthorized
	}

	id := claims["id"].(string)
	roleStr, ok := claims["role"].(string)

	auth := &model.Auth{
		ID:     id,
		RoleID: roleStr,
	}
	return auth, nil

}
