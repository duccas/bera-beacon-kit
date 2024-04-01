// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package http

import (
	"net/http"

	"github.com/berachain/beacon-kit/mod/node-builder/utils/jwt"
	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/node"
)

// NewHeaderWithJWT creates a new HTTP header with a JWT token.
func NewHeaderWithJWT(secret *jwt.Secret) http.Header {
	header := http.Header{}
	if err := AddJWTHeader(header, secret); err != nil {
		panic(err)
	}
	return header
}

// AddJWTHeader adds a JWT header to the provided HTTP header.
func AddJWTHeader(header http.Header, secret *jwt.Secret) error {
	// Authenticate the execution node JSON-RPC endpoint.
	if secret == nil {
		return errors.New("jwt.Secret is nil")
	} else if header == nil {
		return errors.New("http.Header is nil")
	}
	return node.NewJWTAuth(*secret)(header)
}