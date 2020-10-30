package detjson

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ErrorMessageCouldNotReadTestData = "could not read test data: "
const identifier = ".json"
const pathTOTestData = "testdata/"

// ReadFile reads a json file.
func ReadFile(fileName string) (string, error) {
	JSONRequest, err := os.Open(pathTOTestData + fileName + identifier)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := JSONRequest.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	RequestByteValue, err := ioutil.ReadAll(JSONRequest)
	if err != nil {
		return "", err
	}
	isValid := json.Valid(RequestByteValue)
	if !isValid {
		return "", errors.New("no valid json: " + fileName)
	}
	return string(RequestByteValue), nil
}

// TestMarshalOrdered_OutputIsValidJSON checks if output is valid JSON
func TestMarshalOrdered_OutputIsValidJSON(t *testing.T) {
	tests := []struct {
		testName string
		fileName string
		valid    bool
	}{
		{
			testName: "isbn",
			fileName: "isbn",
			valid:    true,
		},
		{
			testName: "stack",
			fileName: "stack",
			valid:    true,
		},
		{
			testName: "stack2 mit integers, float und null",
			fileName: "stack2",
			valid:    true,
		},
		{
			testName: "goo mit integers, float und null",
			fileName: "goo",
			valid:    true,
		},
	}
	for _, test := range tests {
		testData, err := ReadFile(test.fileName)
		if err != nil {
			assert.Fail(t, ErrorMessageCouldNotReadTestData, err)
			return
		}
		marshaller := NewMarshaller(testData)
		err = marshaller.UnMarshal()
		if err != nil {
			assert.Fail(t, "could not unmarshall json", err)
			return
		}
		err = marshaller.MarshalOrdered()
		if err != nil {
			assert.Fail(t, ErrorMessageCouldNotReadTestData, err)
			return
		}
		t.Run(test.testName, func(t *testing.T) {
			assert.Equal(t, test.valid, json.Valid([]byte(marshaller.GetJSONString())))
		})
	}
}

// TestMarshalOrdered_OutputIsSorted checks if output is sorted.
func TestMarshalOrdered_OutputIsSorted(t *testing.T) {
	tests := []struct {
		testName string
		fileName string
	}{
		{
			testName: "goo mit integers, float und null",
			fileName: "goo_sorted",
		},
		{
			testName: "isbn",
			fileName: "isbn_sorted",
		},
		{
			testName: "stack",
			fileName: "stack_sorted",
		},
		{
			testName: "stack2 mit integers, float und null",
			fileName: "stack2_sorted",
		},
	}
	for _, test := range tests {
		testData, err := ReadFile(test.fileName)
		if err != nil {
			assert.Fail(t, ErrorMessageCouldNotReadTestData, err)
			return
		}
		marshaller := NewMarshaller(testData)
		err = marshaller.UnMarshal()
		if err != nil {
			assert.Fail(t, "could not unmarshall json", err)
			return
		}
		err = marshaller.MarshalOrdered()
		if err != nil {
			assert.Fail(t, ErrorMessageCouldNotReadTestData, err)
			return
		}
		jsonSorted, err := ReadFile(test.fileName)
		if err != nil {
			assert.Fail(t, ErrorMessageCouldNotReadTestData, err)
			return
		}
		// delete all whitespaces and all new lines
		// this this is necessary to compare with test data with output
		re := regexp.MustCompile(`[ \n]*`)
		// output
		output := string(re.ReplaceAll([]byte(marshaller.GetJSONString()), []byte("")))
		// test data
		testdata := string(re.ReplaceAll([]byte(jsonSorted), []byte("")))
		t.Run(test.testName, func(t *testing.T) {
			assert.Equal(t, testdata, output)
		})
	}
}
