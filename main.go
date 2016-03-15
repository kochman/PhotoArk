//go:generate esc -o static.go -prefix static static

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/fatih/structs"
	"github.com/nfnt/resize"
	"gopkg.in/yaml.v2"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

var photosMutex = &sync.Mutex{}
var photos = NewSyncMap()
var dirToMetadata = NewSyncMap()
var cache *Cache
var photoDir string
var cacheDir string
var bind string

type metadata struct {
	Event        string
	Photographer string
	Date         string
	Location     string
	Directory    string `json:"-" structs:"directory"`
}

func handleFilters(w http.ResponseWriter, r *http.Request) {
	keys := []string{"Event", "Photographer", "Date", "Location"}

	filters := make(map[string]map[string]struct{})
	for _, k := range keys {
		filters[k] = make(map[string]struct{})
	}

	for meta := range dirToMetadata.Values() {
		m := structs.Map(meta)
		for _, k := range keys {
			if _, ok := m[k]; ok {
				filters[k][m[k].(string)] = struct{}{}
			}
		}
	}

	jFilters := make(map[string][]string)
	for _, k := range keys {
		jFilters[k] = make([]string, 0)
	}

	for i := range filters {
		for j := range filters[i] {
			jFilters[i] = append(jFilters[i], j)
		}
		sort.Strings(jFilters[i])
		if len(jFilters[i]) == 0 {
			delete(jFilters, i)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(jFilters); err != nil {
		log.Print(err)
		return
	}
}

func handleFilter(w http.ResponseWriter, r *http.Request) {
	filterToEligible := make(map[string][]string)

	filters := make(map[string][]string)
	for filterKey := range r.URL.Query() {
		filters[filterKey] = r.URL.Query()[filterKey]
	}

	// find all directories with matching filters
	for filterKey, filterValues := range filters {
		for dir := range dirToMetadata.Keys() {
			meta := dirToMetadata.Get(dir)
			m := structs.Map(meta)
			for k, v := range m {
				if filterKey == k {
					for _, filterValue := range filterValues {
						if filterValue == v {
							filterToEligible[filterKey] = append(filterToEligible[filterKey], dir)
						}
					}
				}
			}
		}
	}

	// find intersection of all filter matches
	eligible := make(map[string]bool)

	for _, dirs := range filterToEligible {
		for _, dir := range dirs {
			count := 0
			for _, dirTests := range filterToEligible {
				for _, dirTest := range dirTests {
					if dir == dirTest {
						count += 1
					}
				}
			}
			if count == len(filterToEligible) {
				eligible[dir] = true
			}
		}
	}

	eligiblePhotos := make([]string, 0)
	for photo := range photos.Keys() {
		dir := filepath.Dir(photo)
		if _, ok := eligible[dir]; ok {
			eligiblePhotos = append(eligiblePhotos, photo)
		}
	}
	sort.Strings(eligiblePhotos)

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(eligiblePhotos); err != nil {
		log.Print(err)
		return
	}
}

func handleMetadata(w http.ResponseWriter, r *http.Request) {
	parentDir := filepath.Dir(r.URL.Query().Get("photo"))

	meta := dirToMetadata.Get(parentDir)

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(meta); err != nil {
		log.Print(err)
		return
	}
}

func handlePhoto(w http.ResponseWriter, r *http.Request) {
	if filepath.Ext(r.URL.Path) != ".jpg" {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// serve photo
	photoFilepath := filepath.Join(photoDir, r.URL.Path[len("/photo/"):])
	http.ServeFile(w, r, photoFilepath)
}

func handleThumb(w http.ResponseWriter, r *http.Request) {
	if filepath.Ext(r.URL.Path) != ".jpg" {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	key := r.URL.Path[len("/cache/"):]
	thumb := cache.Get(key)

	// serve thumbnail
	http.ServeContent(w, r, key, time.Time{}, bytes.NewReader(thumb))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// pull index.html from esc http.FileSystem
	file, err := FS(false).Open("/index.html")
	if err != nil {
		log.Print(err)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		log.Print(err)
		return
	}
	modTime := stat.ModTime()

	http.ServeContent(w, r, "/index.html", modTime, file)
}

func cacheFillFunc(key string) []byte {
	// create thumbnail

	photoFilepath := filepath.Join(photoDir, key)
	file, err := os.Open(photoFilepath)
	if err != nil {
		log.Print(err)
		return nil
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Print(err)
		return nil
	}
	file.Close()

	m := resize.Thumbnail(500, 1000, img, resize.Bicubic)

	out := &bytes.Buffer{}

	jpeg.Encode(out, m, nil)
	return out.Bytes()
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	dir, filename := filepath.Split(path)
	dir, err = filepath.Rel(photoDir, dir)
	if err != nil {
		return err
	}

	if filename == "metadata.yaml" {
		m := metadata{Directory: dir}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, &m)
		dirToMetadata.Put(dir, m)
	} else if filepath.Ext(filename) == ".jpg" {
		filename = filepath.Join(dir, filename)
		photos.Put(filename, struct{}{})
	}
	return nil
}

func walk() {
	for {
		filepath.Walk(photoDir, walkFunc)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	flag.StringVar(&photoDir, "photoDir", "photos", "path to photos")
	flag.StringVar(&cacheDir, "cacheDir", "cache", "path to photo cache directory")
	var cacheSize = flag.Float64("cacheSize", 100, "photo cache size limit in megabytes")
	flag.StringVar(&bind, "bind", "127.0.0.1:8000", "interface and port to bind to")
	flag.Parse()

	go walk()

	cache = NewCache(cacheDir, *cacheSize, cacheFillFunc)
	go cache.PeriodicallyPrune(time.Second * 30)

	// esc for static content. true uses local files, false uses embedded
	http.HandleFunc("/", handleIndex)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(FS(false))))

	http.HandleFunc("/photo/", handlePhoto)
	http.HandleFunc("/thumb/", handleThumb)

	http.HandleFunc("/api/filters", handleFilters)
	http.HandleFunc("/api/filter", handleFilter)
	http.HandleFunc("/api/metadata", handleMetadata)

	log.Print("listening on ", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Print(err)
	}
}
