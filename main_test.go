package main

import (
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
)

func TestEmptyFilters(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	handleFilters(w, req)

	Expect(w.Code).To(Equal(200))
	Expect(len(w.HeaderMap["Content-Type"])).To(Equal(1))
	Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	Expect(w.Body).To(MatchJSON("{}"))
}

func TestFilters(t *testing.T) {
	RegisterTestingT(t)

	// set up test filter
	m := metadata{
		Directory:    "Test Directory",
		Photographer: "Test Photographer",
		Date:         "Test Date",
		Location:     "Test Location",
		Event:        "Test Event",
	}
	dirToMetadata.Put("Test Directory", m)

	// // set up test photos
	// for i := 1; i < 10; i++ {
	// 	filename := "photo" + string(i) + ".jpg"
	// 	filename = filepath.Join("Test Directory", filename)
	// 	photos[filename] = struct{}{}
	// }

	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	handleFilters(w, req)

	Expect(w.Code).To(Equal(200))
	Expect(len(w.HeaderMap["Content-Type"])).To(Equal(1))
	Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	Expect(w.Body).To(MatchJSON(`
		{
			"Date": [
				"Test Date"
			],
			"Event": [
				"Test Event"
			],
			"Photographer": [
				"Test Photographer"
			],
			"Location": [
				"Test Location"
			]
		}
	`))

	// add more test filters
	m = metadata{
		Directory:    "Test Directory 2",
		Photographer: "Test Photographer 2",
		Date:         "Test Date",
		Location:     "Test Location",
		Event:        "Test Event 2",
	}
	dirToMetadata.Put("Test Directory 2", m)

	m = metadata{
		Directory:    "Test Directory 3",
		Photographer: "Test Photographer 2",
		Date:         "Test Date",
		Location:     "Test Location 2",
		Event:        "Test Event",
	}
	dirToMetadata.Put("Test Directory 3", m)

	req, _ = http.NewRequest("GET", "", nil)
	w = httptest.NewRecorder()
	handleFilters(w, req)

	Expect(w.Code).To(Equal(200))
	Expect(len(w.HeaderMap["Content-Type"])).To(Equal(1))
	Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	Expect(w.Body).To(MatchJSON(`
		{
			"Date": [
				"Test Date"
			],
			"Event": [
				"Test Event",
				"Test Event 2"
			],
			"Photographer": [
				"Test Photographer",
				"Test Photographer 2"
			],
			"Location": [
				"Test Location",
				"Test Location 2"
			]
		}
	`))
}

func TestNoFilter(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	handleFilter(w, req)

	Expect(w.Code).To(Equal(200))
	Expect(len(w.HeaderMap["Content-Type"])).To(Equal(1))
	Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	Expect(w.Body).To(MatchJSON("[]"))
}

