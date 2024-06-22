package main

import (
	"context"
	"encoding/json"
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
	homeURL      = "https://nh.xvm.mit.edu/"

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
	session, _ := store.Get(r, "auth")
	session.Values["access_token"] = token.AccessToken

    rawIDToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "No id_token in token response", http.StatusInternalServerError)
        return
    }

    //idToken, err := provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(ctx, rawIDToken)
	if _, err = provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(ctx, rawIDToken); err != nil {
		http.Error(w, "Failed to verify id_token", http.StatusInternalServerError)
		return
	}

	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		//fmt.Fprintf(w, "%s %+v\n", idToken.Subject, idToken)
		// save the idToken.Subject in session
		//fmt.Fprintf(w, "Successfully authenticated as %s.\n", idToken.Subject)
		http.Redirect(w, r, homeURL, http.StatusFound)
	}

    //var claims struct {
    //	Email string `json:"email"`
    //}
    //if err := idToken.Claims(&claims); err != nil {
    //	http.Error(w, "Failed to get claims", http.StatusInternalServerError)
    //	return
    //}
}


func getUserInfo(r *http.Request) (map[string]interface{}, error) {
	session, _ := store.Get(r, "auth")
	token, exists := session.Values["access_token"]
	if !exists {
		return nil, nil
	}

    client := &http.Client{}

    req, err := http.NewRequest("GET", "https://petrock.mit.edu/oidc/userinfo", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer " + token.(string))

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


func authMiddleware(w http.ResponseWriter, r *http.Request) (string, bool) {
	if userInfo, err := getUserInfo(r); err != nil {
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return "", false
	} else if userInfo == nil {
		http.Error(w, "Not logged in", http.StatusForbidden)
		return "", false
	} else {
		return userInfo["email"].(string), true
	}
}

func getProfile(w http.ResponseWriter, r *http.Request) {	
	ret := make(map[string]string)

	if id, cont := authMiddleware(w, r); cont {
		ret["email"] = id
		json.NewEncoder(w).Encode(ret)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {	
	session, _ := store.Get(r, "auth")
	delete(session.Values, "access_token")
	if err := session.Save(r, w); err != nil {
        http.Error(w, "Failed to logout", http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, homeURL, http.StatusFound)
	}
}
