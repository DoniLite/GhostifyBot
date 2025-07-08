package services

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func Browser() {
	// Lance un navigateur Chromium headless
	url := launcher.New().Headless(true).MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// Ouvre la page
	page := browser.MustPage("https://github.com/DoniLite")

	// Exemple : remplir un champ et cliquer
	page.MustElement("#qb-input-query").MustInput("Doni bot")
	page.MustElement("#query-builder-test-results > li:nth-child(1)").Type()

	// Attend qu’un élément apparaisse
	page.MustWaitLoad()
	result := page.MustElement("h1").MustText()

	fmt.Println("Résultat trouvé :", result)
}
