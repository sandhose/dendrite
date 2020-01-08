package orm

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/matrix-org/dendrite/common"
	"github.com/matrix-org/dendrite/federationsender/types"
	"github.com/matrix-org/gomatrixserverlib"
)

type ORM struct {
	common.PartitionOffsetStatements
	db *gorm.DB
}

// NewDatabase opens a new database
func NewDatabase(dataSourceName string) (*ORM, error) {
	var result ORM
	var err error
	if result.db, err = gorm.Open("postgres", dataSourceName); err != nil {
		return nil, err
	}
	if err = result.PartitionOffsetStatements.Prepare(result.db.DB(), "federationsender"); err != nil {
		return nil, err
	}
	result.db.AutoMigrate(&types.Room{}, &types.JoinedHost{})
	return &result, nil
}

func (o *ORM) GetJoinedHosts(ctx context.Context, roomID string) (result []types.JoinedHost, err error) {
	o.db.Where("name = ?", roomID).Find(&result)
	return
}

func (o *ORM) UpdateRoom(ctx context.Context, roomID, oldEventID, newEventID string, addHosts []types.JoinedHost, removeHosts []string) (joinedHosts []types.JoinedHost, err error) {
	err = o.db.Transaction(func(tx *gorm.DB) error {
		// Create the room
		var room types.Room
		tx.FirstOrInit(&room, types.Room{RoomID: roomID})

		// Does the last ID match the new event? Do nothing if so
		if room.LastEventID == newEventID {
			return nil
		}

		// The ID didn't match what we expected, return an error
		if room.LastEventID != oldEventID {
			return types.EventIDMismatchError{
				DatabaseID: room.LastEventID, RoomServerID: oldEventID,
			}
		}

		// Get the joined hosts, as before any updates that follow
		tx.Find(&joinedHosts)

		// Remove the specified hosts
		for _, removeHost := range removeHosts {
			if err := tx.Delete(types.JoinedHost{
				ServerName: gomatrixserverlib.ServerName(removeHost),
			}).Error; err != nil {
				return err
			}
		}

		// Add the specified hosts
		for _, addHost := range addHosts {
			if err := tx.Create(addHost).Error; err != nil {
				return err
			}
		}

		// Update the room in the database
		room.LastEventID = newEventID
		tx.Save(&room)

		return nil
	})
	return
}
