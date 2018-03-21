package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func TestConvertFromOneXlsxIntoOneJson(t *testing.T) {
	dir, _ := os.Getwd()
	inputFiles := []string{
		path.Join(dir, "convert_test.xlsx"),
	}
	outputFile := path.Join(dir, "convert_test.json")

	c := &converter{}
	c.Convert(inputFiles, outputFile, false, false)

	bytes, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	result := make(map[string][]map[string]string)
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Fatal(err)
	}

	if _, ok := result["sheet"]; !ok {
		t.Errorf("outputed json should have a key of sheet name")
	}

	contents, _ := result["sheet"]
	if len(contents) == 0 {
		t.Errorf("contents array is empty")
	}

	except := map[string]string{
		"id":          "1.0",
		"characterId": "1001.0",
		"name":        "アルファ",
		"hp":          "100.0",
		"mp":          "50.0",
		"attack":      "1.0",
		"defense":     "1.0",
	}

	for k, v := range contents[0] {
		if v != except[k] {
			t.Errorf("Mismatch contents. key %s, except %s, actual %s", k, except[k], v)
		}
	}

	exceptMaxSize := 4
	if len(contents) != exceptMaxSize {
		t.Errorf("Invalid contents size (has empty row data?). except %d, actual %d", exceptMaxSize, len(contents))
	}
}

func TestConvertFromMultiXlsxIntoOneJson(t *testing.T) {
	dir, _ := os.Getwd()
	inputFiles := []string{
		path.Join(dir, "convert_test.xlsx"),
		path.Join(dir, "convert_test2.xlsx"),
	}
	outputFile := path.Join(dir, "convert_test.json")

	c := &converter{}
	c.Convert(inputFiles, outputFile, false, false)

	bytes, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	result := make(map[string][]map[string]string)
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatal(err)
	}

	for _, sheetName := range []string{"sheet", "nextSheet"} {
		if _, ok := result[sheetName]; !ok {
			t.Errorf("outputed json should have a key of sheet name %s", sheetName)
		}
		contents, _ := result[sheetName]
		if len(contents) == 0 {
			t.Errorf("contents array is empty %s", sheetName)
		}
	}
}

func TestConvertFromOneXlsxDirIntoOneJson(t *testing.T) {
	dir, _ := os.Getwd()
	inputDir := []string{
		path.Join(dir, "test"),
	}
	outputFile := path.Join(dir, "convert_test.json")

	c := &converter{}
	c.Convert(inputDir, outputFile, false, false)

	bytes, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	result := make(map[string][]map[string]string)
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatal(err)
	}

	for _, sheetName := range []string{"sheet", "nextSheet"} {
		if _, ok := result[sheetName]; !ok {
			t.Errorf("outputed json should have a key of sheet name %s", sheetName)
		}
		contents, _ := result[sheetName]
		if len(contents) == 0 {
			t.Errorf("contents array is empty %s", sheetName)
		}
	}
}

func TestConcurrencyConvertFromOneXlsxDirIntoOneJson(t *testing.T) {
	dir, _ := os.Getwd()
	inputDir := []string{
		path.Join(dir, "test", "concurrency"),
	}
	outputFile := path.Join(dir, "convert_test.json")

	c := &converter{}
	c.ConvertConcurrency(inputDir, outputFile, false, false)

	bytes, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	result := make(map[string][]map[string]string)
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatal(err)
	}

	for _, sheetName := range []string{"sheet", "nextSheet"} {
		if _, ok := result[sheetName]; !ok {
			t.Errorf("outputed json should have a key of sheet name %s", sheetName)
		}
		contents, _ := result[sheetName]
		if len(contents) == 0 {
			t.Errorf("contents array is empty %s", sheetName)
		}
	}
}
