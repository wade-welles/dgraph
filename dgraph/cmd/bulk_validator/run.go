package bulkvalidator

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/dgraph-io/dgraph/x"
	"github.com/spf13/cobra"
)

var BulkValidator x.SubCommand

func init() {
	BulkValidator.Cmd = &cobra.Command{
		Use:   "bulk_validator",
		Short: "Validate input file",
		Run: func(cmd *cobra.Command, args []string) {
			defer x.StartProfile(BulkValidator.Conf).Stop()
			run()
		},
	}
	BulkValidator.EnvPrefix = "BULK_VALIDATOR"

	flag := BulkValidator.Cmd.Flags()
	flag.StringP("files", "f", "",
		"Location of *.rdf(.gz) or *.json(.gz) files(s) to validate.")
	flag.StringP("schema", "s", "",
		"Location of schema file.")
	flag.IntP("num_go_routines", "j", int(math.Ceil(float64(runtime.NumCPU())/4.0)),
		"Number of worker threads to use. MORE THREADS LEAD TO HIGHER RAM USAGE.")
	flag.String("format", "",
		"Specify file format (rdf or json) instead of getting it from filename.")
}

func run() {
	opt := options{
		DataFiles:     BulkValidator.Conf.GetString("files"),
		SchemaFile:    BulkValidator.Conf.GetString("schema"),
		NumGoroutines: BulkValidator.Conf.GetInt("num_go_routines"),
		DataFormat:    BulkValidator.Conf.GetString("format"),
	}

	x.PrintVersion()
	if opt.SchemaFile == "" {
		fmt.Fprint(os.Stderr, "Schema file must be specified.\n")
		os.Exit(1)
	} else if _, err := os.Stat(opt.SchemaFile); err != nil && os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Schema path(%v) does not exist.\n", opt.SchemaFile)
		os.Exit(1)
	}
	if opt.DataFiles == "" {
		fmt.Fprint(os.Stderr, "RDF or JSON file(s) location must be specified.\n")
		os.Exit(1)
	} else {
		fileList := strings.Split(opt.DataFiles, ",")
		for _, file := range fileList {
			if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Data path(%v) does not exist.\n", file)
				os.Exit(1)
			}
		}
	}

	// fmt.Println(opt.DataFiles)
	loader := newLoader(opt)
	loader.mapStage()
}