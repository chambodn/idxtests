package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	indexSetup   bool
	fpath, dpath string
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index test results file into Elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("index called")
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
	indexCmd.Flags().BoolVar(&indexSetup, "setup", false, "Create Elasticsearch index")
	indexCmd.Flags().StringVar(&fpath, "file", "", "File to index")
	indexCmd.Flags().StringVar(&dpath, "directory", "", "")
}
