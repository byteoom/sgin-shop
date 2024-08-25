package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey = "your-secret-key"
)

// GenerateToken 生成 JWT token
func GenerateToken(userID string) (string, error) {
	// 定义 JWT 的有效期限
	expirationTime := time.Now().Add(24 * time.Hour) // 设置为 24 小时有效期，可根据需求调整

	// 创建 token 的声明部分
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	// 使用 HS256 算法进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥对 token 进行签名，生成字符串格式的 token
	tokenString, err := token.SignedString([]byte(secretKey)) // 使用与验证时相同的密钥
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil // 使用与生成 token 时相同的密钥
	})
	if err != nil {
		return nil, err
	}

	// 获取 token 中的声明部分
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

// 解析token返回user_id
func ParseTokenGetUserID(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", err
	}

	return userID, nil
}

// 生成验证码
func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999)
	return strconv.Itoa(code)
}

// SignBody 签名
func SignBody(body, secretKey []byte) string {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// 数组转换成json字符串
func ArrayToJsonString(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	jsonBytes, _ := json.Marshal(arr)
	return string(jsonBytes)
}

func MapGetString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func MapGetFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

// GetFileMd5
func GetFileMd5(fileinfo multipart.File) (string, error) {
	md5h := md5.New()
	if _, err := io.Copy(md5h, fileinfo); err != nil {
		return "", err
	}
	return hex.EncodeToString(md5h.Sum(nil)), nil
}

// GenerateOrderID generates a unique order ID based on the current date and time including nanoseconds.
func GenerateOrderID() string {
	// Set the seed for random number generation
	rand.Seed(time.Now().UnixNano())

	// Get the current date and time including nanoseconds
	now := time.Now()
	dateStr := now.Format("20060102150405") // Format as YYYYMMDDHHMMSS
	nanoStr := fmt.Sprintf("%09d", now.Nanosecond())

	// Generate a random 4-digit number
	randomNum := rand.Intn(10000)
	randomStr := fmt.Sprintf("%04d", randomNum)

	// Combine the date, time, nanoseconds, and random number to form the order ID
	orderID := dateStr + nanoStr + randomStr
	return orderID
}
