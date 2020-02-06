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

package storage

import (
	"context"
	"net/url"

	eventDatabase "github.com/matrix-org/dendrite/roomserver/input"
	"github.com/matrix-org/dendrite/roomserver/storage/postgres"
	"github.com/matrix-org/dendrite/roomserver/types"
	"github.com/matrix-org/gomatrixserverlib"
)

type Database interface {
	eventDatabase.RoomEventDatabase
	EventStateKeys(ctx context.Context, eventStateKeyNIDs []types.EventStateKeyNID) (map[types.EventStateKeyNID]string, error)
	EventNIDs(ctx context.Context, eventIDs []string) (map[string]types.EventNID, error)
	RoomNID(ctx context.Context, roomID string) (types.RoomNID, error)
	LatestEventIDs(ctx context.Context, roomNID types.RoomNID) ([]gomatrixserverlib.EventReference, types.StateSnapshotNID, int64, error)
	GetInvitesForUser(ctx context.Context, roomNID types.RoomNID, targetUserNID types.EventStateKeyNID) (senderUserIDs []types.EventStateKeyNID, err error)
	SetRoomAlias(ctx context.Context, alias string, roomID string, creatorUserID string) error
	GetRoomIDForAlias(ctx context.Context, alias string) (string, error)
	GetAliasesForRoomID(ctx context.Context, roomID string) ([]string, error)
	GetCreatorIDForAlias(ctx context.Context, alias string) (string, error)
	RemoveRoomAlias(ctx context.Context, alias string) error
	GetMembership(ctx context.Context, roomNID types.RoomNID, requestSenderUserID string) (membershipEventNID types.EventNID, stillInRoom bool, err error)
	GetMembershipEventNIDsForRoom(ctx context.Context, roomNID types.RoomNID, joinOnly bool) ([]types.EventNID, error)
	EventsFromIDs(ctx context.Context, eventIDs []string) ([]types.Event, error)
	GetRoomVersionForRoom(ctx context.Context, roomNID types.RoomNID) (int64, error)
}

// NewPublicRoomsServerDatabase opens a database connection.
func Open(dataSourceName string) (Database, error) {
	uri, err := url.Parse(dataSourceName)
	if err != nil {
		return postgres.Open(dataSourceName)
	}
	switch uri.Scheme {
	case "postgres":
		return postgres.Open(dataSourceName)
	default:
		return postgres.Open(dataSourceName)
	}
}
