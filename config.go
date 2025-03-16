package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBEngine         string          `json:"db_engine"`
	LmdbMapSize      int64           `json:"lmdb_map_size"`
	RelayURL         string          `json:"relay_url"`
	RelayPort        int             `json:"relay_port"`
	RelayBindAddress string          `json:"relay_bind_address"`
	RelaySoftware    string          `json:"relay_software"`
	RelayVersion     string          `json:"relay_version"`
	RelayName        string          `json:"relay_name"`
	RelayNpub        string          `json:"relay_npub"`
	RelayDescription string          `json:"relay_description"`
	RelayIcon        string          `json:"relay_icon"`
	UserNpubs        map[string]bool `json:"user_npubs"`
}

func loadConfig() Config {
	_ = godotenv.Load(".env")

	return Config{
		DBEngine:         getEnvString("DB_ENGINE", "lmdb"),
		LmdbMapSize:      getEnvInt64("LMDB_MAPSIZE", 0),
		RelayURL:         getEnv("RELAY_URL"),
		RelayPort:        getEnvInt("RELAY_PORT", 2121),
		RelayBindAddress: getEnvString("RELAY_BIND_ADDRESS", "0.0.0.0"),
		RelaySoftware:    "https://github.com/ciori/mattress",
		RelayVersion:     "v0.0.1",
		RelayName:        getEnv("RELAY_NAME"),
		RelayNpub:        getEnv("RELAY_NPUB"),
		RelayDescription: getEnv("RELAY_DESCRIPTION"),
		RelayIcon:        getEnv("RELAY_ICON"),
		UserNpubs:        getUserNpubsFromFile(getEnv("USER_NPUBS_FILE")),
	}
}

func getUserNpubsFromFile(filePath string) map[string]bool {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}
	var userNpubs map[string]bool
	if err := json.Unmarshal(file, &userNpubs); err != nil {
		log.Fatalf("Failed to parse JSON: %s", err)
	}
	for npub := range userNpubs {
		userNpubs[strings.TrimSpace(npub)] = true
	}
	return userNpubs
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

func getEnvString(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return intValue
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
		return intValue
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		return boolValue
	}
	return defaultValue
}

var art = `
__________________________________________________
 __  __    _  _____ _____ ____  _____ ____ ____
|  \/  |  / \|_   _|_   _|  _ \| ____/ ___/ ___|
| |\/| | / _ \ | |   | | | |_) |  _| \___ \___ \
| |  | |/ ___ \| |   | | |  _ <| |___ ___) |__) |
|_|  |_/_/   \_\_|   |_| |_| \_\_____|____/____/
__________________________________________________
Nostr relay for Cashu tokens of whitelisted npubs.
`
