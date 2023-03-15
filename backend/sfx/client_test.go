package sfx

import (
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestIsValidXML(t *testing.T) {
	var tests = []struct {
		testXMLFile string
		expected    bool
	}{
		{"./testdata/sfx-context-object-valid.xml", true},
		{"./testdata/sfx-context-object-invalid-truncated.xml", false},
		{"", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.testXMLFile)
		t.Run(testname, func(t *testing.T) {
			data, _ := ioutil.ReadFile(tt.testXMLFile)
			ans := isValidXML(data)
			if ans != tt.expected {
				t.Errorf("isValidXML returned '%v', expecting '%v'", ans, tt.expected)
			}
		})
	}
}

func TestValidGenre(t *testing.T) {
	var tests = []struct {
		genre       []string
		expected    string
		expectedErr error
	}{
		{[]string{"book", "book"}, "book", nil},
		{[]string{"journal", "book"}, "journal", nil},
		{[]string{"book"}, "book", nil},
		{[]string{"unknown"}, "unknown", nil},
		{[]string{"wrong"}, "", errors.New("error")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.genre)
		t.Run(testname, func(t *testing.T) {
			ans, err := validGenre(tt.genre)
			if ans != tt.expected {
				t.Errorf("validGenre returned '%v', expecting '%v'", ans, tt.expected)
			}
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("validGenre err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
		})
	}
}