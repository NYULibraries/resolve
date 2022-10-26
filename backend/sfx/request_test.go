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
		// TODO: This was previously a false negative, as the XML-illegal "<" char prevented
		// the creation of the request XML.  Do we want this to pass or fail?
		// Might require testing GetIt and SFX, or just making a decision about what to do.
		// If we want it to pass, we will need to figure out how we want that done.
		// {map[string][]string{"rft.genre": {"book"}, "rft.aulast": {"<rft:"}}, nil},
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

func TestToRequestXML(t *testing.T) {
	var tests = []struct {
		sfxContext  *MultipleObjectsRequest
		tpl         multipleObjectsRequestBody
		expectedErr error
	}{
		{&MultipleObjectsRequest{}, multipleObjectsRequestBody{RftValues: &openURL{"genre": {"book"}, "btitle": {"a book"}}, Timestamp: mockTimestamp, Genre: "book"}, nil},
		{&MultipleObjectsRequest{}, multipleObjectsRequestBody{}, errors.New("error")},
		{&MultipleObjectsRequest{}, multipleObjectsRequestBody{RftValues: &openURL{"genre": {"<rft:"}}, Timestamp: mockTimestamp, Genre: "book"}, errors.New("error")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.sfxContext)
		t.Run(testname, func(t *testing.T) {
			c := tt.sfxContext
			err := c.toRequestXML(tt.tpl)
			if tt.expectedErr == nil && !strings.HasPrefix(c.RequestXML, `<?xml version="1.0" encoding="UTF-8"?>`) {
				t.Errorf("toRequestXML didn't return an XML document")
			}
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("toRequestXML err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
		})
	}
}

// func (c MultipleObjectsRequest) Do() (body string, err error) {
// func Init(qs url.Values) (MultipleObjectsRequest *MultipleObjectsRequest, err error) {
