package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)


var (
    clientID     = "REDACTED"
    clientSecret = "REDACTED"
    redirectURL  = "https://nh.xvm.mit.edu/oidc-response"
    providerURL  = "https://petrock.mit.edu"

    ctx context.Context
    provider *oidc.Provider
    oauth2Config oauth2.Config
)


func setupOIDC() {
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")

	ctx = context.Background()

    var err error
    provider, err = oidc.NewProvider(ctx, providerURL)
	if err != nil {
        log.Fatalf("Failed to get provider: %v", err)
    }

	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
}


func initOIDCAuth(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, oauth2Config.AuthCodeURL("state"), http.StatusFound)
}


func handleOIDCResponse(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    code := r.FormValue("code")
    if code == "" {
        http.Error(w, "No code in request", http.StatusBadRequest)
        return
    }

    token, err := oauth2Config.Exchange(ctx, code)
    if err != nil {
        http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
        return
    }

    rawIDToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "No id_token in token response", http.StatusInternalServerError)
        return
    }

    idToken, err := provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(ctx, rawIDToken)
    //fmt.Fprintf(w, "%s %+v\n", idToken.Subject, idToken)
    // save the idToken.Subject in session
    fmt.Fprintf(w, "Successfully authenticated as %s. Click the browser BACK button.\n", idToken.Subject)
    if err != nil {
        http.Error(w, "Failed to verify id_token", http.StatusInternalServerError)
        return
    }

    //var claims struct {
    //	Email string `json:"email"`
    //}
    //if err := idToken.Claims(&claims); err != nil {
    //	http.Error(w, "Failed to get claims", http.StatusInternalServerError)
    //	return
    //}

    //fmt.Fprintf(w, "Hello, %s!", claims.Email)
    //userInfo, err := getUserInfo(token, "https://petrock.mit.edu/oidc/userinfo")
    //if err != nil {
    //    http.Error(w, "Failed to get user info", http.StatusInternalServerError)
    //    return
    //}

    //userInfoJSON, err := json.MarshalIndent(userInfo, "", "  ")
    //if err != nil {
    //    http.Error(w, "Failed to marshal user info", http.StatusInternalServerError)
    //    return
    //}

    //fmt.Fprintf(w, "Hello, %s!\n\nUser Info:\n%s", claims.Email, userInfoJSON)
}

func getUserInfo(token *oauth2.Token, userInfoEndpoint string) (map[string]interface{}, error) {
    client := &http.Client{}

    req, err := http.NewRequest("GET", userInfoEndpoint, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer " + token.AccessToken)

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var userInfo map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        return nil, err
    }

    return userInfo, nil
}

func getProfile(w http.ResponseWriter, r *http.Request) {
}
