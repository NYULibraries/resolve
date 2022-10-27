package sfx

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

const mockTimestamp = "2017-10-27T10:49:40-04:00"

func TestNewMultipleObjectsRequest(t *testing.T) {
	var tests = []struct {
		querystring   url.Values
		expectedError error
	}{
		{map[string][]string{"genre": {"book"}}, errors.New("could not parse OpenURL: no valid querystring values to parse")},
		{map[string][]string{"rft.genre": {"podcast"}}, errors.New("genre is not valid: genre not in list of allowed genres: [podcast]")},
		{map[string][]string{"rft.genre": {"book"}, "rft.aulast": {"<rft:"}}, errors.New("could not convert multiple objects request to XML: request multiple objects XML is not valid XML: <nil>")},
		{map[string][]string{"rft.genre": {"book"}, "rft.btitle": {"dune"}}, nil},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.querystring)
		t.Run(testName, func(t *testing.T) {
			ans, err := NewMultipleObjectsRequest(testCase.querystring)
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("NewMultipleObjectsRequest returned no error, expecting '%v'", testCase.expectedError)
				}
				if err.Error() != testCase.expectedError.Error() {
					t.Errorf("NewMultipleObjectsRequest returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("NewMultipleObjectsRequest returned error '%v', expecting no errors", err)
			}
			if err == nil {
				if !strings.HasPrefix(ans.RequestXML, `<?xml version="1.0" encoding="UTF-8"?>`) {
					t.Errorf("requestXML isn't an XML document")
				}
			}
		})
	}
}

func TestRequestXML(t *testing.T) {
	var tests = []struct {
		name        string
		tpl         multipleObjectsRequestBodyParams
		expectedErr error
	}{
		{"genre=\"book\"; btitle=\"a book\"", multipleObjectsRequestBodyParams{RftValues: &openURL{"genre": {"book"}, "btitle": {"a book"}}, Timestamp: mockTimestamp, Genre: "book"}, nil},
		{"[empty request body]", multipleObjectsRequestBodyParams{}, errors.New("error")},
		{"genre=\"<rft:\"", multipleObjectsRequestBodyParams{RftValues: &openURL{"genre": {"<rft:"}}, Timestamp: mockTimestamp, Genre: "book"}, errors.New("error")},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actualXML, err := requestXML(testCase.tpl)
			if testCase.expectedErr == nil && !strings.HasPrefix(actualXML, `<?xml version="1.0" encoding="UTF-8"?>`) {
				t.Errorf("toRequestXML didn't return an XML document")
			}
			if testCase.expectedErr != nil {
				if err == nil {
					t.Errorf("toRequestXML err was '%v', expecting '%v'", err, testCase.expectedErr)
				}
			}
		})
	}
}

// func (c MultipleObjectsRequest) Do() (body string, err error) {
// func Init(qs url.Values) (MultipleObjectsRequest *MultipleObjectsRequest, err error) {
