package main

import (
	"./config"
	"./logger"
	"github.com/urfave/cli"
)

// Commands cli.Command object list
var Commands = []cli.Command{
	commandConvert,
}

var commandConvert = cli.Command{
	Name:      "convert",
	Usage:     "Convert xlsx to json file(s)",
	ArgsUsage: "[--verbose | -v] [--only-header] [--multiple-output] --from <xlsxFileName|xlsxDir> --to <jsonFileName|jsonDir>",
	Description: `
    Convert single or multiple xlsx file to single or multiple json file.
    You can designate multiple xlsx file names/dirs and also json files.
`,
	Action: doConvert,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
		cli.BoolFlag{Name: "only-header", Usage: "Only header output mode"},
		cli.BoolFlag{
			Name:  "multiple-output",
			Usage: "Output multiple json files each xlsx sheets",
		},
		cli.StringSliceFlag{
			Name:  "from",
			Value: &cli.StringSlice{},
			Usage: "Input xlsx files or directory which includes some xlsx files. Multiple choices are allowed.",
		},
		cli.StringFlag{
			Name:  "to",
			Usage: "Output json file or directory. Directory choices required --multiple-output mode.",
		},
	},
}

func doConvert(c *cli.Context) error {
	conffile := c.GlobalString("conf")
	conf, err := config.LoadConfigFile(conffile)
	logger.DieIf(err)

	from := c.StringSlice("from")
	to := c.String("to")
	isOnlyHeader := c.Bool("only-header")
	isMultipleOutput := c.Bool("multiple-output")

	if len(from) < 1 || to == "" {
		cli.ShowCommandHelpAndExit(c, "convert", 1)
	}

	converter := NewConverter(conf)
	if isOnlyHeader {
		converter.ConvertIntoHeader(from, to, isMultipleOutput)
	} else {
		converter.Convert(from, to, isMultipleOutput)
	}

	return nil
}
