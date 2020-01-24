// Copyright 2017-2018 New Vector Ltd
// Copyright 2019-2020 The Matrix.org Foundation C.I.C.
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

package postgreswithdht

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/matrix-org/dendrite/publicroomsapi/storage/postgres"
	"github.com/matrix-org/dendrite/publicroomsapi/types"
	"github.com/matrix-org/gomatrixserverlib"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

// PublicRoomsServerDatabase represents a public rooms server database.
type PublicRoomsServerDatabase struct {
	dht *dht.IpfsDHT
	postgres.PublicRoomsServerDatabase
	ourRoomsContext context.Context             // our current value in the DHT
	ourRoomsCancel  context.CancelFunc          // cancel when we want to expire our value
	foundRooms      map[string]types.PublicRoom // additional rooms we have learned about from the DHT
	foundRoomsMutex sync.RWMutex                // protects foundRooms
}

// NewPublicRoomsServerDatabase creates a new public rooms server database.
func NewPublicRoomsServerDatabase(dataSourceName string, dht *dht.IpfsDHT) (*PublicRoomsServerDatabase, error) {
	pg, err := postgres.NewPublicRoomsServerDatabase(dataSourceName)
	if err != nil {
		return nil, err
	}
	provider := PublicRoomsServerDatabase{
		dht:                       dht,
		PublicRoomsServerDatabase: *pg,
	}
	return &provider, nil
}

func (d *PublicRoomsServerDatabase) GetRoomVisibility(ctx context.Context, roomID string) (bool, error) {
	return d.PublicRoomsServerDatabase.GetRoomVisibility(ctx, roomID)
}

func (d *PublicRoomsServerDatabase) SetRoomVisibility(ctx context.Context, visible bool, roomID string) error {
	defer d.AdvertiseRoomsIntoDHT(ctx)
	return d.PublicRoomsServerDatabase.SetRoomVisibility(ctx, visible, roomID)
}

func (d *PublicRoomsServerDatabase) CountPublicRooms(ctx context.Context) (int64, error) {
	return d.PublicRoomsServerDatabase.CountPublicRooms(ctx)
}

func (d *PublicRoomsServerDatabase) GetPublicRooms(ctx context.Context, offset int64, limit int16, filter string) ([]types.PublicRoom, error) {
	return d.PublicRoomsServerDatabase.GetPublicRooms(ctx, offset, limit, filter)
}

func (d *PublicRoomsServerDatabase) UpdateRoomFromEvents(ctx context.Context, eventsToAdd []gomatrixserverlib.Event, eventsToRemove []gomatrixserverlib.Event) error {
	defer d.AdvertiseRoomsIntoDHT(ctx)
	return d.PublicRoomsServerDatabase.UpdateRoomFromEvents(ctx, eventsToAdd, eventsToRemove)
}

func (d *PublicRoomsServerDatabase) UpdateRoomFromEvent(ctx context.Context, event gomatrixserverlib.Event) error {
	defer d.AdvertiseRoomsIntoDHT(ctx)
	return d.PublicRoomsServerDatabase.UpdateRoomFromEvent(ctx, event)
}

func (d *PublicRoomsServerDatabase) AdvertiseRoomsIntoDHT(ctx context.Context) error {
	if d.ourRoomsContext != nil {
		d.ourRoomsCancel()
	}
	if count, err := d.PublicRoomsServerDatabase.CountPublicRooms(ctx); err != nil || count == 0 {
		fmt.Println("Not advertising rooms:", err)
		return err
	}
	ourRooms, err := d.GetPublicRooms(ctx, 0, 1024, "")
	if err != nil {
		fmt.Println("Error getting public rooms:", err)
		return err
	}
	fmt.Println("Going to advertise", len(ourRooms), "rooms into the DHT")
	if j, err := json.Marshal(ourRooms); err == nil {
		fmt.Println("Marshalled:", string(j))
		d.ourRoomsContext, d.ourRoomsCancel = context.WithCancel(context.Background())
		d.dht.PutValue(d.ourRoomsContext, "/matrix/publicRooms", j)
	}
	return nil
}
