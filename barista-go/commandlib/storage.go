package commandlib

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Data struct {
	gorm.Model
	Key        string
	Value      string
	Kind       int
	LocationID string
}

var StorageDB *gorm.DB

func init() {
	var err error
	StorageDB, err = gorm.Open("sqlite3", "storage/data.db")
	if err != nil {
		panic(err.Error())
	}
	StorageDB.AutoMigrate(&Data{})
}

type Scope int

const (
	Invalid Scope = iota
	User
	Channel
	Community
	Global
)

func StoreData(c Context, key string, value string, scope Scope) {
	var id string
	switch scope {
	case User:
		id = c.AuthorIdentifier()
	case Channel:
		id = c.RoomIdentifier()
	case Community:
		id = c.CommunityIdentifier()
	case Global:
		id = "__global__"
	}

	var data Data
	StorageDB.Where("kind = ?", int(scope)).Where("location_id = ?", id).FirstOrCreate(&data, "key = ?", key)
	data.Kind = int(scope)
	data.Key = key
	data.Value = value
	data.LocationID = id
	StorageDB.Save(&data)
}

func RecallData(c Context, key string, scope Scope) string {
	var id string
	switch scope {
	case User:
		id = c.AuthorIdentifier()
	case Channel:
		id = c.RoomIdentifier()
	case Community:
		id = c.CommunityIdentifier()
	}

	var data Data
	StorageDB.Where("kind = ?", int(scope)).Where("location_id = ?", id).First(&data, "key = ?", key)
	return data.Value
}
