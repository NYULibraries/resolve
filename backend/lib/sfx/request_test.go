package sfx

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

const mockTimestamp = "2017-10-27T10:49:40-04:00"

func TestValidGenre(t *testing.T) {
	var tests = []struct {
		genre       []string
		expected    string
		expectedErr error
	}{
		{[]string{"book", "book"}, "book", nil},
		{[]string{"journal", "book"}, "journal", nil},
		{[]string{"book"}, "book", nil},
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

func TestParseOpenURL(t *testing.T) {
	var tests = []struct {
		queryString map[string][]string
		expected    map[string][]string
		expectedErr error
	}{
		{map[string][]string{"genre": {"book"}, "rft.genre": {"book"}}, OpenURL{"genre": {"book"}}, nil},
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal", "book"}}, OpenURL{"genre": {"journal", "book"}}, nil},
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal"}}, OpenURL{"genre": {"journal"}}, nil},
		{map[string][]string{"genre": {"book"}}, OpenURL{}, errors.New("error")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.queryString)

		t.Run(testname, func(t *testing.T) {
			ans, err := parseOpenURL(tt.queryString)
			if reflect.DeepEqual(ans, tt.expected) {
				t.Errorf("parseOpenURL returned '%v', expecting '%v'", ans, tt.expected)
			}
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("parseOpenURL err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestIsValidXML(t *testing.T) {
	var tests = []struct {
		testXMLFile string
		expected    bool
	}{
		{"../../testdata/ctxObj_good.xml", true},
		{"../../testdata/ctxObj_bad.xml", false},
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

func TestToRequestXML(t *testing.T) {
	var tests = []struct {
		ctx         *CtxObjReq
		tpl         ctxObjTpl
		expectedErr error
	}{
		{&CtxObjReq{}, ctxObjTpl{RftValues: &OpenURL{"genre": {"book"}, "btitle": {"a book"}}, Timestamp: mockTimestamp, Genre: "book"}, nil},
		{&CtxObjReq{}, ctxObjTpl{}, errors.New("error")},
		{&CtxObjReq{}, ctxObjTpl{RftValues: &OpenURL{"genre": {"<rft:"}}, Timestamp: mockTimestamp, Genre: "book"}, errors.New("error")},
	}

	// Create the templates/index.goxml in the current test context temporarily
	// and delete after the test completes
	err := os.Mkdir("templates", 0755)
	if err != nil {
		t.Errorf("could not create temp templates dir")
	}
	_, err = copy("../../templates/index.goxml", "./templates/index.goxml")
	if err != nil {
		t.Errorf("could not copy template file")
	}
	defer os.RemoveAll("templates")

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.ctx)
		t.Run(testname, func(t *testing.T) {
			c := tt.ctx
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

func TestSetCtxObjReq(t *testing.T) {
	var tests = []struct {
		querystring url.Values
		expected    *CtxObjReq
		expectedErr error
	}{
		{map[string][]string{"genre": {"book"}, "rft.genre": {"book"}, "rft.btitle": {"a book"}}, "book", map[string][]string{"genre": {"book"}, "btitle": {"a book"}}, nil},
		{map[string][]string{}, "book", map[string][]string{"rft:genre": {"book"}, "rft:btitle": {"a book"}}, nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.querystring)
		t.Run(testname, func(t *testing.T) {
			ans, err := setCtxObjReq(tt.querystring)

			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("setCtxObjReq err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
			if err == nil {
				if !strings.HasPrefix(ans.RequestXML, `<?xml version="1.0" encoding="UTF-8"?>`) {
					t.Errorf("requestXML isn't an XML document")
				}
			}
		})
	}
}

// type ContextObjectReq struct {
// 	RftValues map[string][]string
// 	Timestamp string
// 	Genre     string
// }

// func setContextObjectReq(qs url.Values) (ctx *ContextObjectReq, err error) {
// func toJson(from []byte) (to string, err error) {
// func ToCtxObjReq(qs url.Values) (ctxObjReqXml string, err error) {
// func Post(requestXML string) (body string, err error) {

// Util function for copying a file from a source to a new dest
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