func TestFilter(t *testing.T) {
	RegisterTestingT(t)

	// set up test photos
	for i := 1; i < 10; i++ {
		filename := "photo" + strconv.Itoa(i) + ".jpg"
		filename = filepath.Join("Test Directory", filename)
		photos.Put(filename, struct{}{})
	}
	for i := 10; i < 20; i++ {
		filename := "photo" + strconv.Itoa(i) + ".jpg"
		filename = filepath.Join("Test Directory 2", filename)
		photos.Put(filename, struct{}{})
	}
	for i := 20; i < 30; i++ {
		filename := "photo" + strconv.Itoa(i) + ".jpg"
		filename = filepath.Join("Test Directory 3", filename)
		photos.Put(filename, struct{}{})
	}

	req, _ := http.NewRequest("GET", "?Photographer=Test+Photographer", nil)
	w := httptest.NewRecorder()
	handleFilter(w, req)

	Expect(w.Code).To(Equal(200))
	Expect(len(w.HeaderMap["Content-Type"])).To(Equal(1))
	Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	Expect(w.Body).To(MatchJSON(`
		["Test Directory/photo1.jpg", "Test Directory/photo2.jpg", "Test Directory/photo3.jpg", "Test Directory/photo4.jpg", "Test Directory/photo5.jpg", "Test Directory/photo6.jpg", "Test Directory/photo7.jpg", "Test Directory/photo8.jpg", "Test Directory/photo9.jpg"]
	`))

	req, _ = http.NewRequest("GET", "?Photographer=Test+Photographer&Photographer=Test+Photographer%202", nil)
	w = httptest.NewRecorder()
	handleFilter(w, req)

	Expect(w.Code).To(Equal(200))
	Expect(len(w.HeaderMap["Content-Type"])).To(Equal(1))
	Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	Expect(w.Body).To(MatchJSON(`
		[
			"Test Directory 2/photo10.jpg",
			"Test Directory 2/photo11.jpg",
			"Test Directory 2/photo12.jpg",
			"Test Directory 2/photo13.jpg",
			"Test Directory 2/photo14.jpg",
			"Test Directory 2/photo15.jpg",
			"Test Directory 2/photo16.jpg",
			"Test Directory 2/photo17.jpg",
			"Test Directory 2/photo18.jpg",
			"Test Directory 2/photo19.jpg",
			"Test Directory 3/photo20.jpg",
			"Test Directory 3/photo21.jpg",
			"Test Directory 3/photo22.jpg",
			"Test Directory 3/photo23.jpg",
			"Test Directory 3/photo24.jpg",
			"Test Directory 3/photo25.jpg",
			"Test Directory 3/photo26.jpg",
			"Test Directory 3/photo27.jpg",
			"Test Directory 3/photo28.jpg",
			"Test Directory 3/photo29.jpg",
			"Test Directory/photo1.jpg",
			"Test Directory/photo2.jpg",
			"Test Directory/photo3.jpg",
			"Test Directory/photo4.jpg",
			"Test Directory/photo5.jpg",
			"Test Directory/photo6.jpg",
			"Test Directory/photo7.jpg",
			"Test Directory/photo8.jpg",
			"Test Directory/photo9.jpg"
		]
	`))

	req, _ = http.NewRequest("GET", "?Photographer=Test+Photographer&Photographer=Test%20Photographer%202&Location=Test%20Location%202", nil)
	w = httptest.NewRecorder()
	handleFilter(w, req)

	Expect(w.Code).To(Equal(200))
	Expect(len(w.HeaderMap["Content-Type"])).To(Equal(1))
	Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	Expect(w.Body).To(MatchJSON(`
		[
			"Test Directory 3/photo20.jpg",
			"Test Directory 3/photo21.jpg",
			"Test Directory 3/photo22.jpg",
			"Test Directory 3/photo23.jpg",
			"Test Directory 3/photo24.jpg",
			"Test Directory 3/photo25.jpg",
			"Test Directory 3/photo26.jpg",
			"Test Directory 3/photo27.jpg",
			"Test Directory 3/photo28.jpg",
			"Test Directory 3/photo29.jpg"
		]
	`))

	req, _ = http.NewRequest("GET", "?Photographer=Test+Photographer&Photographer=Test+Photographer%202&Location=Test+Location&Location=Test+Location+2", nil)
	w = httptest.NewRecorder()
	handleFilter(w, req)

	Expect(w.Code).To(Equal(200))
	Expect(len(w.HeaderMap["Content-Type"])).To(Equal(1))
	Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json"))
	Expect(w.Body).To(MatchJSON(`
		[
			"Test Directory 2/photo10.jpg",
			"Test Directory 2/photo11.jpg",
			"Test Directory 2/photo12.jpg",
			"Test Directory 2/photo13.jpg",
			"Test Directory 2/photo14.jpg",
			"Test Directory 2/photo15.jpg",
			"Test Directory 2/photo16.jpg",
			"Test Directory 2/photo17.jpg",
			"Test Directory 2/photo18.jpg",
			"Test Directory 2/photo19.jpg",
			"Test Directory 3/photo20.jpg",
			"Test Directory 3/photo21.jpg",
			"Test Directory 3/photo22.jpg",
			"Test Directory 3/photo23.jpg",
			"Test Directory 3/photo24.jpg",
			"Test Directory 3/photo25.jpg",
			"Test Directory 3/photo26.jpg",
			"Test Directory 3/photo27.jpg",
			"Test Directory 3/photo28.jpg",
			"Test Directory 3/photo29.jpg",
			"Test Directory/photo1.jpg",
			"Test Directory/photo2.jpg",
			"Test Directory/photo3.jpg",
			"Test Directory/photo4.jpg",
			"Test Directory/photo5.jpg",
			"Test Directory/photo6.jpg",
			"Test Directory/photo7.jpg",
			"Test Directory/photo8.jpg",
			"Test Directory/photo9.jpg"
		]
	`))
}
