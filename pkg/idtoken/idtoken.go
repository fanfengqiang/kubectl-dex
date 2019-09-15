package idtoken

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const exampleAppState = "Hello everyone"

type app struct {
	clientID     string
	clientSecret string
	redirectURI  string
	scopes       []string

	client *http.Client
	ctx    context.Context

	verifier *oidc.IDTokenVerifier
	provider *oidc.Provider
}

func (a *app) oauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.clientID,
		ClientSecret: a.clientSecret,
		Endpoint:     a.provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  a.redirectURI,
	}
}

// Config is the function that will return a Config.
func Config(clientID, clientSecret, redirectURI, issuerURL string, scopes []string) app {
	var a app
	a.clientID = clientID
	a.clientSecret = clientSecret
	a.redirectURI = redirectURI
	a.scopes = scopes

	a.client = http.DefaultClient
	a.ctx = oidc.ClientContext(context.Background(), a.client)
	provider, err := oidc.NewProvider(a.ctx, issuerURL)
	if err != nil {
		log.Fatalln("Failed to query provider:", issuerURL, err)
	}
	a.provider = provider
	a.verifier = provider.Verifier(&oidc.Config{ClientID: a.clientID})

	return a
}

// GetIDToken is the method that gets the ID Token.
func (a *app) GetIDToken(user, passwd string) (*oauth2.Token, error) {

	authCodeURL := a.oauth2Config(a.scopes).AuthCodeURL(exampleAppState)
	resp, err := http.Get(authCodeURL)
	if err != nil {
		// handle error
		log.Fatalln("Failed to connetc to issuer", err)
	}
	defer resp.Body.Close()
	postURL := resp.Request.URL.String()

	data := url.Values{"login": {user}, "password": {passwd}}
	body := strings.NewReader(data.Encode())
	clt := http.Client{
		// forbid to redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp2, err := clt.Post(postURL, "application/x-www-form-urlencoded", body)
	if err != nil {
		log.Fatalln("Fail to post data", err)
	}
	defer resp2.Body.Close()

	l, err := resp2.Location()
	if err != nil {
		log.Fatalln("Fail to get redirectURI1", err)
	}
	redirectURI1 := l.String()

	resp3, err := clt.Get(redirectURI1)
	if err != nil {
		log.Fatalln("Fail to get redirectURI2", err)
	}
	defer resp3.Body.Close()

	redirectURI2, err := resp3.Location()
	if err != nil {
		log.Fatalln("Fail to get redirectURI2", err)
	}
	parses, err := url.ParseQuery(redirectURI2.RawQuery)
	if err != nil {
		log.Fatalln("Fail to Parse redirectURI2", err)
	}
	code := parses["code"][0]

	oauth2Config := a.oauth2Config(nil)

	token, err := oauth2Config.Exchange(a.ctx, code)
	if err != nil {
		log.Fatalln("failed to Exchange code", err)
	}
	return token, nil
}

// ParseToken is the function that parses the token
func (a *app) ParseToken(token *oauth2.Token) string {
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Fatalln("no id_token in token response")
	}

	rawIDToken, err := a.verifier.Verify(context.Background(), idToken)
	if err != nil {
		log.Fatalln("failed to verify ID token: ", err)
	}

	accessToken, ok := token.Extra("access_token").(string)
	if !ok {
		log.Fatalln("no access_token in token response")
	}

	var claims json.RawMessage
	if err := rawIDToken.Claims(&claims); err != nil {
		log.Fatalln("error decoding ID token claims:", err)
	}

	buff := new(bytes.Buffer)
	if err := json.Indent(buff, []byte(claims), "", "  "); err != nil {
		log.Fatalln("error indenting ID token claims:", err)
	}

	fmt.Println()
	fmt.Println("idtoken:", idToken)
	fmt.Println()
	fmt.Println("Access Token:", accessToken)
	fmt.Println()
	fmt.Println("Refresh Token:", token.RefreshToken)
	fmt.Println()
	fmt.Println("claims:", buff.String())

	return idToken
}
