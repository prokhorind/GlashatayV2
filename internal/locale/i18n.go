package locale

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
)

//go:embed lang/*.json
var embeddedLocalesFiles embed.FS

type Translator struct {
	translations map[language.Tag]map[string]string
}

func NewTranslator() (*Translator, error) {
	t := &Translator{
		translations: make(map[language.Tag]map[string]string),
	}
	err := t.loadTranslations()
	return t, err
}

func (t *Translator) loadTranslations() error {
	// Load and unmarshal translation files for supported locales
	supportedLocales := []language.Tag{
		language.Ukrainian,
		// Add more supported locales as needed
	}

	for _, locale := range supportedLocales {
		localeTag := locale
		localeStr := locale.String()
		localeJSON, err := embeddedLocalesFiles.ReadFile(fmt.Sprintf("lang/%s.json", localeStr))
		if err != nil {
			return err
		}

		translations := make(map[string]string)
		if err := json.Unmarshal(localeJSON, &translations); err != nil {
			return err
		}

		t.translations[localeTag] = translations
	}
	return nil
}

func (t *Translator) Get(key string, locale language.Tag) (string, error) {
	translations, ok := t.translations[locale]
	if !ok {
		return "", fmt.Errorf("locale not found: %s", locale)
	}

	translation, ok := translations[key]
	if !ok {
		return "", fmt.Errorf("translation not found for key: %s", key)
	}

	return translation, nil
}
