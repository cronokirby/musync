package internal

import (
	"reflect"
	"testing"
)

const DOC = `
[[source]]
  name = "Scars of Love"
  artist = "FRACTIONS"
  url = "https://www.youtube.com/watch?v=E899THEd3Ao"
  path = "industrial/"
  namestamps = ["Millennials","Scars of Love"]
  timestamps = ["0:00","5:16"]

[[source]]
  name = "1982"
  artist = "Haircuts for Men"
  url = "https://www.youtube.com/watch?v=HSp0E0kCzVc"
  path = "vapor/chirpy/"
  timestamps = ["00:00"]
  namestamps = ["Henry's Lunch Money"]
`

func TestLoadLibrary(t *testing.T) {
	lib, err := LoadLibrary(DOC)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := Library{
		Sources: []Source{
			Source{
				Name:   "Scars of Love",
				Artist: "FRACTIONS",
				URL:    "https://www.youtube.com/watch?v=E899THEd3Ao",
				Path:   "industrial/",
				Sections: []Section{
					Section{Name: "Millennials", HasEnd: true, StartTime: "0:00", EndTime: "5:16"},
					Section{Name: "Scars of Love", HasEnd: false, StartTime: "5:16"},
				},
			},
			Source{
				Name:   "1982",
				Artist: "Haircuts for Men",
				URL:    "https://www.youtube.com/watch?v=HSp0E0kCzVc",
				Path:   "vapor/chirpy/",
				Sections: []Section{
					Section{Name: "Henry's Lunch Money", HasEnd: false, StartTime: "00:00"},
				},
			},
		},
	}
	if !reflect.DeepEqual(*lib, expected) {
		t.Errorf("Expected %+v but got %+v", expected, *lib)
	}
}
