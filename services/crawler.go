package services

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

var (
	pw *playwright.Playwright
)

func init() {
	err := playwright.Install()

	if err != nil {
		panic(err)
	}
}

// func main() {
// 	pw, err := playwright.Run()
// 	if err != nil {
// 		log.Fatalf("could not start playwright: %v", err)
// 	}
// 	browser, err := pw.Chromium.Launch()
// 	if err != nil {
// 		log.Fatalf("could not launch browser: %v", err)
// 	}
// 	page, err := browser.NewPage()
// 	if err != nil {
// 		log.Fatalf("could not create page: %v", err)
// 	}
// 	if _, err = page.Goto("https://news.ycombinator.com"); err != nil {
// 		log.Fatalf("could not goto: %v", err)
// 	}
// 	entries, err := page.Locator(".athing").All()
// 	if err != nil {
// 		log.Fatalf("could not get entries: %v", err)
// 	}
// 	for i, entry := range entries {
// 		title, err := entry.Locator("td.title > span > a").TextContent()
// 		if err != nil {
// 			log.Fatalf("could not get text content: %v", err)
// 		}
// 		fmt.Printf("%d: %s\n", i+1, title)
// 	}
// 	if err = browser.Close(); err != nil {
// 		log.Fatalf("could not close browser: %v", err)
// 	}
// 	if err = pw.Stop(); err != nil {
// 		log.Fatalf("could not stop Playwright: %v", err)
// 	}
// }

func ConnectToYGG() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://www.ygg.re/"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	defer page.Close()
	defer browser.Close()
	defer pw.Stop()
	entry := page.Locator(".main-wrapper #content input[type=checkbox]")
	err = entry.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		log.Fatalf("could not wait for entry: %v", err)
	}
	err = entry.Click()
	if err != nil {
		log.Fatalf("could not click entry: %v", err)
	}
	page.WaitForURL("https://www.ygg.re/auth/login")
}

