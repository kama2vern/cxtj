package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tealeg/xlsx"

	"./logger"
)

type converter struct {
	InputFiles       []string
	OutputFile       string
	IsOnlyHeader     bool
	IsMultipleOutput bool
}

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

func (c *converter) sheet2Map(sheet *xlsx.Sheet) SheetDataList {
	headers := make([]string, len(sheet.Rows[0].Cells))
	for i, c := range sheet.Rows[0].Cells {
		headers[i] = c.Value
	}

	converts := make(SheetDataList, len(sheet.Rows[1:]))
	for i, r := range sheet.Rows[3:] {
		convertMap := RowMap{}

		for j := 0; j < len(headers); j++ {
			if j >= len(r.Cells) {
				convertMap[headers[j]] = ""
			} else {
				convertMap[headers[j]] = r.Cells[j].Value
			}
		}
		converts[i] = convertMap
	}

	return converts
}

func (c *converter) xlsx2Map(xFile *xlsx.File) XlsxMap {
	resultJSON := XlsxMap{}
	for _, s := range xFile.Sheets {
		resultJSON[s.Name] = c.sheet2Map(s)
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

func (c *converter) convertXlsxFile(filename string) XlsxMap {
	xlsxFile, err := xlsx.OpenFile(filename)
	logger.DieIf(err)

	return c.xlsx2Map(xlsxFile)
}

func (c *converter) Convert() {
	resultJSON := XlsxMap{}
	for _, inputFile := range c.InputFiles {
		fi, err := os.Stat(inputFile)
		logger.DieIf(err)

		if fi.IsDir() {
			filepath.Walk(inputFile, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					resultJSON = c.mergeXlsxMap(resultJSON, c.convertXlsxFile(path))
				}
				return nil
			})
		} else {
			resultJSON = c.mergeXlsxMap(resultJSON, c.convertXlsxFile(inputFile))
		}
	}

	bytes, err := json.Marshal(resultJSON)
	logger.DieIf(err)

	err = ioutil.WriteFile(c.OutputFile, bytes, 0644)
	logger.DieIf(err)
}
