package main

import (
	"os"
	"path"
	"testing"
)

func TestLoadExcelFormatFromConfig(t *testing.T) {
	dir, _ := os.Getwd()
	conffle := path.Join(dir, "test/cxtj.conf")

	excelFormats := LoadExcelFormatsFromConfig(conffle)

	lineOneFormat := excelFormats[0]
	if lineOneFormat.RowType != ExcelFormatRowTypeKey {
		t.Errorf("Format of line:1 should be key")
	}

	lineTwoFormat := excelFormats[1]
	if lineTwoFormat.RowType != ExcelFormatRowTypeValueType {
		t.Errorf("Format of line:2 should be type")
	}

	lineThreeFormat := excelFormats[2]
	if lineThreeFormat.RowType != ExcelFormatRowTypeComment {
		t.Errorf("Format of line:1 should be comment")
	}
}

func TestLoadJsonFormatFromConfig(t *testing.T) {
}
