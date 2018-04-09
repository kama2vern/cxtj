package config

import (
	"os"
	"path"
	"testing"
)

func TestLoadExcelFormatFromConfig(t *testing.T) {
	dir, _ := os.Getwd()
	conffle := path.Join(dir, "..", "test", "cxtj.conf")

	conf, err := LoadConfigFile(conffle)
	if err != nil {
		panic(err)
	}

	excelExts := conf.ExcelExts
	if len(excelExts) <= 0 {
		t.Errorf("Excel Extension is not found")
	}
	if len(excelExts) > 0 &&
		excelExts[0] != ".xlsx" {
		t.Errorf("Invalid excel ext. expect: .xlsx actual: %s", excelExts[0])
	}

	lineOneFormat, _ := conf.GetExcelFormatByLine(1)
	if lineOneFormat.RowType != ExcelFormatRowTypeKey {
		t.Errorf("Format of line:1 should be key")
	}
	keyFormat, _ := conf.GetExcelFormatByRowType(ExcelFormatRowTypeKey)
	if keyFormat.RowLine != 1 {
		t.Errorf("Key Format: should be row_line: 1")
	}

	lineTwoFormat, _ := conf.GetExcelFormatByLine(2)
	if lineTwoFormat.RowType != ExcelFormatRowTypeValueType {
		t.Errorf("Format of line:2 should be type")
	}
	valueTypeFormat, _ := conf.GetExcelFormatByRowType(ExcelFormatRowTypeValueType)
	if valueTypeFormat.RowLine != 2 {
		t.Errorf("ValueType Format: should be row_line: 2")
	}

	lineThreeFormat, _ := conf.GetExcelFormatByLine(3)
	if lineThreeFormat.RowType != ExcelFormatRowTypeComment {
		t.Errorf("Format of line:1 should be comment")
	}
	commentFormat, _ := conf.GetExcelFormatByRowType(ExcelFormatRowTypeComment)
	if commentFormat.RowLine != 3 {
		t.Errorf("Comment Format: should be row_line: 3")
	}
}

func TestLoadJsonFormatFromConfig(t *testing.T) {
}
