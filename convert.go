package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tealeg/xlsx"

	"./config"
	"./logger"
)

type Converter struct {
	config *config.Config
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

// ColumnInfo is the maximum information about one of the column
type ColumnInfo struct {
	Index     int    `json:"index"`
	ValueType string `json:"valueType"`
}

/*
	sheet: {
		column1: {
			Index: 0,
			ValueType: "int",
		},
		column2: {
			Index: 1,
			ValueType: "string",
		},
		column3: {
			Index: 2,
			ValueType: "float",
		},
	}
*/
// SheetColumns is information list of columns
type SheetColumns map[string]ColumnInfo

// XlsxHeaderMap is information list of columns from one xlsx
type XlsxHeaderMap map[string]map[string]ColumnInfo

func (c *Converter) sheet2Map(sheet *xlsx.Sheet) SheetDataList {
	headers := make([]string, len(sheet.Rows[0].Cells))
	for i, c := range sheet.Rows[0].Cells {
		headers[i] = c.Value
	}

	excelFormats := c.config.ExcelFormats

	size := 0
	converts := make(SheetDataList, len(sheet.Rows[len(excelFormats):]))
	for i, r := range sheet.Rows {
		if i < len(excelFormats) && excelFormats[i].RowType != config.ExcelFormatRowTypeData {
			continue
		}

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
				converts[size] = convertMap
				size++
				break
			}
		}
	}

	return converts[:size]
}

func (c *Converter) xlsx2Map(xFile *xlsx.File) XlsxMap {
	resultJSON := XlsxMap{}
	for _, s := range xFile.Sheets {
		resultJSON[s.Name] = c.sheet2Map(s)
	}
	return resultJSON
}

func (c *Converter) sheet2HeaderMap(sheet *xlsx.Sheet) SheetColumns {
	keyExcelFormat, err := c.config.GetExcelFormatByRowType(config.ExcelFormatRowTypeKey)
	logger.DieIf(err)

	valueTypeExcelFormat, err := c.config.GetExcelFormatByRowType(config.ExcelFormatRowTypeValueType)
	logger.DieIf(err)

	headers := make(map[string]ColumnInfo, len(sheet.Rows[keyExcelFormat.RowLine-1].Cells))
	for i, c := range sheet.Rows[keyExcelFormat.RowLine-1].Cells {
		headers[c.Value] = ColumnInfo{
			Index:     i,
			ValueType: sheet.Rows[valueTypeExcelFormat.RowLine-1].Cells[i].Value,
		}
	}
	return headers
}

func (c *Converter) xlsx2HeaderMap(xFile *xlsx.File) XlsxHeaderMap {
	ret := XlsxHeaderMap{}
	for _, s := range xFile.Sheets {
		ret[s.Name] = c.sheet2HeaderMap(s)
	}
	return ret
}

func (c *Converter) convertXlsxFileIntoHeader(filename string) XlsxHeaderMap {
	xlsxFile, err := xlsx.OpenFile(filename)
	if logger.ErrorIf(err) {
		logger.Log("convert.go", fmt.Sprintf("error file: %s", filename))
		return XlsxHeaderMap{}
	}

	return c.xlsx2HeaderMap(xlsxFile)
}

func (c *Converter) mergeXlsxMap(m1 XlsxMap, m2 XlsxMap) XlsxMap {
	ret := XlsxMap{}
	for k, v := range m1 {
		ret[k] = v
	}
	for k, v := range m2 {
		ret[k] = v
	}
	return ret
}

func (c *Converter) mergeXlsxHeaderMap(m1 XlsxHeaderMap, m2 XlsxHeaderMap) XlsxHeaderMap {
	ret := XlsxHeaderMap{}
	for k, v := range m1 {
		ret[k] = v
	}
	for k, v := range m2 {
		ret[k] = v
	}
	return ret
}

func (c *Converter) convertXlsxFile(filename string) XlsxMap {
	xlsxFile, err := xlsx.OpenFile(filename)
	if logger.ErrorIf(err) {
		logger.Log("convert.go", fmt.Sprintf("ignored error file: %s", filename))
		return XlsxMap{}
	}

	return c.xlsx2Map(xlsxFile)
}

func (c *Converter) traversalInputFiles(inputDirsOrFiles []string) []string {
	ret := []string{}
	for _, inputDirOrFile := range inputDirsOrFiles {
		fi, err := os.Stat(inputDirOrFile)
		logger.DieIf(err)

		if fi.IsDir() {
			filepath.Walk(inputDirOrFile, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() && c.isExcelFile(path) {
					ret = append(ret, path)
				}
				return nil
			})
			continue
		}

		if c.isExcelFile(inputDirOrFile) {
			ret = append(ret, inputDirOrFile)
		}
	}
	return ret
}

func (c *Converter) isExcelFile(filename string) bool {
	targetExt := filepath.Ext(filename)
	for _, ext := range c.config.ExcelExts {
		if targetExt == ext {
			return true
		}
	}
	return false
}

// ConvertConcurrency executes as the same logic as Convert in concurrently
func (c *Converter) ConvertConcurrency(inputDirsOrFiles []string, outputFile string, isMultipleOutput bool) {
	resultJSON := XlsxMap{}

	il := c.traversalInputFiles(inputDirsOrFiles)

	resultJSON = DispatchConcurrencyWorkers(il, func(path string) XlsxMap {
		return c.convertXlsxFile(path)
	})

	bytes, err := json.Marshal(resultJSON)
	logger.DieIf(err)

	err = ioutil.WriteFile(outputFile, bytes, 0644)
	logger.DieIf(err)
}

// Convert executes convertion from xlsx files or directories into json file(s)
func (c *Converter) Convert(inputDirsOrFiles []string, outputFile string, isMultipleOutput bool) {
	resultJSON := XlsxMap{}

	for _, inputFile := range c.traversalInputFiles(inputDirsOrFiles) {
		resultJSON = c.mergeXlsxMap(resultJSON, c.convertXlsxFile(inputFile))
	}

	bytes, err := json.Marshal(resultJSON)
	logger.DieIf(err)

	err = ioutil.WriteFile(outputFile, bytes, 0644)
	logger.DieIf(err)
}

// ConvertIntoHeader executes convertion from xlsx files or directories into header only json file(s)
func (c *Converter) ConvertIntoHeader(inputDirsOrFiles []string, outputFile string, isMultipleOutput bool) {
	resultJSON := XlsxHeaderMap{}

	for _, inputFile := range c.traversalInputFiles(inputDirsOrFiles) {
		resultJSON = c.mergeXlsxHeaderMap(resultJSON, c.convertXlsxFileIntoHeader(inputFile))
	}

	bytes, err := json.Marshal(resultJSON)
	logger.DieIf(err)

	err = ioutil.WriteFile(outputFile, bytes, 0644)
	logger.DieIf(err)
}

// NewConverter creates new Converter instance
func NewConverter(conf *config.Config) *Converter {
	var c *config.Config
	if conf == nil {
		c = config.DefaultConfig
	} else {
		c = conf
	}

	ret := &Converter{
		config: c,
	}
	return ret
}
