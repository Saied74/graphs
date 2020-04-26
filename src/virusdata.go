//this package contains two functions.  One to collect the virus data and
//the other is to analyze and parse the JSON files from the website

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	lexer "github.com/Saied74/Lexer2"
)

//Pick type is used to pick the data item of interest from the website
type Pick struct {
	Date         string
	State        string
	FieldName    string
	FieldValue   string
	InterimFiles Interim
	StateList    []string
	DateList     []string
}

//Interim holds the value of the data by state and the date.  It is a map of
//maps.  The first index is the date.  The second index is the state abbriviation
//the data item is the value for the state as a string
type Interim map[string]map[string]string

//GetData collects data from the covid tracking website and returns a string
func GetData(targetURL string) (*string, error) {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	// Make request
	response, err := client.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("request api failed to %s with error %v",
			targetURL, err)
	}
	defer response.Body.Close()

	// Copy data from the response to a byte buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return nil, fmt.Errorf("Copy into IO buffer failed %v", err)
	}
	inputData := buf.String()
	return &inputData, nil
}

//LexInputData parses the JSON file from the website api call
func (p *Pick) LexInputData(pattern [][]string, inputData *string) error {
	var start, done bool
	// b.pickFile.fieldName = "death"

	item := lexer.Lex(pattern, *inputData)

	for {
		newItem := <-item
		switch newItem.ItemKey {
		case "nodeType":
			start = true
		case "object":
			start = false
		case "EOF":
			done = true
		}
		if start {
			switch newItem.ItemKey {
			case "dateChecked":
				tmpDate := strings.Split(newItem.ItemValue, "T")
				if len(tmpDate) != 2 {
					return fmt.Errorf("encounted badly formatted date %v", newItem.ItemValue)
				}
				p.Date = strings.TrimPrefix(tmpDate[0], "\"")
			case "state":
				p.State = newItem.ItemValue
				p.State = strings.TrimPrefix(p.State, `"`)
				p.State = strings.TrimSuffix(p.State, `"`)
			case p.FieldName:
				p.FieldValue = newItem.ItemValue
				if newItem.ItemValue == "null" {
					p.FieldValue = "0"
				}
			}
		}
		if !start && inSlice(p.State, p.StateList) {
			p.processItem()
		}
		if done {
			break
		}
	}
	return nil
}

//After each JSON element is extracted, it appends data to the Interim map
func (p *Pick) processItem() {
	_, ok := p.InterimFiles[p.Date]
	if ok {
		p.InterimFiles[p.Date][p.State] = p.FieldValue
		return
	}
	p.InterimFiles[p.Date] = map[string]string{
		p.State: p.FieldValue,
	}
	return
}

//InSlice checks to see if the candidate is in the target
func inSlice(candidate string, target []string) bool {
	for _, element := range target {
		if element == candidate {
			return true
		}
	}
	return false
}

//BuildDateIndex picks out the datas from the interimFiles and returns a strings
//of dates
func (p *Pick) BuildDateIndex() []string {
	var dateIndex []string
	for key := range p.InterimFiles {
		dateIndex = append(dateIndex, key)
	}
	sort.Strings(dateIndex)
	return dateIndex
}

//GetPattern readds a csv pattern seperated by "|", and returns a [][]string
//one field for each line and one line from each file so organized.
func GetPattern(fileName string) ([][]string, error) {
	var pattern [][]string
	content, err := ioutil.ReadFile(fileName) //get the whole file
	if err != nil {
		return [][]string{}, fmt.Errorf("open error %v on file %s", err, fileName)
	}
	pat1 := strings.Split(string(content), "\n") //split into lines
	for _, pat2 := range pat1 {                  //scan the lines
		pat3 := strings.Split(pat2, "|") //split into comma seperated fields
		pattern = append(pattern, pat3)  //append to the output
	}
	return pattern, nil
}
