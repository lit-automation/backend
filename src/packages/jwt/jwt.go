package jwt

import (
	"fmt"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

func CreateJWTToken(inXMinutes int64, userID uuid.UUID, scopes []string) (string, error) {
	token := jwtgo.New(jwtgo.SigningMethodRS512)

	now := time.Now()

	tokenClaims := jwtgo.MapClaims{
		"iss":    "System",                         // who creates the token and signs it
		"aud":    "User",                           // to whom the token is intended to be sent
		"exp":    inXMinutes,                       // time when the token will expire (x minutes from now)
		"jti":    uuid.Must(uuid.NewV4()),          // a unique identifier for the token
		"iat":    now.Unix(),                       // when the token was issued/created (now)
		"nbf":    now.Add(-2 * time.Minute).Unix(), // time before which the token is not yet valid (2 minutes ago)
		"sub":    userID,                           // the subject/principal is whom the token is about
		"scopes": scopes,                           // token scope - not a standard claim
	}

	token.Claims = tokenClaims

	privKey := PrivateKey()
	if privKey == nil {
		return "", fmt.Errorf("no private key set")
	}

	signedToken, err := token.SignedString(privKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token")
	}

	return signedToken, nil
}
