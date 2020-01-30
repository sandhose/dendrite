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
	"time"

	"github.com/matrix-org/dendrite/publicroomsapi/storage/postgres"
	"github.com/matrix-org/dendrite/publicroomsapi/types"
	"github.com/matrix-org/gomatrixserverlib"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

const DHTInterval = time.Second * 10

// PublicRoomsServerDatabase represents a public rooms server database.
type PublicRoomsServerDatabase struct {
	dht *dht.IpfsDHT
	postgres.PublicRoomsServerDatabase
	ourRoomsContext  context.Context             // our current value in the DHT
	ourRoomsCancel   context.CancelFunc          // cancel when we want to expire our value
	foundRooms       map[string]types.PublicRoom // additional rooms we have learned about from the DHT
	foundRoomsMutex  sync.RWMutex                // protects foundRooms
	maintenanceTimer *time.Timer
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
	go provider.ResetDHTMaintenance()
	return &provider, nil
}

func (d *PublicRoomsServerDatabase) GetRoomVisibility(ctx context.Context, roomID string) (bool, error) {
	return d.PublicRoomsServerDatabase.GetRoomVisibility(ctx, roomID)
}

func (d *PublicRoomsServerDatabase) SetRoomVisibility(ctx context.Context, visible bool, roomID string) error {
	d.ResetDHTMaintenance()
	return d.PublicRoomsServerDatabase.SetRoomVisibility(ctx, visible, roomID)
}

func (d *PublicRoomsServerDatabase) CountPublicRooms(ctx context.Context) (int64, error) {
	count, err := d.PublicRoomsServerDatabase.CountPublicRooms(ctx)
	if err != nil {
		return 0, err
	}
	d.foundRoomsMutex.RLock()
	defer d.foundRoomsMutex.RUnlock()
	return count + int64(len(d.foundRooms)), nil
}

func (d *PublicRoomsServerDatabase) GetPublicRooms(ctx context.Context, offset int64, limit int16, filter string) ([]types.PublicRoom, error) {
	rooms, err := d.PublicRoomsServerDatabase.GetPublicRooms(ctx, offset, limit, filter)
	if err != nil {
		return []types.PublicRoom{}, err
	}
	d.foundRoomsMutex.RLock()
	defer d.foundRoomsMutex.RUnlock()
	for _, room := range d.foundRooms {
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (d *PublicRoomsServerDatabase) UpdateRoomFromEvents(ctx context.Context, eventsToAdd []gomatrixserverlib.Event, eventsToRemove []gomatrixserverlib.Event) error {
	return d.PublicRoomsServerDatabase.UpdateRoomFromEvents(ctx, eventsToAdd, eventsToRemove)
}

func (d *PublicRoomsServerDatabase) UpdateRoomFromEvent(ctx context.Context, event gomatrixserverlib.Event) error {
	return d.PublicRoomsServerDatabase.UpdateRoomFromEvent(ctx, event)
}

func (d *PublicRoomsServerDatabase) ResetDHTMaintenance() {
	if d.maintenanceTimer != nil && !d.maintenanceTimer.Stop() {
		<-d.maintenanceTimer.C
	}
	d.Interval()
}

func (d *PublicRoomsServerDatabase) Interval() {
	if err := d.AdvertiseRoomsIntoDHT(); err != nil {
		//	fmt.Println("Failed to advertise room in DHT:", err)
	}
	if err := d.FindRoomsInDHT(); err != nil {
		//	fmt.Println("Failed to find rooms in DHT:", err)
	}
	d.maintenanceTimer = time.AfterFunc(DHTInterval, d.Interval)
}

func (d *PublicRoomsServerDatabase) AdvertiseRoomsIntoDHT() error {
	if d.ourRoomsContext != nil {
		d.ourRoomsCancel()
	}
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer dbCancel()
	if count, err := d.PublicRoomsServerDatabase.CountPublicRooms(dbCtx); err != nil || count == 0 {
		return err
	}
	ourRooms, err := d.GetPublicRooms(dbCtx, 0, 1024, "")
	if err != nil {
		return err
	}
	fmt.Println("Going to advertise", len(ourRooms), "rooms into the DHT")
	if j, err := json.Marshal(ourRooms); err == nil {
		d.ourRoomsContext, d.ourRoomsCancel = context.WithCancel(context.Background())
		if err := d.dht.PutValue(d.ourRoomsContext, "/matrix/publicRooms", j); err != nil {
			return err
		}
	}
	return nil
}

func (d *PublicRoomsServerDatabase) FindRoomsInDHT() error {
	d.foundRoomsMutex.Lock()
	searchCtx, searchCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer searchCancel()
	defer d.foundRoomsMutex.Unlock()
	results, err := d.dht.GetValues(searchCtx, "/matrix/publicRooms", 1024)
	if err != nil {
		return err
	}
	d.foundRooms = make(map[string]types.PublicRoom)
	for _, result := range results {
		var received []types.PublicRoom
		if err := json.Unmarshal(result.Val, &received); err != nil {
			return err
		}
		for _, room := range received {
			/*
				if room.RoomID[len(room.RoomID)-52:] != result.From.String() {
					continue
				}
			*/
			d.foundRooms[room.RoomID] = room
		}
	}
	fmt.Println("Found", len(d.foundRooms), "rooms")
	return nil
}
