package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/fiatjaf/eventstore/badger"
	"github.com/fiatjaf/eventstore/lmdb"
	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/policies"
	"github.com/nbd-wtf/go-nostr"
)

var (
	relay   = khatru.NewRelay()
	db = newDBBackend("db/relay")
)

type DBBackend interface {
	Init() error
	Close()
	CountEvents(ctx context.Context, filter nostr.Filter) (int64, error)
	DeleteEvent(ctx context.Context, evt *nostr.Event) error
	QueryEvents(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error)
	SaveEvent(ctx context.Context, evt *nostr.Event) error
	ReplaceEvent(ctx context.Context, evt *nostr.Event) error
	Serial() []byte
}

func newDBBackend(path string) DBBackend {
	switch config.DBEngine {
	case "lmdb":
		return newLMDBBackend(path)
	case "badger":
		return &badger.BadgerBackend{
			Path: path,
		}
	default:
		return newLMDBBackend(path)
	}
}

func newLMDBBackend(path string) *lmdb.LMDBBackend {
	return &lmdb.LMDBBackend{
		Path:    path,
		MapSize: config.LmdbMapSize,
	}
}

func initRelays() {
	if err := db.Init(); err != nil {
		panic(err)
	}

	initRelayLimits()

	relay.Info.Name = config.RelayName
	relay.Info.PubKey = nPubToPubkey(config.RelayNpub)
	relay.Info.Description = config.RelayDescription
	relay.Info.Icon = config.RelayIcon
	relay.Info.Version = config.RelayVersion
	relay.Info.Software = config.RelaySoftware
	relay.ServiceURL = "https://" + config.RelayURL + "/relay"

	if !relayLimits.AllowEmptyFilters {
		relay.RejectFilter = append(relay.RejectFilter, policies.NoEmptyFilters)
	}

	if !relayLimits.AllowComplexFilters {
		relay.RejectFilter = append(relay.RejectFilter, policies.NoComplexFilters)
	}

	relay.RejectEvent = append(relay.RejectEvent,
		policies.RejectEventsWithBase64Media,
		policies.EventIPRateLimiter(
			relayLimits.EventIPLimiterTokensPerInterval,
			time.Minute*time.Duration(relayLimits.EventIPLimiterInterval),
			relayLimits.EventIPLimiterMaxTokens,
		),
	)

	relay.RejectConnection = append(relay.RejectConnection,
		policies.ConnectionRateLimiter(
			relayLimits.ConnectionRateLimiterTokensPerInterval,
			time.Minute*time.Duration(relayLimits.ConnectionRateLimiterInterval),
			relayLimits.ConnectionRateLimiterMaxTokens,
		),
	)

	relay.OnConnect = append(relay.OnConnect, func(ctx context.Context) {
		khatru.RequestAuth(ctx)
	})

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.ReplaceEvent = append(relay.ReplaceEvent, db.ReplaceEvent)

	relay.RejectFilter = append(relay.RejectFilter, func(ctx context.Context, filter nostr.Filter) (bool, string) {
		authenticatedUser := khatru.GetAuthed(ctx)
		if authenticatedUser == nPubToPubkey(config.OwnerNpub) {
			return false, ""
		}

		return true, "auth-required: this query requires you to be authenticated"
	})

	relay.RejectEvent = append(relay.RejectEvent, func(ctx context.Context, event *nostr.Event) (bool, string) {
		authenticatedUser := khatru.GetAuthed(ctx)

		if authenticatedUser == nPubToPubkey(config.OwnerNpub) {
			return false, ""
		}

		return true, "auth-required: publishing this event requires authentication"
	})

	mux := relay.Router()

	mux.HandleFunc("GET /relay", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		data := struct {
			RelayName        string
			RelayPubkey      string
			RelayDescription string
			RelayURL         string
		}{
			RelayName:        config.RelayName,
			RelayPubkey:      nPubToPubkey(config.RelayNpub),
			RelayDescription: config.RelayDescription,
			RelayURL:         "wss://" + config.RelayURL + "/relay",
		}
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
