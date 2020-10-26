package commandlib

import (
	"context"
	"fmt"

	"github.com/appadeia/barista/ent"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	gorm.Model
	Key        string
	Value      string
	Kind       int
	LocationID string
}

var StorageDB *ent.Client

func init() {
	var err error
	StorageDB, err = ent.Open("sqlite3", "file:data?_fk=1")
	StorageDB.Schema.Create(context.Background())
	if err != nil {
		panic(err.Error())
	}
}

type Scope int

const (
	Invalid Scope = iota
	User
	Channel
	Community
)

func StoreNote(c Context, key string, value string, scope Scope) {
	var id string
	switch scope {
	case User:
		id = c.AuthorIdentifier()
	case Channel:
		id = c.RoomIdentifier()
	case Community:
		id = c.CommunityIdentifier()
	}
	id = fmt.Sprintf("%s-%s", id, key)

	note, err := StorageDB.Note.Get(context.Background(), id)
	if err != nil {
		if _, ok := err.(*ent.NotFoundError); ok {
			StorageDB.Note.Create().
				SetID(id).
				SetContent(value).
				SaveX(context.Background())
		}

		return
	}

	note.Update().SetContent(value).SaveX(context.Background())

	StorageDB.Note.UpdateOneID(id).
		SetContent(value).
		SaveX(context.Background())
}

func GetNote(c Context, key string, scope Scope) (string, bool) {
	var id string
	switch scope {
	case User:
		id = c.AuthorIdentifier()
	case Channel:
		id = c.RoomIdentifier()
	case Community:
		id = c.CommunityIdentifier()
	}
	id = fmt.Sprintf("%s-%s", id, key)

	data, err := StorageDB.Note.Get(context.Background(), id)
	if data == nil {
		if err != nil {
			println(err.Error())
		}
		return "", false
	}

	return data.Content, true
}
