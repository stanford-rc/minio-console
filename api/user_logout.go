// This file is part of MinIO Console Server
// Copyright (c) 2021 MinIO, Inc.
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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/minio/console/api/operations"
	authApi "github.com/minio/console/api/operations/auth"
	"github.com/minio/console/models"
	"github.com/minio/console/pkg/auth/idp/oauth2"
)

func registerLogoutHandlers(api *operations.ConsoleAPI) {
	// logout from console
	api.AuthLogoutHandler = authApi.LogoutHandlerFunc(func(params authApi.LogoutParams, session *models.Principal) middleware.Responder {
		redirect, err := getLogoutResponse(session, params)
		if err != nil {
			api.Logger("IDP logout failed: %v", err.APIError.DetailedMessage)
		}
		// Custom response writer to expire the session cookies
		return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
			if err != nil {
				w.Header().Set("IDP-Logout", fmt.Sprintf("%v", err.APIError.DetailedMessage))
			}
			if redirect != "" {
				w.Header().Set(PostLogoutRedirectHeader, redirect)
			}
			expiredCookie := ExpireSessionCookie()
			// this will tell the browser to clear the cookie and invalidate user session
			// additionally we are deleting the cookie from the client side
			http.SetCookie(w, &expiredCookie)
			http.SetCookie(w, &http.Cookie{
				Path:     "/",
				Name:     "idp-refresh-token",
				Value:    "",
				MaxAge:   -1,
				Expires:  time.Now().Add(-100 * time.Hour),
				HttpOnly: true,
				Secure:   len(GlobalPublicCerts) > 0,
				SameSite: http.SameSiteLaxMode,
			})
			authApi.NewLogoutOK().WriteResponse(w, p)
		})
	})
}

// logout() call Expire() on the provided ConsoleCredentials
func logout(credentials ConsoleCredentialsI) {
	credentials.Expire()
}

// getLogoutResponse performs logout() and returns the URL the browser should
// be redirected to afterwards, which is empty when the provider advertises no
// end-session endpoint, plus nil or errors.
func getLogoutResponse(session *models.Principal, params authApi.LogoutParams) (string, *CodedAPIError) {
	ctx, cancel := context.WithCancel(params.HTTPRequest.Context())
	defer cancel()

	redirect, err := postLogoutRedirectURL(params.Body.State)
	if err != nil {
		return "", ErrorWithContext(ctx, err)
	}

	creds := getConsoleCredentialsFromSession(session)
	credentials := ConsoleCredentials{ConsoleCredentials: creds}
	logout(credentials)

	return redirect, nil
}

// PostLogoutRedirectHeader carries the URL the browser should be sent to once
// the console has finished logging out, so that the session is also ended at
// the identity provider. It is a header rather than a response field because
// the /logout response carries no schema, and adding one means regenerating
// both the Go and the TypeScript clients for a single string.
const PostLogoutRedirectHeader = "X-Console-Post-Logout-Url"

// postLogoutRedirectURL returns the provider's end-session URL, or "" when no
// provider is configured or the resolved one advertises no such endpoint.
//
// The login state is optional. It carries the IDPName chosen at login and is
// written only by LoginCallback, which runs when the browser completes the
// console's OWN OIDC round trip and lands on /oauth_callback. A deployment
// that authenticates the user elsewhere and then logs them into the console
// directly never visits that route, so localStorage holds no auth-state and
// the provider has to be resolved from configuration instead.
//
// Treating an absent state as "nothing to do" is what made sign-out a silent
// local-only logout for those deployments: no redirect, no error, and the
// user left with a live IdP session.
//
// This replaces a back channel POST of client_id, client_secret and
// refresh_token to that same URL, which required a 204 to consider the logout
// successful. That was wrong on its own terms: OpenID Connect RP-Initiated
// Logout defines end_session_endpoint as a URL the USER AGENT is sent to, not
// a back channel. The POST therefore reached an endpoint that does not answer
// it, failed every logout it was enabled for, and sent our client secret
// somewhere that has no use for it.
//
// Returning the URL for the browser to follow also means the refresh token is
// no longer needed here, so a missing idp-refresh-token cookie stops being a
// reason to fail the logout.
func postLogoutRedirectURL(state string) (string, error) {
	if state == "" {
		// More than one provider is a genuine ambiguity: without the state
		// there is nothing to say which session the user holds.
		if n := len(GlobalMinIOConfig.OpenIDProviders); n > 1 {
			return "", fmt.Errorf("no login state supplied and %d identity "+
				"providers are configured, cannot determine which session to end", n)
		}

		// Exactly one provider needs no disambiguating. None at all means
		// there is no IdP session to end, so "" is the right answer.
		for _, providerCfg := range GlobalMinIOConfig.OpenIDProviders {
			return providerCfg.EndSessionEndpoint, nil
		}

		return "", nil
	}

	decodedRState, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return "", err
	}

	var requestItems oauth2.LoginURLParams
	if err := json.Unmarshal(decodedRState, &requestItems); err != nil {
		return "", err
	}

	providerCfg, ok := GlobalMinIOConfig.OpenIDProviders[requestItems.IDPName]
	if !ok {
		return "", fmt.Errorf("unrecognized identity provider %q", requestItems.IDPName)
	}

	return providerCfg.EndSessionEndpoint, nil
}
