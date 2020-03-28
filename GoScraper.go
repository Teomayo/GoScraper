package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andlabs/ui"
)

var response = "string"

func main() {
	err := ui.Main(func() {

		// creates textbody box for HTML response
		textbox := ui.NewHorizontalBox()
		textbox.SetPadded(true)
		textbody := ui.NewMultilineEntry()
		textbox.Append(textbody, false)
		// End of textbody

		// Creates and handles search box logic
		input := ui.NewSearchEntry()
		searchbutton := ui.NewButton("Search")
		urlbox := ui.NewVerticalBox()
		urlbox.Append(ui.NewLabel("Insert URL:"), false)
		urlbox.Append(input, false)
		urlbox.Append(searchbutton, false)
		searchbutton.OnClicked(func(button *ui.Button) {
			response = httpRequest(input.Text())
			textbody.SetText(response)
		})
		// End of Search Box

		// Start of Combobox
		combobox := ui.NewCombobox()
		combobox.Append("HTML")
		combobox.Append("TXT")
		combobox.Append("CSV")
		combobox.SetSelected(0)
		vbox := ui.NewVerticalBox()
		vbox.SetPadded(true)
		vbox.Append(combobox, false)
		// End of Combobox

		// Start of Save button
		savebutton := ui.NewButton("Save")
		savebox := ui.NewHorizontalBox()
		savebox.Append(savebutton, false)

		switch combobox.Selected() {
		case 0:
			savebutton.OnClicked(func(button *ui.Button) {
				linesToWrite := response
				err := ioutil.WriteFile("temp.html", []byte(linesToWrite), 0777)
				if err != nil {
					log.Fatal(err)
				}
			})
		case 1:
			log.Fatal()
		case 2:
			log.Fatal()
		}

		urlbox.Append(vbox, false)
		urlbox.Append(textbox, false)
		urlbox.Append(savebox, false)

		window := ui.NewWindow("GoScrape", 500, 500, false)
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

	// Copy data from the response to standard output
	n, err := io.Copy(os.Stdout, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Number of bytes copied to STDOUT:", n)

	responsetext := string(body)
	return responsetext
}