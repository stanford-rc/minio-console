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
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/minio/console/pkg/auth/idp/oauth2"
)

// mock function of Get()
func (ac consoleCredentialsMock) Expire() {
	// Do nothing
	// Implementing this method for the consoleCredentials interface
}

const (
	stanfordEndSession = "https://test.elm.stanford.edu/logout"
	otherEndSession    = "https://idp.example/end-session"
)

// withProviders installs an OpenID provider map for the duration of one test.
// GlobalMinIOConfig is process wide and is normally written once at startup by
// the MinIO server, so it has to be put back.
func withProviders(t *testing.T, providers oauth2.OpenIDPCfg) {
	t.Helper()

	saved := GlobalMinIOConfig
	t.Cleanup(func() { GlobalMinIOConfig = saved })

	GlobalMinIOConfig = MinIOConfig{OpenIDProviders: providers}
}

// loginState builds the state the console's own OIDC callback stores in
// localStorage under auth-state, which is a base64 JSON LoginURLParams.
func loginState(t *testing.T, idpName string) string {
	t.Helper()

	enc, err := json.Marshal(oauth2.LoginURLParams{State: "irrelevant", IDPName: idpName})
	if err != nil {
		t.Fatalf("marshalling the login state: %s", err)
	}

	return base64.StdEncoding.EncodeToString(enc)
}

// The Elm case, and the reason this function accepts an empty state at all.
//
// elm-console performs the OIDC login itself and then logs the user into the
// console directly with an access key and secret key, so the browser never
// reaches /oauth_callback, LoginCallback never runs, and localStorage holds no
// auth-state. Sign-out therefore posts an empty state.
//
// Treating that as "nothing to do" made sign-out a silent local-only logout:
// the console cleared its own cookies, returned no redirect and logged no
// error, and the user kept a live session at the IdP.
func TestPostLogoutRedirectURLResolvesTheSoleProviderWithoutState(t *testing.T) {
	withProviders(t, oauth2.OpenIDPCfg{
		// "_" is the key a single unnamed provider is registered under, which
		// is what MINIO_IDENTITY_OPENID_CONFIG_URL with no _<NAME> suffix
		// produces.
		"_": {EndSessionEndpoint: stanfordEndSession},
	})

	got, err := postLogoutRedirectURL("")
	if err != nil {
		t.Fatalf("postLogoutRedirectURL(\"\") = error %s, want the sole provider's endpoint", err)
	}
	if got != stanfordEndSession {
		t.Errorf("postLogoutRedirectURL(\"\") = %q, want %q", got, stanfordEndSession)
	}
}

// A console with no OIDC at all logs in against the root credentials, so there
// is no IdP session to end. That is an ordinary logout and must not report an
// error, which would put "IDP logout failed" in the log on every sign-out.
func TestPostLogoutRedirectURLWithoutStateOrProviders(t *testing.T) {
	withProviders(t, oauth2.OpenIDPCfg{})

	got, err := postLogoutRedirectURL("")
	if err != nil {
		t.Fatalf("unexpected error with no providers configured: %s", err)
	}
	if got != "" {
		t.Errorf("postLogoutRedirectURL(\"\") = %q, want \"\"", got)
	}
}

// Several providers and no state is a real ambiguity rather than something to
// guess at, since picking the wrong one would send the browser to an IdP the
// user has no session with.
func TestPostLogoutRedirectURLWithoutStateRefusesToGuessBetweenProviders(t *testing.T) {
	withProviders(t, oauth2.OpenIDPCfg{
		"_":     {EndSessionEndpoint: stanfordEndSession},
		"other": {EndSessionEndpoint: otherEndSession},
	})

	got, err := postLogoutRedirectURL("")
	if err == nil {
		t.Fatalf("postLogoutRedirectURL(\"\") = %q, want an error naming the ambiguity", got)
	}
	if got != "" {
		t.Errorf("a refused lookup must not also return a target, got %q", got)
	}
}

// The upstream path, which the empty-state handling must leave intact: a state
// is present, so it selects the provider even when several are configured.
func TestPostLogoutRedirectURLSelectsTheNamedProvider(t *testing.T) {
	withProviders(t, oauth2.OpenIDPCfg{
		"_":     {EndSessionEndpoint: stanfordEndSession},
		"other": {EndSessionEndpoint: otherEndSession},
	})

	for name, want := range map[string]string{
		"_":     stanfordEndSession,
		"other": otherEndSession,
	} {
		got, err := postLogoutRedirectURL(loginState(t, name))
		if err != nil {
			t.Errorf("provider %q: unexpected error: %s", name, err)
			continue
		}
		if got != want {
			t.Errorf("provider %q: got %q, want %q", name, got, want)
		}
	}
}

// A provider advertising no end_session_endpoint is not an error. It means the
// IdP supports no RP-initiated logout, so a local logout is all there is.
func TestPostLogoutRedirectURLEmptyWhenProviderAdvertisesNoEndpoint(t *testing.T) {
	withProviders(t, oauth2.OpenIDPCfg{"_": {EndSessionEndpoint: ""}})

	for _, state := range []string{"", loginState(t, "_")} {
		got, err := postLogoutRedirectURL(state)
		if err != nil {
			t.Errorf("state %q: unexpected error: %s", state, err)
		}
		if got != "" {
			t.Errorf("state %q: got %q, want \"\"", state, got)
		}
	}
}

// A state that arrives but names nothing we know is a different situation from
// no state at all, and stays an error: something went wrong rather than the
// deployment simply not using the callback route.
func TestPostLogoutRedirectURLRejectsUnusableState(t *testing.T) {
	withProviders(t, oauth2.OpenIDPCfg{"_": {EndSessionEndpoint: stanfordEndSession}})

	for name, state := range map[string]string{
		"unknown provider": loginState(t, "no-such-provider"),
		"not base64":       "!!!not base64!!!",
		"not json":         base64.StdEncoding.EncodeToString([]byte("not json")),
	} {
		got, err := postLogoutRedirectURL(state)
		if err == nil {
			t.Errorf("%s: got %q, want an error", name, got)
		}
		if got != "" {
			t.Errorf("%s: a failed lookup must not return a target, got %q", name, got)
		}
	}
}
