package i18n

var i18n = map[string]map[string]string{
	"tpo": {
		"Next":     "tawa kama",
		"Previous": "tawa pini",
	},
	"es": {
		"Next":                        "Siguiente",
		"Previous":                    "Anterior",
		"Package":                     "Paquete",
		"%s %d out of %d":             "%s %d de %d",
		"%s Package Search":           "Búsqueda de paquetes %s",
		"Name":                        "Nombre",
		"Description":                 "Descripción",
		"Version":                     "Versión",
		"Download Size":               "Tamaño de descarga",
		"Install Size":                "Tamaño de instalación",
		"Search results for %s in %s": "Resultados de la búsqueda para %s en %s",
		"%d packages found":           "%d paquetes encontrados",
		"No packages were found.":     "No se encontraron paquetes.",
		"There was an issue searching for packages: ": "Hay un error: ",
	},
}

func I18n(locale, text string) string {
	if val, ok := i18n[locale]; ok {
		if str, ok := val[text]; ok {
			return str
		}
	}
	return text
}
