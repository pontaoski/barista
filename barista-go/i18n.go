package barista

import (
	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/i18n"
)

func l10n(c commandlib.Context, text string) string {
	return i18n.I18n(schemas["locale"].ReadValue(c), text)
}
