// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwt

import (
	"fmt"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

type Verifier interface {
	Verify(token string) (*Claims, error)
}

type verifier struct {
	key    interface{}
	method jwtgo.SigningMethod
}

// NewVerifier returns a new verifier using given signing method.
func NewVerifier(method jwtgo.SigningMethod, keyFile string) (Verifier, error) {
	key, err := readKeyFile(method, keyFile, false)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %v", err)
	}
	return &verifier{
		key:    key,
		method: method,
	}, nil
}

func (v *verifier) Verify(tokenString string) (*Claims, error) {
	// NOTE: The issuedAt and notBefore claims are set to "used if exists" by default.
	// ref: https://github.com/golang-jwt/jwt/issues/411#issuecomment-2423818974
	parser := jwtgo.NewParser(
		jwtgo.WithIssuer(Issuer),
		jwtgo.WithIssuedAt(),
		jwtgo.WithExpirationRequired(),
	)

	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwtgo.Token) (interface{}, error) {
		if v.method != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return v.key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("unexpected claims type: %T", token.Claims)
	}
	return claims, nil
}
