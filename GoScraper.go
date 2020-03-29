package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/andlabs/ui"
)

var response = ""
var stringmatches = ""
//var cleanmatch []string
var matches []string

func main() {
	err := ui.Main(func() {
		// Set regexp pattern
		var textregex = regexp.MustCompile(`(?m)<p>(.+)`)
		//var nosymbolregex = regexp.MustCompile(`(?m)\w+`)


		// creates textbody box for HTML response
		textbox := ui.NewHorizontalBox()
		textbox.SetPadded(true)
		textbody := ui.NewMultilineEntry()
		textbox.Append(textbody, false)
		// End of textbody

		// Creates and handles url input
		input := ui.NewSearchEntry()
		urlbox := ui.NewVerticalBox()
		urlbox.Append(ui.NewLabel("Insert URL:"), false)
		urlbox.SetPadded(true)
		urlbox.Append(input,false)
		// End of url input

		// Start of Combobox
		combobox := ui.NewCombobox()
		combobox.Append("HTML")
		combobox.Append("TXT")
		combobox.Append("CSV")
		combobox.SetSelected(0)
		vbox := ui.NewHorizontalBox()
		vbox.SetPadded(true)
		vbox.Append(combobox, false)
		// End of Combobox

		// Start of search box
		searchbutton := ui.NewButton("Search")
		searchbox := ui.NewHorizontalBox()
		searchbox.SetPadded(true)
		searchbox.Append(searchbutton, true)
		searchbutton.OnClicked(func(button *ui.Button) {
			response = httpRequest(input.Text())
			for _, match := range textregex.FindAllString(response, -1) {
				matches = append(matches, match)
				stringmatches = strings.Join(matches," ")
				//for _,match := range nosymbolregex.FindAllString(stringmatches,-1) {
				//	cleanmatch = append(cleanmatch, match)
				//	stringmatches = strings.Join(cleanmatch," ")
				//}
			}
			if combobox.Selected() == 0 {
				textbody.SetText(response)
			}	else if combobox.Selected() == 1 {
				textbody.SetText(stringmatches)
			}	else {
				// temporary filler until table extraction is done
				textbody.SetText(response)
			}
		})
		// End of Search Box

		// Start of Save button
		savebutton := ui.NewButton("Save")
		savebox := ui.NewHorizontalBox()
		savebox.SetPadded(true)
		savebox.Append(savebutton, false)
		// End of Savebox

		combobox.OnSelected(func(combobox *ui.Combobox) {
			switch combobox.Selected() {
			case 0:
				textbody.SetText(response)
				savebutton.OnClicked(func(button *ui.Button) {
					linesToWrite := response
					err := ioutil.WriteFile("temp.html", []byte(linesToWrite), 0777)
					if err != nil {
						log.Fatal(err)
					}
				})
			case 1:
				textbody.SetText(stringmatches)
				savebutton.OnClicked(func(button *ui.Button) {
					linesToWrite := stringmatches
					err := ioutil.WriteFile("temp.txt", []byte(linesToWrite), 0777)
					if err != nil {
						log.Fatal(err)
					}
				})
			case 2:
				savebutton.OnClicked(func(button *ui.Button) {
					linesToWrite := response
					err := ioutil.WriteFile("temp.csv", []byte(linesToWrite), 0777)
					if err != nil {
						log.Fatal(err)
					}
				})
			}
		})

		urlbox.Append(vbox,false)
		vbox.Append(searchbox, true)
		urlbox.Append(textbox,false)
		urlbox.Append(savebox,false)

		window := ui.NewWindow("GoScrape", 700, 500, false)
		window.SetMargined(true)
		window.SetChild(urlbox)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()

	})
	if err != nil {
		panic(err)

	}
}


func httpRequest(url string) string {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	response, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responsetext := string(body)
	return responsetext
}