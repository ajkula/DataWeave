package main

import (
	"embed"
	"fmt"

	lang "github.com/cloudfoundry/jibber_jabber"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

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

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "db_meta",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
