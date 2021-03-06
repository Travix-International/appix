package appcatalog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Travix-International/appix/appixLogger"
	"github.com/Travix-International/appix/config"
)

/**
 * - happy flow
 * - retries (1 - 3)
 * - catalog error (authentication)
 * - break prepare
 */

func TestPushToCatalog(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bodyBytes := []byte(`{ "messages": [ "success" ], "links": { "upload": "http://localhost:3001" } }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	appManifestFile := "../mocks/mockApp.manifest"

	conf := config.Config{
		AuthFilePath:            "../mocks/mockAuth.json",
		FirebaseApiKey:          "firebaseapikey",
		FirebaseRefreshTokenUrl: testServer.URL + "/firebase?key=",
		MaxRetryAttempts:        3,
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	uploadURI, err := PushToCatalog(testServer.URL, 10, appManifestFile, false, conf, logger)

	if err != nil {
		t.Fatalf("An error occured when testing PushToCatalog. Details: %s\n", err.Error())
	} else {
		if uploadURI == "http://localhost:3001" {
			t.Logf("The test for PushToCatalog went well.")
		} else {
			t.Fatalf("Unexpected returned value: %s\n", uploadURI)
		}
	}
}

func TestPushToCatalogWithRetries(t *testing.T) {
	counter := 0

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte

		if strings.Contains(r.URL.Path, "firebase") {
			w.WriteHeader(http.StatusOK)
			bodyBytes = []byte(`{ "access_token": "therightone", "expires_in": "3600", "expires_at": "1489675496997", "id_token": "idtoken", "project_id": "projectid", "refresh_token": "refreshtoken", "token_type": "Bearer", "user_id": "toto" }`)
		} else {
			if counter == 2 {
				w.WriteHeader(http.StatusOK)
				bodyBytes = []byte(`{ "messages": [ "success" ], "links": { "acc:query": "http://localhost:3001" } }`)
			} else {
				w.WriteHeader(http.StatusNotAcceptable)
				bodyBytes = []byte(`{ "messages": [ "failed" ], "links": { "acc:query": "http://localhost:3001" } }`)
			}
			counter++
		}

		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	appManifestFile := "../mocks/mockApp.manifest"

	conf := config.Config{
		AuthFilePath:            "../mocks/mockAuth.json",
		FirebaseApiKey:          "firebaseapikey",
		FirebaseRefreshTokenUrl: testServer.URL + "/firebase?key=",
		MaxRetryAttempts:        3,
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	_, err := PushToCatalog(testServer.URL, 10, appManifestFile, false, conf, logger)

	if err != nil {
		t.Fatalf("An error occured when testing PushToCatalog. Details: %s\n", err.Error())
	} else {
		t.Logf("The test for PushToCatalog went well.")
	}
}

func TestPushToCatalogFailAuthentication(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		bodyBytes := []byte(`{ "messages": [ "failed" ], "links": { "upload": "http://localhost:3001" } }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	appManifestFile := "./../mocks/mockApp.manifest"

	conf := config.Config{
		AuthFilePath:            "../mocks/mockAuth.json",
		FirebaseApiKey:          "firebaseapikey",
		FirebaseRefreshTokenUrl: testServer.URL + "/firebase?key=",
		MaxRetryAttempts:        3,
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	uploadURI, err := PushToCatalog(testServer.URL, 10, appManifestFile, false, conf, logger)

	if err != nil {
		t.Logf("An error occured when testing PushToCatalog as expected. Details: %s\n", err.Error())
	} else {
		t.Fatalf("Expecting an error. Unexpected returned value: %s\n", uploadURI)
	}
}

func TestPushToCatalogFailServer(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		bodyBytes := []byte(`{ "messages": [ "failed" ], "links": { "upload": "http://localhost:3001" } }`)
		w.Write(bodyBytes)
	}))

	defer testServer.Close()

	appManifestFile := "../mocks/mockApp.manifest"

	conf := config.Config{
		AuthFilePath:            "../mocks/mockAuth.json",
		FirebaseApiKey:          "firebaseapikey",
		FirebaseRefreshTokenUrl: testServer.URL + "/firebase?key=",
		MaxRetryAttempts:        3,
	}

	logger := appixLogger.NewAppixLogger(testServer.URL)
	logger.Start()
	defer logger.Stop()

	uploadURI, err := PushToCatalog(testServer.URL, 10, appManifestFile, false, conf, logger)

	if err != nil {
		t.Logf("An error occured when testing PushToCatalog as expected. Details: %s\n", err.Error())
	} else {
		t.Fatalf("Expecting an error. Unexpected returned value: %s\n", uploadURI)
	}
}
