package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	ErrEmptyAuthorityList = fmt.Errorf("authority list is empty")
	ErrEmptyMethodList    = fmt.Errorf("method list is empty")
	ErrEmptyURIList       = fmt.Errorf("URI list is empty")
)

// TestCaseYAML define the list of TestCase
type TestCaseYAML struct {
	TestCases []*TestCase `yaml:"testCases"`
}

// TestCase defines the API for declaring unit tests
type TestCase struct {
	Description string       `yaml:"description"`
	Request     *Request     `yaml:"request"`
	Destination *Destination `yaml:"destination"`
}

// Request define the crafted http request present in the test case file.
type Request struct {
	Authority []string          `yaml:"authority"`
	Method    []string          `yaml:"method"`
	URI       []string          `yaml:"uri"`
	Headers   map[string]string `yaml:"headers"`
}

// Input contains the data structure which will be used to assert
type Input struct {
	Authority string
	Method    string
	URI       string
	Headers   map[string]string
}

// Destination define the destination we should assert
type Destination struct {
	Host string `yaml:"host"`
	Port Port   `yaml:"port"`
}

// Port define the port of a given Destination
type Port struct {
	Number int16 `yaml:"number"`
}

// Unfold returns a list of Input objects constructed by all possibilities defined in the Request object. Ex:
// Request{Authority: {"www.example.com", "example.com"}, Method: {"GET", "OPTIONS"}}
// returns []Input{
// 	{Authority:"www.example.com", Method: "GET"},
// 	{Authority:"www.example.com", Method: "OPTIONS"}
// 	{Authority:"example.com", Method: "GET"},
// 	{Authority:"example.com", Method: "OPTIONS"},
// }
func (r *Request) Unfold() ([]Input, error) {
	out := []Input{}

	if len(r.Authority) == 0 {
		return out, ErrEmptyAuthorityList
	}
	if len(r.Method) == 0 {
		return out, ErrEmptyMethodList
	}
	if len(r.URI) == 0 {
		return out, ErrEmptyURIList
	}

	for _, auth := range r.Authority {
		for _, method := range r.Method {
			for _, uri := range r.URI {
				out = append(out, Input{Authority: auth, Method: method, URI: uri, Headers: r.Headers})
			}
		}
	}

	return out, nil
}

func parseTestCases(rootDir string) ([]*TestCase, error) {
	out := []*TestCase{}
	err := filepath.Walk(rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			yamlFile := &TestCaseYAML{}
			fileContet, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			err = yaml.Unmarshal(fileContet, yamlFile)
			if err != nil {
				return err
			}

			if len(yamlFile.TestCases) == 0 {
				return nil
			}

			out = append(out, yamlFile.TestCases...)

			return nil
		})
	if err != nil {
		return nil, err
	}
	return out, nil
}
