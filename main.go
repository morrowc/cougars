// A siple webserver to show cougar angels.
// Displays a simple page with background and foreground images selected
// at random from 2 distinct lists.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	host    = flag.String("host", "127.0.0.1", "Host address to listen on.")
	port    = flag.Int("port", 9010, "Port upon which to listen for requests.")
	docroot = flag.String("docroot", "/prod/cougar-angels.com", "DocumentRoot.")

	// Single template makes the milkshake.
	page = template.Must(template.New("page").Parse(pageTemplate))

	image = []string{
		"/bk/cougar.jpg",
		"/bk/cougar-03.png",
		"/bk/cougarangel.jpg",
		"/bk/cougar-rock.jpg",
		"/bk/crayon.jpg",
		"/bk/garfield.png",
		"/bk/narla.jpg",
		"/bk/snowy.jpg",
		"/bk/kitty.jpg",
		"/bk/stone.jpg",
	}

	background = []string{
		"/bk/christmas_lights_texture_seamless.jpg",
		"/bk/colorful_christmas_lights_seamless_texture.jpg",
		"/bk/green_christmas_lights_out_of_focus_seamless_texture.jpg",
		"/bk/pink_and_purple_christmas_lights_texture_seamless.jpg",
		"/bk/pink_christmas_lights_texture_seamless.jpg",
		"/bk/purple_satin_love_bats.gif",
		"/bk/red_yellow_and_green_stars.gif",
	}

	// Set a random seed object up.
	r = rand.New(rand.NewSource(time.Now().Unix()))
)

type handler struct{}

func newHandler() (*handler, error) {
	return &handler{}, nil
}

func selectRandom(s string) string {
	switch {
	case s == "/bk/img.jpg":
		return image[r.Intn(len(image))]
	case s == "/bk/background.jpg":
		return background[r.Intn(len(background))]
	default:
		return "#fail"
	}
}
func writeFile(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join(*docroot, selectRandom(r.URL.Path))
	if _, err := os.Stat(fp); err != nil {
		fmt.Fprintf(w, "<!-- invalid image: %v -->\n", r.URL.Path)
	}
	// Set a cache-control header prior to sending the file.
	w.Header().Add("Cache-Control", "no-cache")

	http.ServeFile(w, r, fp)
}

func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Background string
		Main       string
	}{
		Background: fmt.Sprintf("/bk/background.jpg?%d", time.Now().Unix()),
		Main:       fmt.Sprintf("/bk/img.jpg?%d", time.Now().Unix()),
	}

	// Set a cache-control header prior to sending the file.
	w.Header().Add("Cache-Control", "no-cache")

	var b bytes.Buffer
	err := page.Execute(&b, data)
	if err != nil {
		fmt.Fprintf(w, "failed to parse template: %v", err)
	}
	fmt.Fprintf(w, b.String())
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/bk"):
		writeFile(w, r)
	case strings.HasPrefix(r.URL.Path, "/"):
		index(w, r)
	}
}

func main() {
	flag.Parse()
	worker, err := newHandler()
	if err != nil {
		fmt.Printf("failed to create server: %v\n", err)
		return
	}

	h := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", *host, *port),
		Handler: worker,
	}
	log.Fatal(h.ListenAndServe())
}
