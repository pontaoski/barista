package barista

import badger "github.com/dgraph-io/badger"

// MemKey : Get key value prefixed with unique member information
func (cmd *LexedCommand) MemKey(key string) string {
	return cmd.CommandMessage.Author.ID + key
}

// GetGlobalKey : Get key stored by entire bot.
func GetGlobalKey(key string) string {
	ret := ""
	db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("global" + key))
		if err != nil {
			ret = ""
			return err
		}
		item.Value(func(val []byte) error {
			ret = string(val)
			return nil
		})
		return nil
	})
	return ret
}

// SetGlobalKey : Set key stored by entire bot.
func SetGlobalKey(key string, value string) {
	db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("global"+key), []byte(value))
		return err
	})
}

// GetGuildKey : Get key stored by guild.
func (cmd *LexedCommand) GetGuildKey(key string) string {
	var ret string
	db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(cmd.CommandMessage.GuildID + key))
		if err != nil {
			ret = ""
			return err
		}
		item.Value(func(val []byte) error {
			ret = string(val)
			return nil
		})
		return nil
	})
	return ret
}

// SetGuildKey : Set key stored by guild.
func (cmd *LexedCommand) SetGuildKey(key string, value string) {
	db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(cmd.CommandMessage.GuildID+key), []byte(value))
		return err
	})
}
