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

	lineOneFormat, _ := conf.GetExcelFormatByLine(1)
	if lineOneFormat.RowType != ExcelFormatRowTypeKey {
		t.Errorf("Format of line:1 should be key")
	}

	lineTwoFormat, _ := conf.GetExcelFormatByLine(2)
	if lineTwoFormat.RowType != ExcelFormatRowTypeValueType {
		t.Errorf("Format of line:2 should be type")
	}

	lineThreeFormat, _ := conf.GetExcelFormatByLine(3)
	if lineThreeFormat.RowType != ExcelFormatRowTypeComment {
		t.Errorf("Format of line:1 should be comment")
	}
}

func TestLoadJsonFormatFromConfig(t *testing.T) {
}
