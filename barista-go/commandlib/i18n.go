package commandlib

var i18nschema = Schema{
	Name:           "Preferred Locale",
	Description:    "The preferred language of this channel.",
	ID:             "locale",
	DefaultValue:   "en",
	PossibleValues: []string{"en", "de", "es", "fr", "it", "nl", "pl", "tpo"},
}

func GetLanguage(context Context) string {
	return i18nschema.ReadValue(context)
}
