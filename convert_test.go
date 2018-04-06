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
		path.Join(dir, "test", "excels", "convert_test.xlsx"),
	}
	outputFile := path.Join(dir, "test", "output", "convert_test.json")

	c := &converter{}
	c.Convert(inputFiles, outputFile, false)

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
		path.Join(dir, "test", "excels", "convert_test.xlsx"),
		path.Join(dir, "test", "excels", "convert_test2.xlsx"),
	}
	outputFile := path.Join(dir, "test", "output", "convert_test.json")

	c := &converter{}
	c.Convert(inputFiles, outputFile, false)

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
		path.Join(dir, "test", "excels"),
	}
	outputFile := path.Join(dir, "test", "output", "convert_test.json")

	c := &converter{}
	c.Convert(inputDir, outputFile, false)

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
		path.Join(dir, "test", "excels"),
	}
	outputFile := path.Join(dir, "test", "output", "convert_test.json")

	c := &converter{}
	c.ConvertConcurrency(inputDir, outputFile, false)

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

func TestConvertFromOneXlsxIntoOneJsonOnlyHeader(t *testing.T) {
	dir, _ := os.Getwd()
	inputFiles := []string{
		path.Join(dir, "test", "excels", "convert_test.xlsx"),
	}
	outputFile := path.Join(dir, "test", "output", "convert_test.json")

	c := &converter{}
	c.ConvertIntoHeader(inputFiles, outputFile, false)

	bytes, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	result := make(XlsxHeaderMap)
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Fatal(err)
	}

	if _, ok := result["sheet"]; !ok {
		t.Errorf("outputed json should have a key of sheet name")
	}

	headerInfo := result["sheet"]

	except := map[string]ColumnInfo{
		"id": ColumnInfo{
			Index:     0,
			ValueType: "int",
		},
		"characterId": ColumnInfo{
			Index:     1,
			ValueType: "int",
		},
		"name": ColumnInfo{
			Index:     2,
			ValueType: "string",
		},
		"hp": ColumnInfo{
			Index:     3,
			ValueType: "long",
		},
		"mp": ColumnInfo{
			Index:     4,
			ValueType: "long",
		},
		"attack": ColumnInfo{
			Index:     5,
			ValueType: "int",
		},
		"defense": ColumnInfo{
			Index:     6,
			ValueType: "int",
		},
	}

	for k1, v1 := range except {
		if _, ok := headerInfo[k1]; !ok {
			t.Errorf("column not found: %s", k1)
		}

		if v1.Index != headerInfo[k1].Index {
			t.Errorf("invalid column info. column name: %s, attribute: Index, expect: %d, actual: %d", k1, v1.Index, headerInfo[k1].Index)
		}
		if v1.ValueType != headerInfo[k1].ValueType {
			t.Errorf("invalid column info. column name: %s, attribute: ValueType, expect: %s, actual: %s", k1, v1.ValueType, headerInfo[k1].ValueType)
		}
	}
}
