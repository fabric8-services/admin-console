package jwt

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
)

// NewJWTContext creates a context with a JWT having the given `subject`
// using the private key that is located in the 'privateKeyPath'
func NewJWTContext(subject, kid, privateKey string) (context.Context, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = subject
	tk := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tk.Header["kid"] = kid
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, err
	}
	signed, err := tk.SignedString(key)
	if err != nil {
		return nil, err
	}
	tk.Raw = signed
	return goajwt.WithJWT(context.Background(), tk), nil
}
