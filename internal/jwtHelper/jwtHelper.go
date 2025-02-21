package jwtHelper

import (
	"crypto/rsa"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//go:embed public-key.pem
var publicKey []byte

// jwtHelper now receives a pointer to the Application struct
type JWTService struct { // No interface, just a struct
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTService(secretKey string) (*JWTService, error) {
	privateKey, err := parsePrivateKey(secretKey)
	if err != nil {
		return nil, err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	return &JWTService{privateKey: privateKey, publicKey: pubKey}, nil
}

// Load RSA private key from an environment variable
func parsePrivateKey(pemKey string) (*rsa.PrivateKey, error) {
	if pemKey == "" {
		return nil, errors.New("RSA_PRIVATE_KEY environment variable not set")
	}

	pemKey = fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----\n", pemKey)

	// Decode the PEM block containing the private key
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// Parse the RSA private key from the PEM block
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pemKey))
	if err != nil {
		return nil, fmt.Errorf("could not parse RSA private key: %v", err)
	}

	return privKey, nil
}

func (h *JWTService) CreateJwtCookie(userId int64, username string, rememberMe bool) (*http.Cookie, error) {
	daysToExpire := 1
	if rememberMe {
		daysToExpire = 28
	}

	expiry := time.Now().Add(time.Hour * 24 * time.Duration(daysToExpire))

	jwtToken, err := h.createSignedJWT(userId, username, expiry)

	if err != nil {
		return nil, err
	}

	cookie, err := h.createCookieFromJWT(jwtToken, rememberMe, expiry)

	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func (h *JWTService) ClearCookie() (*http.Cookie, error) {
	return h.createCookieFromJWT("", false, time.Now())
}

func (h *JWTService) VerifyJWT(signedToken string) error {
	// Parse the public key
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return fmt.Errorf("error parsing public key: %v", err)
	}

	// Verify the token
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return fmt.Errorf("error verifying token: %v", err)
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}
	return nil
}

// Function to create and sign a JWT
func (h *JWTService) createSignedJWT(userId int64, username string, expiry time.Time) (string, error) {
	// Create a new token object with claims
	claims := jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      expiry.Unix(),
	}

	// Create token with expiry:
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the private key
	signedToken, err := token.SignedString(h.privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing the token: %v", err)
	}

	return signedToken, nil
}

func (h *JWTService) createCookieFromJWT(signedToken string, rememberMe bool, expiry time.Time) (*http.Cookie, error) {

	isDev := os.Getenv("ENV") == "local"

	// Create the JWT cookie
	jwtCookie := &http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   !isDev,
		SameSite: func() http.SameSite {
			if !isDev {
				// TODO: Change this to Strict once we are in PROD (and determine a way to handle staging)
				return http.SameSiteNoneMode
			}
			return http.SameSiteLaxMode
		}(),
	}

	if rememberMe {
		jwtCookie.Expires = expiry
	}

	return jwtCookie, nil
}

// ParseAndVerifyJWT parses a JWT token and verifies it with the private key stored in the environment
func (h *JWTService) ParseAndVerifyJWT(token string) (*jwt.Token, error) {

	// Parse the JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.privateKey.Public(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT: %v", err)
	}
	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	return parsedToken, nil
}

// GetUserIdFromJWT retrieves the userId from a JWT token, using the ParseAndVerifyJWT function
func (h *JWTService) GetUserIdFromJWT(token string) (int64, error) {
	// Parse the JWT token
	parsedToken, err := h.ParseAndVerifyJWT(token)
	if err != nil {
		return 0, err
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	userId := claims["userId"].(float64)

	if userId == 0 {
		return 0, errors.New("invalid token")
	}
	return int64(userId), nil
}
