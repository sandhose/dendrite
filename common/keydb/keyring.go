// Copyright 2017 New Vector Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keydb

import (
	"context"
	"crypto/ed25519"

	"github.com/matrix-org/gomatrixserverlib"
)

// CreateKeyRing creates and configures a KeyRing object.
//
// It creates the necessary key fetchers and collects them into a KeyRing
// backed by the given KeyDatabase.
func CreateKeyRing(client gomatrixserverlib.Client,
	keyDB gomatrixserverlib.KeyDatabase) gomatrixserverlib.KeyRing {
	return gomatrixserverlib.KeyRing{
		KeyFetchers: []gomatrixserverlib.KeyFetcher{
			// TODO: Use perspective key fetchers for production.
			&gomatrixserverlib.DirectKeyFetcher{
				Client: client,
			},
		},
		KeyDatabase: keyDB,
	}
}

// Everything below might seem really strange, but amazingly, libp2p doesn't let
// us dial ourselves and seemingly gomatrixserverlib's DirectKeyFetcher always
// tries to do a HTTP dance even to get our own keys. The implementation below
// stores our own keys internally so that we can answer for them and avoid
// having to make any HTTP requests in that instance. This is pretty critical to
// federation working properly in the P2P branch.
func CreateDirectAvoidingLoopbackKeyRing(
	client gomatrixserverlib.Client,
	keyDB gomatrixserverlib.KeyDatabase,
	serverName gomatrixserverlib.ServerName,
	privateKey ed25519.PrivateKey,
	keyID gomatrixserverlib.KeyID,
) gomatrixserverlib.KeyRing {
	return gomatrixserverlib.KeyRing{
		KeyFetchers: []gomatrixserverlib.KeyFetcher{
			// TODO: Use perspective key fetchers for production.
			&DirectAvoidingLoopbackFetcher{
				DirectKeyFetcher: gomatrixserverlib.DirectKeyFetcher{
					Client: client,
				},
				serverName:  serverName,
				serverKey:   privateKey.Public().(ed25519.PublicKey),
				serverKeyID: string(keyID),
			},
		},
		KeyDatabase: keyDB,
	}
}

type DirectAvoidingLoopbackFetcher struct {
	gomatrixserverlib.DirectKeyFetcher
	serverName  gomatrixserverlib.ServerName
	serverKey   ed25519.PublicKey
	serverKeyID string
}

func (f *DirectAvoidingLoopbackFetcher) FetchKeys(
	ctx context.Context,
	requests map[gomatrixserverlib.PublicKeyLookupRequest]gomatrixserverlib.Timestamp,
) (map[gomatrixserverlib.PublicKeyLookupRequest]gomatrixserverlib.PublicKeyLookupResult, error) {
	results := make(map[gomatrixserverlib.PublicKeyLookupRequest]gomatrixserverlib.PublicKeyLookupResult)

	for request, timestamp := range requests {
		if request.ServerName == f.serverName {
			delete(requests, request)
			results[request] = gomatrixserverlib.PublicKeyLookupResult{
				VerifyKey: gomatrixserverlib.VerifyKey{
					Key: gomatrixserverlib.Base64String(f.serverKey),
				},
				ValidUntilTS: timestamp,
				ExpiredTS:    gomatrixserverlib.PublicKeyNotExpired,
			}
		}
	}

	directResults, err := f.DirectKeyFetcher.FetchKeys(ctx, requests)
	if err != nil {
		return nil, err
	}

	for request, result := range directResults {
		results[request] = result
	}

	return results, nil
}

func (f *DirectAvoidingLoopbackFetcher) FetcherName() string {
	return "DirectAvoidingLoopbackFetcher"
}
