package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/fiatjaf/eventstore/lmdb"
	"github.com/fiatjaf/khatru"
	"github.com/joho/godotenv"
	"github.com/nbd-wtf/go-nostr"
)

// Struct of whitelisted users pubkeys
type Users struct {
	Pubkeys []string `json:"pubkeys"`
}

// Load whitelisted users from json file
func loadUsers(filename string) (*Users, error) {
	// Try to open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()
	// Try to read file
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	// Save and return whitelisted users
	var users Users
	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, fmt.Errorf("could not parse JSON: %w", err)
	}
	return &users, nil
}

// Struct for the relay configuration
type Config struct {
	RelayURL         string `json:"relay_url"`
	RelayPort        int    `json:"relay_port"`
	RelayBindAddress string `json:"relay_bind_address"`
	RelayName        string `json:"relay_name"`
	RelayPubkey      string `json:"relay_pubkey"`
	RelayDescription string `json:"relay_description"`
	RelayIcon        string `json:"relay_icon"`
	RelayContact     string `json:"relay_contact"`
	DBDir            string `json:"db_dir"`
	UsersFile        string `json:"users_file"`
	RelaySoftware    string `json:"relay_software"`
	RelayVersion     string `json:"relay_version"`
}

// Custom get env functions with defaults
// Get env as integer
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

// Get env as string
func getEnvString(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// Load configuration from environment variables
func loadConfig() Config {
	_ = godotenv.Load(".env")
	return Config{
		RelayURL:         os.Getenv("RELAY_URL"),
		RelayPort:        getEnvInt("RELAY_PORT", 2121),
		RelayBindAddress: getEnvString("RELAY_BIND_ADDRESS", "0.0.0.0"),
		RelayName:        os.Getenv("RELAY_NAME"),
		RelayPubkey:      os.Getenv("RELAY_PUBKEY"),
		RelayDescription: os.Getenv("RELAY_DESCRIPTION"),
		RelayIcon:        os.Getenv("RELAY_ICON"),
		RelayContact:     os.Getenv("RELAY_CONTACT"),
		DBDir:            getEnvString("DB_DIR", "db/"),
		UsersFile:        getEnvString("USERS_FILE", "users.json"),
		RelaySoftware:    "https://github.com/ciori/mattress",
		RelayVersion:     "v0.1.0",
	}
}

// Banner
var art = `
__________________________________________________
 __  __    _  _____ _____ ____  _____ ____ ____
|  \/  |  / \|_   _|_   _|  _ \| ____/ ___/ ___|
| |\/| | / _ \ | |   | | | |_) |  _| \___ \___ \
| |  | |/ ___ \| |   | | |  _ <| |___ ___) |__) |
|_|  |_/_/   \_\_|   |_| |_| \_\_____|____/____/
__________________________________________________
Nostr relay for Cashu tokens of whitelisted pubkeys.
`

func main() {
	// Initial welcome
	fmt.Println(art)
	fmt.Println("🚀 Mattress is booting up")

	// Load the environment variables
	fmt.Println("Loading the environment variables...")
	config := loadConfig()

	// Initialize the database
	fmt.Println("Setting up the database...")
	db := lmdb.LMDBBackend{
		Path: config.DBDir,
	}
	if err := db.Init(); err != nil {
		panic(err)
	}

	// Initialize the relay
	fmt.Println("Initializing the relay...")
	// Create the relay
	relay := khatru.NewRelay()
	// Set relay information
	relay.Info.Name = config.RelayName
	relay.Info.PubKey = config.RelayPubkey
	relay.Info.Icon = config.RelayIcon
	relay.Info.Contact = config.RelayContact
	relay.Info.Description = config.RelayDescription
	relay.Info.Software = config.RelaySoftware
	relay.Info.Version = config.RelayVersion

	// Load whitelisted users
	fmt.Println("Loading the whitelisted users...")
	users, err := loadUsers(config.UsersFile)
	if err != nil {
		fmt.Println("Error loading users:", err)
		return
	}

	// Configure the relay
	fmt.Println("Configuring the relay...")
	// Set relay functions
	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)
	relay.ReplaceEvent = append(relay.ReplaceEvent, db.ReplaceEvent)
	// Set required authentication on connect
	relay.OnConnect = append(relay.OnConnect, func(ctx context.Context) {
		khatru.RequestAuth(ctx)
	})
	// Only allow writing events from whitelisted users (NEEDS TO BE CHANGED)
	relay.RejectEvent = append(relay.RejectEvent, func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		// Event must have a pubkey
		if event.PubKey == "" {
			return true, "no pubkey"
		}
		// Allow all if there are no whitelisted users
		if len(users.Pubkeys) == 0 {
			return false, ""
		}
		// Check whether the user is whitelisted
		if slices.Contains(users.Pubkeys, event.PubKey) {
			return false, ""
		}
		return true, "pubkey not whitelisted"
	})
	// ...

	// Serve the relay
	endpoint := fmt.Sprintf("%s:%d", config.RelayBindAddress, config.RelayPort)
	fmt.Println("Running on", endpoint)
	http.ListenAndServe(endpoint, relay)
}
