package sfx

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestNewSFXRequest(t *testing.T) {
	var tests = []struct {
		querystring   url.Values
		expectedError error
	}{
		{map[string][]string{"sid": {"genre=article&isbn=&issn=19447485&title=Community%20Development&volume=49&issue=5&date=20181020&atitle=Can%20community%20task%20groups%20learn%20from%20the%20principles%20of%20group%20therapy?&aulast=Zanbar,%20L.&spage=574&sid=EBSCO:Scopus\\\\u00ae&pid=Zanbar,%20L.edselc.2-52.0-8505573399120181020Scopus\\\\u00ae"}}, nil},
		// TODO: Are we able to generate an error to test?
		// {map[string][]string{}, errors.New("TODO")},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.querystring)
		t.Run(testName, func(t *testing.T) {
			// TODO: Test sfxRequest
			_, err := NewSFXRequest(testCase.querystring)
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("NewSFXRequest returned no error, expecting '%v'", testCase.expectedError)
				} else if err.Error() != testCase.expectedError.Error() {
					t.Errorf("NewSFXRequest returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("NewSFXRequest returned error '%v', expecting no errors", err)
			}
		})
	}
}

func TestFilterOpenURLParams(t *testing.T) {
	var testCases = []struct {
		testName    string
		queryString url.Values //map[string][]string
		expected    url.Values //map[string][]string
	}{
		//{"querystring contains sid", map[string][]string{"sid": {"unicode+garbage"}}, map[string][]string{"rfr_id": {"unicode+garbage"}}},
		//{"querystring doesn't contain sid", map[string][]string{"id": {"unicode+garbage"}}, map[string][]string{"id": {"unicode+garbage"}}},
		{"querystring contains sid", map[string][]string{"sid": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}, map[string][]string{"rfr_id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}},
		{"querystring doesn't contain sid", map[string][]string{"id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}, map[string][]string{"id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			actual := filterOpenURLParams(testCase.queryString)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("filterOpenURLParams returned '%v', expecting '%v'", actual, testCase.expected)
			}
		})
	}
}
