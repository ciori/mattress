package main

import (
	"github.com/nbd-wtf/go-nostr/nip19"
)

func npubToPubkey(nPub string) string {
	_, v, err := nip19.Decode(nPub)
	if err != nil {
		panic(err)
	}
	return v.(string)
}

func pubkeyToNpub(pubkey string) string {
	v, err := nip19.EncodePublicKey(pubkey)
	if err != nil {
		panic(err)
	}
	return v
}
