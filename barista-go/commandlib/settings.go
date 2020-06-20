package commandlib

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Setting struct {
	gorm.Model
	Key    string
	Value  string
	RoomID string
}

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("sqlite3", "storage/settings.db")
	if err != nil {
		panic(err.Error())
	}
	DB.AutoMigrate(&Setting{})
}

type Schema struct {
	Name           string
	Description    string
	ID             string
	DefaultValue   string
	PossibleValues []string
}

func (s Schema) ReadValue(c Context) string {
	var set Setting
	DB.Where("room_id = ?", c.RoomIdentifier()).First(&set, "key = ?", s.ID)
	if set.RoomID == "" {
		return s.DefaultValue
	}
	return set.Value
}

func (s Schema) WriteValue(c Context, value string) {
	var set Setting
	DB.Where("room_id = ?", c.RoomIdentifier()).FirstOrCreate(&set, "key = ?", s.ID)
	set.Key = s.ID
	set.Value = value
	set.RoomID = c.RoomIdentifier()
	DB.Save(&set)
	return
}

func (s Schema) ToEmbed(c Context) Embed {
	return Embed{
		Title: EmbedHeader{
			Text: fmt.Sprintf("%s (%s)", s.Name, s.ID),
		},
		Fields: []EmbedField{
			{
				Title: "Value",
				Body:  s.ReadValue(c),
			},
		},
		Body: s.Description,
	}
}

func (s Schema) ValueValid(check string) bool {
	if len(s.PossibleValues) == 0 {
		return true
	}
	for _, value := range s.PossibleValues {
		if value == check {
			return true
		}
	}
	return false
}
