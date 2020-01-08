// Copyright 2017 Vector Creations Ltd
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

package types

import (
	"fmt"

	"github.com/matrix-org/gomatrixserverlib"
)

// A JoinedHost is a server that is joined to a matrix room.
type JoinedHost struct {
	RoomID        string                       `gorm:"PRIMARY_KEY"` // The room ID
	MemberEventID string                       `gorm:"PRIMARY_KEY"` // The MemberEventID of a m.room.member join event.
	ServerName    gomatrixserverlib.ServerName `gorm:""`            // The domain part of the state key of the m.room.member join event
}

type Room struct {
	RoomID      string `gorm:"PRIMARY_KEY"` // The room ID
	LastEventID string `gorm:"NOT NULL"`    // Last event in the room
}

// A EventIDMismatchError indicates that we have got out of sync with the
// room server.
type EventIDMismatchError struct {
	DatabaseID   string // The event ID we have stored in our local database.
	RoomServerID string // The event ID received from the room server.
}

func (e EventIDMismatchError) Error() string {
	return fmt.Sprintf(
		"mismatched last sent event ID: had %q in database got %q from room server",
		e.DatabaseID, e.RoomServerID,
	)
}
