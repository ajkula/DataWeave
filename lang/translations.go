package lang

import (
	"fmt"

	lang "github.com/cloudfoundry/jibber_jabber"
)

const (
	defaultLanguage = "en"
)

var supportedLanguages = []string{
	"en", "fr", // Ajoutez d'autres langues prises en charge ici
}

func getSystemLanguage() string {
	detectedLang, err := lang.DetectLanguage()
	if err != nil {
		fmt.Printf("Error detecting system language: %v. Using default language: %s\n", err, defaultLanguage)
		return defaultLanguage
	}

	for _, supportedLang := range supportedLanguages {
		if detectedLang == supportedLang {
			return supportedLang
		}
	}

	fmt.Printf("Detected language %s is not supported. Using default language: %s\n", detectedLang, defaultLanguage)
	return defaultLanguage
}
