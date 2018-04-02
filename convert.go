package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tealeg/xlsx"

	"./logger"
)

type converter struct{}

// XlsxMap is converted data structure from xlsx file
/*
	{
		#{sheet name}: {
			[
				#{column1 name}: #{row1 column1 value},
				#{column2 name}: #{row1 column2 value},
				#{column3 name}: #{row1 column3 value},
				...
			],
			...
		} // same as SheetDataList
	}
*/
type XlsxMap map[string][]map[string]string

// SheetDataList is converted data structure from one of xlsx sheet
/*
	[
		{
			#{column1 name}: #{row1 column1 value},
			#{column2 name}: #{row1 column2 value},
			#{column3 name}: #{row1 column3 value},
		},
		{
			#{column1 name}: #{row2 column1 value},
			#{column2 name}: #{row2 column2 value},
			#{column3 name}: #{row2 column3 value},
		}
	]
*/
type SheetDataList []map[string]string

// RowMap is one of the record from sheet
/*
	{
		#{column1 name}: #{row1 column1 value},
		#{column2 name}: #{row1 column2 value},
		#{column3 name}: #{row1 column3 value},
	},
*/
type RowMap map[string]string

// ColumnInfo is the maximum information about one of the column
type ColumnInfo struct {
	valueType string
}

// SheetColumns is information list of columns from one xlsxs
type SheetColumns map[string]ColumnInfo

func (c *converter) sheet2Map(sheet *xlsx.Sheet, isOnlyHeader bool) SheetDataList {
	headers := make([]string, len(sheet.Rows[0].Cells))
	for i, c := range sheet.Rows[0].Cells {
		headers[i] = c.Value
	}

	if isOnlyHeader {
		// TODO: it's to hard to specify the type of header data dynamically
		logger.Log("sheet2Map", "work in progress")
		return SheetDataList{}
	}

	size := 0
	converts := make(SheetDataList, len(sheet.Rows[3:]))
	for i, r := range sheet.Rows[3:] {
		convertMap := RowMap{}

		for j := 0; j < len(headers); j++ {
			if j >= len(r.Cells) {
				convertMap[headers[j]] = ""
			} else {
				convertMap[headers[j]] = r.Cells[j].Value
			}
		}

		// ignore row which has all empty values
		for _, v := range convertMap {
			if len(v) > 0 {
				converts[i] = convertMap
				size++
				break
			}
		}
	}

	return converts[:size]
}

func (c *converter) xlsx2Map(xFile *xlsx.File, isOnlyHeader bool) XlsxMap {
	resultJSON := XlsxMap{}
	for _, s := range xFile.Sheets {
		resultJSON[s.Name] = c.sheet2Map(s, isOnlyHeader)
	}
	return resultJSON
}

func (c *converter) mergeXlsxMap(m1 XlsxMap, m2 XlsxMap) XlsxMap {
	ret := XlsxMap{}
	for k, v := range m1 {
		ret[k] = v
	}
	for k, v := range m2 {
		ret[k] = v
	}
	return ret
}

func (c *converter) convertXlsxFile(filename string, isOnlyHeader bool) XlsxMap {
	xlsxFile, err := xlsx.OpenFile(filename)
	if logger.ErrorIf(err) {
		logger.Log("convert.go", fmt.Sprintf("error file: %s", filename))
		return XlsxMap{}
	}

	return c.xlsx2Map(xlsxFile, isOnlyHeader)
}

func (c *converter) ConvertConcurrency(inputFiles []string, outputFile string, isOnlyHeader bool, isMultipleOutput bool) {
	resultJSON := XlsxMap{}

	il := []string{}
	for _, inputFile := range inputFiles {
		fi, err := os.Stat(inputFile)
		logger.DieIf(err)

		if fi.IsDir() {
			filepath.Walk(inputFile, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() && filepath.Ext(path) == ".xlsx" {
					il = append(il, path)
				}
				return nil
			})
		} else {
			il = append(il, inputFile)
		}
	}

	resultJSON = DispatchConcurrencyWorkers(il, func(path string) XlsxMap {
		return c.convertXlsxFile(path, isOnlyHeader)
	})

	bytes, err := json.Marshal(resultJSON)
	logger.DieIf(err)

	err = ioutil.WriteFile(outputFile, bytes, 0644)
	logger.DieIf(err)
}

func (c *converter) Convert(inputFiles []string, outputFile string, isOnlyHeader bool, isMultipleOutput bool) {
	resultJSON := XlsxMap{}
	for _, inputFile := range inputFiles {
		fi, err := os.Stat(inputFile)
		logger.DieIf(err)

		if fi.IsDir() {
			filepath.Walk(inputFile, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() && filepath.Ext(path) == ".xlsx" {
					resultJSON = c.mergeXlsxMap(resultJSON, c.convertXlsxFile(path, isOnlyHeader))
				}
				return nil
			})
		} else {
			resultJSON = c.mergeXlsxMap(resultJSON, c.convertXlsxFile(inputFile, isOnlyHeader))
		}
	}

	bytes, err := json.Marshal(resultJSON)
	logger.DieIf(err)

	err = ioutil.WriteFile(outputFile, bytes, 0644)
	logger.DieIf(err)
}
