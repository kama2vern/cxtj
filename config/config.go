package config

import (
	"fmt"

	"../logger"

	"github.com/BurntSushi/toml"
)

var DefaultConfig *Config

// Config represents cxtj's configuration file.
type Config struct {
	ExcelFormats []ExcelFormat `toml:"excel"`

	// TODO: output json config
}

// ExcelFormat represents an input excel format
type ExcelFormat struct {
	RowType ExcelFormatRowType `toml:"row_type"`
	RowLine int                `toml:"row_line"`
}

// ExcelFormatRowType is an enum to represent a type of excel row.
type ExcelFormatRowType int

// ExcelFormatRowType enum values
const (
	ExcelFormatRowTypeData ExcelFormatRowType = iota
	ExcelFormatRowTypeKey
	ExcelFormatRowTypeValueType
	ExcelFormatRowTypeComment
)

func (c ExcelFormatRowType) String() string {
	switch c {
	case ExcelFormatRowTypeData:
		return "data"
	case ExcelFormatRowTypeKey:
		return "key"
	case ExcelFormatRowTypeValueType:
		return "value-type"
	case ExcelFormatRowTypeComment:
		return "comment"
	}
	return ""
}

// UnmarshalText is used by toml unmarshaller
func (c *ExcelFormatRowType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "data", "":
		*c = ExcelFormatRowTypeData
		return nil
	case "key":
		*c = ExcelFormatRowTypeKey
		return nil
	case "value-type":
		*c = ExcelFormatRowTypeValueType
		return nil
	case "comment":
		*c = ExcelFormatRowTypeComment
		return nil
	default:
		*c = ExcelFormatRowTypeData // Avoid panic
		return fmt.Errorf("failed to parse")
	}
}

func init() {
	DefaultConfig = &Config{
		ExcelFormats: []ExcelFormat{
			ExcelFormat{
				RowType: ExcelFormatRowTypeKey,
			},
			ExcelFormat{
				RowType: ExcelFormatRowTypeValueType,
			},
			ExcelFormat{
				RowType: ExcelFormatRowTypeComment,
			},
		},
	}
}

// LoadExcelFormatsFromConfig gets array of excel formats from config file
func LoadExcelFormatsFromConfig(conffile string) []ExcelFormat {
	conf, err := LoadConfigFile(conffile)
	if err != nil {
		return []ExcelFormat{}
	}
	return conf.ExcelFormats
}

// GetExcelFormatByLine finds specific excel format by row line
func (c *Config) GetExcelFormatByLine(line int) (ExcelFormat, error) {
	for _, excelFormat := range c.ExcelFormats {
		if excelFormat.RowLine == line {
			return excelFormat, nil
		}
	}
	return ExcelFormat{}, fmt.Errorf("not found excel format. row_line: %d", line)
}

// GetExcelFormatByRowType finds specific excel format by row type
func (c *Config) GetExcelFormatByRowType(rowType ExcelFormatRowType) (ExcelFormat, error) {
	for _, excelFormat := range c.ExcelFormats {
		if excelFormat.RowType == rowType {
			return excelFormat, nil
		}
	}
	return ExcelFormat{}, fmt.Errorf("not found excel format. row_type: %s", rowType.String())
}

// LoadConfigFile gets Config
func LoadConfigFile(file string) (*Config, error) {
	if len(file) == 0 {
		return DefaultConfig, nil
	}

	config := &Config{}
	if _, err := toml.DecodeFile(file, config); err != nil {
		logger.ErrorIf(err)
		return nil, err
	}

	return config, nil
}
