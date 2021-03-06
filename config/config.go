package config

import (
	"time"
)

const DevFileName = ".appixDevSettings"
const IgnoreFileName = ".appixignore"

type Config struct {
	Version         string
	BuildDate       string
	ParsedBuildDate time.Time
	GitHash         string
	CatalogURL      string

	DirectoryPath string
	AuthFilePath  string

	DeveloperProfileUrl string

	FirebaseApiKey            string
	FirebaseAuthDomain        string
	FirebaseDatabaseUrl       string
	FirebaseStorageBucket     string
	FirebaseMessagingSenderId string
	FirebaseRefreshTokenUrl   string
	TravixLoggerUrl           string

	AuthServerPort   string
	MaxRetryAttempts int
}
