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

func main() {
	createUI()
}

//type areaHandler struct {}
//
//	func(h areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
//	return
//}
//
//	func (h areaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
//	return
//}
//
//	func (h areaHandler) MouseCrossed(a *ui.Area, left bool) {
//	return
//}
//
//	func (h areaHandler) DragBroken(a *ui.Area) {
//	return
//}
//
//	func (h areaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
//	return
//}

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

func createUI() {
	err := ui.Main(func() {
		box, searchbutton, input := createSearchBox()
		//scroll := ui.NewScrollingArea(areaHandler{},300,300 )
		hbox := ui.NewHorizontalBox()
		hbox.SetPadded(true)
		textbody := ui.NewMultilineEntry()
		hbox.Append(textbody, false)
		//hbox.Append(scroll, false)
		box.Append(hbox, false)
		window := ui.NewWindow("GoScrape", 500, 500, false)
		window.SetMargined(true)
		window.SetChild(box)
		searchbutton.OnClicked(func(*ui.Button) {
			response := httpRequest(input.Text())
			textbody.SetText(response)
		})
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

func createSearchBox() (*ui.Box, *ui.Button, *ui.Entry) {
	input := ui.NewSearchEntry()
	button := ui.NewButton("Search")
	box := ui.NewVerticalBox()
	box.Append(ui.NewLabel("Insert URL:"), false)
	box.Append(input, false)
	box.Append(button, false)
	return box, button, input
}
