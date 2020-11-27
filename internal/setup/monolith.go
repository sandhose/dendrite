// Copyright 2020 The Matrix.org Foundation C.I.C.
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

package setup

import (
	"github.com/gorilla/mux"
	"github.com/ory/fosite"

	appserviceAPI "github.com/matrix-org/dendrite/appservice/api"
	"github.com/matrix-org/dendrite/authapi"
	authstorage "github.com/matrix-org/dendrite/authapi/storage"
	"github.com/matrix-org/dendrite/clientapi"
	"github.com/matrix-org/dendrite/clientapi/api"
	eduServerAPI "github.com/matrix-org/dendrite/eduserver/api"
	"github.com/matrix-org/dendrite/federationapi"
	federationSenderAPI "github.com/matrix-org/dendrite/federationsender/api"
	"github.com/matrix-org/dendrite/internal/config"
	"github.com/matrix-org/dendrite/internal/transactions"
	keyAPI "github.com/matrix-org/dendrite/keyserver/api"
	"github.com/matrix-org/dendrite/mediaapi"
	roomserverAPI "github.com/matrix-org/dendrite/roomserver/api"
	serverKeyAPI "github.com/matrix-org/dendrite/signingkeyserver/api"
	"github.com/matrix-org/dendrite/syncapi"
	userapi "github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/dendrite/userapi/storage/accounts"
	"github.com/matrix-org/dendrite/wellknownapi"
	"github.com/matrix-org/gomatrixserverlib"
)

// Monolith represents an instantiation of all dependencies required to build
// all components of Dendrite, for use in monolith mode.
type Monolith struct {
	Config         *config.Dendrite
	AccountDB      accounts.Database
	AuthDB         authstorage.Database
	OAuth2Provider fosite.OAuth2Provider
	KeyRing        *gomatrixserverlib.KeyRing
	Client         *gomatrixserverlib.Client
	FedClient      *gomatrixserverlib.FederationClient

	AppserviceAPI       appserviceAPI.AppServiceQueryAPI
	EDUInternalAPI      eduServerAPI.EDUServerInputAPI
	FederationSenderAPI federationSenderAPI.FederationSenderInternalAPI
	RoomserverAPI       roomserverAPI.RoomserverInternalAPI
	ServerKeyAPI        serverKeyAPI.SigningKeyServerAPI
	UserAPI             userapi.UserInternalAPI
	KeyAPI              keyAPI.KeyInternalAPI

	// Optional
	ExtPublicRoomsProvider api.ExtraPublicRoomsProvider
}

// AddAllPublicRoutes attaches all public paths to the given router
func (m *Monolith) AddAllPublicRoutes(authMux, csMux, ssMux, keyMux, mediaMux, wkMux *mux.Router) {
	authapi.AddPublicRoutes(authMux, &m.Config.AuthAPI, m.AuthDB, m.OAuth2Provider)
	clientapi.AddPublicRoutes(
		csMux, &m.Config.ClientAPI, m.AccountDB,
		m.FedClient, m.RoomserverAPI,
		m.EDUInternalAPI, m.AppserviceAPI, transactions.New(),
		m.FederationSenderAPI, m.UserAPI, m.KeyAPI, m.ExtPublicRoomsProvider,
	)
	federationapi.AddPublicRoutes(
		ssMux, keyMux, &m.Config.FederationAPI, m.UserAPI, m.FedClient,
		m.KeyRing, m.RoomserverAPI, m.FederationSenderAPI,
		m.EDUInternalAPI, m.KeyAPI,
	)
	mediaapi.AddPublicRoutes(mediaMux, &m.Config.MediaAPI, m.UserAPI, m.Client)
	syncapi.AddPublicRoutes(
		csMux, m.UserAPI, m.RoomserverAPI,
		m.KeyAPI, m.FedClient, &m.Config.SyncAPI,
	)
	wellknownapi.AddPublicRoutes(wkMux, &m.Config.WellKnownAPI)
}
