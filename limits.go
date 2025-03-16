package main

import (
	"encoding/json"
	"log"
)

var (
	relayLimits RelayLimits
)

type RelayLimits struct {
	EventIPLimiterTokensPerInterval        int
	EventIPLimiterInterval                 int
	EventIPLimiterMaxTokens                int
	AllowEmptyFilters                      bool
	AllowComplexFilters                    bool
	ConnectionRateLimiterTokensPerInterval int
	ConnectionRateLimiterInterval          int
	ConnectionRateLimiterMaxTokens         int
}

func initRelayLimits() {
	relayLimits = RelayLimits{
		EventIPLimiterTokensPerInterval:        getEnvInt("RELAY_EVENT_IP_LIMITER_TOKENS_PER_INTERVAL", 50),
		EventIPLimiterInterval:                 getEnvInt("RELAY_EVENT_IP_LIMITER_INTERVAL", 1),
		EventIPLimiterMaxTokens:                getEnvInt("RELAY_EVENT_IP_LIMITER_MAX_TOKENS", 100),
		AllowEmptyFilters:                      getEnvBool("RELAY_ALLOW_EMPTY_FILTERS", true),
		AllowComplexFilters:                    getEnvBool("RELAY_ALLOW_COMPLEX_FILTERS", true),
		ConnectionRateLimiterTokensPerInterval: getEnvInt("RELAY_CONNECTION_RATE_LIMITER_TOKENS_PER_INTERVAL", 3),
		ConnectionRateLimiterInterval:          getEnvInt("RELAY_CONNECTION_RATE_LIMITER_INTERVAL", 5),
		ConnectionRateLimiterMaxTokens:         getEnvInt("RELAY_CONNECTION_RATE_LIMITER_MAX_TOKENS", 9),
	}

	prettyPrintLimits("Relay limits", relayLimits)
}

func prettyPrintLimits(label string, value interface{}) {
	b, _ := json.MarshalIndent(value, "", "  ")
	log.Printf("🚧 %s:\n%s\n", label, string(b))
}
