# Mattress

_"Hide your cash under the mattress!"_

Nostr relay for Cashu tokens of whitelisted pubkeys.

Based on and inspired by the [haven](https://github.com/bitvora/haven) and [sw2](https://github.com/bitvora/sw2) relays.

## TODOs

- [ ] Only allow cashu wallet related events and kinds
- [ ] Allow correct "public" read and write permissions based on cashu flow
- [ ] Add rate limiters and other useful policies
- [ ] check edge cases where there are no filters/pubkeys/etc...
- [ ] Change user pubkyes list with a map anc change the auth condition to leverage it

## Events/Filters Logic

- 17375 (wallet event): read/write only whitelisted user for own events
- 7375 (proofs): "
- 7376 (history): "
- 7374 (quotes): "
- 10019 (public info):
  - pubkey write own events
  - all read events
- 9321 (nutzap):
  - sender pubkey: write/read only own events if dest is whitelisted
  - dest whitelisted pubkey: read
