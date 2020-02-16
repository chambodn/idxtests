package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/joshdk/go-junit"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	indexSetup   bool
	fpath, dpath string
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "index test results file into Elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {
		worker := TestResultProcessor{
			log: zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
				Level(func() zerolog.Level {
					if os.Getenv("DEBUG") != "" {
						return zerolog.DebugLevel
					}
					return zerolog.InfoLevel
				}()).
				With().
				Timestamp().
				Logger(),
			path: args[0],
		}

		cfg := elasticsearch.Config{
			Addresses: ElasticsearchUrls,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   10,
				ResponseHeaderTimeout: time.Second,
			},
		}
		es, err := elasticsearch.NewClient(cfg)
		if err != nil {
			worker.log.Fatal().Err(err).Msg("Error creating Elasticsearch client")
		}

		config := ResultsConfig{Client: es, IndexName: IndexName}
		res := NewResults(config)

		worker.results = res

		if indexSetup {
			worker.log.Info().Msg("Creating index with mapping")
			if err := worker.setupIndex(); err != nil {
				worker.log.Fatal().Err(err).Msg("Cannot create Elasticsearch index")
			}
		}

		worker.log.Info().Msgf("Starting the processing for <%s>", args[0])
		worker.Run()
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a path or a file to index")
		}
		if len(args) > 1 {
			return errors.New("require only one argument representing file or directory to index")
		}
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return errors.New("this path is not valid")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
	indexCmd.Flags().BoolVar(&indexSetup, "setup", false, "Create Elasticsearch index")
}

// TestResultProcessor allow to process files or directory
// containing xml test results and index them into elasticsearch
type TestResultProcessor struct {
	results *Results
	log     zerolog.Logger

	indexName string
	path      string
}

// Run launches the processor
func (p *TestResultProcessor) Run() {
	rand.Seed(time.Now().Unix())
	statPath, err := os.Stat(p.path)
	if err != nil {
		p.log.Fatal().Str("path", p.path).Err(err).Msg("Cannot retrieve stat from path")
		return
	}
	var suites []junit.Suite
	if statPath.IsDir() {
		suites, err = junit.IngestDir(p.path)
	} else {
		suites, err = junit.IngestFile(p.path)
	}
	if err != nil {
		p.log.Fatal().Str("path", p.path).Err(err).Msg("Cannot read files. Invalid format")
		return
	}
	for _, suite := range suites {
		fmt.Println(suite.Name)
		for _, test := range suite.Tests {
			fmt.Printf("  %s\n", test.Name)
			if test.Error != nil {
				fmt.Printf("    %s: %s\n", test.Status, test.Error.Error())
			} else {
				fmt.Printf("    %s\n", test.Status)
			}
		}
	}
}

func (p *TestResultProcessor) setupIndex() error {
	mapping := `{
		"mappings": {
		  "_doc": {
			"properties": {
			  "id":         { "type": "keyword" },
			  "image_url":  { "type": "keyword" },
			  "title":      { "type": "text", "analyzer": "english" },
			  "alt":        { "type": "text", "analyzer": "english" },
			  "transcript": { "type": "text", "analyzer": "english" },
			  "published":  { "type": "date" },
			  "link":       { "type": "keyword" },
			  "news":       { "type": "text", "analyzer": "english" }
			}
		  }
		}
			}`
	return p.results.CreateIndex(mapping)
}
