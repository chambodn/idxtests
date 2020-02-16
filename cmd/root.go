package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const layoutISO = "2006-01-02"

var cfgFile string

// IndexName is the Elasticsearch index name
var IndexName string

// ElasticsearchUrls is the Elasticsearch cluster endpoints
var ElasticsearchUrls []string

var tWidth int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "idxtests",
	Short: "idxtest allows you to index Test Results",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	currentTime := time.Now()
	indexName := fmt.Sprintf("test-results-%s", currentTime.Format(layoutISO))
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.idxtests.yaml)")
	rootCmd.PersistentFlags().StringVarP(&IndexName, "index", "i", indexName, "elasticsearch index name")
	rootCmd.PersistentFlags().StringArrayVarP(&ElasticsearchUrls, "esUrls", "e", []string{"http://localhost:9200"}, "elasticsearch cluster endpoints")
	tWidth, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".idxtests" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".idxtests")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
