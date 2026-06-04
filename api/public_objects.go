// This file is part of MinIO Console Server
// Copyright (c) 2024 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"encoding/base64"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/minio/console/api/operations"
	"github.com/minio/console/api/operations/public"
	xnet "github.com/minio/pkg/v3/net"
)

func registerPublicObjectsHandlers(api *operations.ConsoleAPI) {
	api.PublicDownloadSharedObjectHandler = public.DownloadSharedObjectHandlerFunc(func(params public.DownloadSharedObjectParams) middleware.Responder {
		resp, err := getDownloadPublicObjectResponse(params)
		if err != nil {
			return public.NewDownloadSharedObjectDefault(err.Code).WithPayload(err.APIError)
		}
		return resp
	})
}

// getDownloadPublicObjectResponse previously had the console fetch a presigned
// object from MinIO server-side and stream it back. Stanford RC: this public,
// unauthenticated endpoint is disabled because it (and the presigned URL it
// served) bypassed elm-proxy. See getShareObjectResponse for the companion
// change. decodeMinIOStringURL is retained (referenced by tests).
func getDownloadPublicObjectResponse(params public.DownloadSharedObjectParams) (middleware.Responder, *CodedAPIError) {
	ctx := params.HTTPRequest.Context()
	return nil, ErrorWithContext(ctx, ErrForbidden)
}

// decodeMinIOStringURL decodes url and validates is a MinIO url endpoint
func decodeMinIOStringURL(inputURL string) (*string, error) {
	decodedURL, err := base64.RawURLEncoding.DecodeString(inputURL)
	if err != nil {
		return nil, err
	}

	// Validate input URL
	parsedURL, err := xnet.ParseHTTPURL(string(decodedURL))
	if err != nil {
		return nil, err
	}
	// Ensure incoming url points to MinIO Server
	minIOHost := getMinIOEndpoint()
	if parsedURL.Host != minIOHost {
		return nil, ErrForbidden
	}
	return swag.String(string(decodedURL)), nil
}
