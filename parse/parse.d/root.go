// Copyright © 2017 Slotix s.r.o. <dm@slotix.sk>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"os"

	"github.com/slotix/dataflowkit/healthcheck"
	"github.com/slotix/dataflowkit/parse"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//VERSION               string // VERSION is set during build
	DFKParse string
	DFKFetch string

	storageType    string
	storageExpires int64 //how long in seconds object stay in a cache before expiration.
	diskvBaseDir   string
	fetchBucket    string

	redisHost       string
	redisExpire     int
	redisNetwork    string
	redisPassword   string
	redisDB         int
	redisSocketPath string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dataflowkit",
	Short: "DataFlow Kit html parser",
	Long: `DataFlow Kit html parser serves for scraping data from websites according to chosen css selectors.
	Here is an example of payload structure:
	
	{"format":"json",
		"collections": [
				{
				"name": "collection1",
				"url": "http://example1.com",
				"fields": [
					{
						"field_name": "link",
						"css_selector": ".link a"
					},
					{
						"field_name": "Text",
						"css_selector": ".text"
					},
					{
						"field_name": "Image",
						"css_selector": ".foto img"
					}
				]
			}
		]
	}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking services ... ")

		services := []healthcheck.Checker{
			healthcheck.FetchConn{
				Host: DFKFetch,
			},
		}
		if storageType == "Redis" {
			services = append(services, healthcheck.RedisConn{
				Network: redisNetwork,
				Host:    redisHost})
		}
		status := healthcheck.CheckServices(services...)
		allAlive := true

		for k, v := range status {
			fmt.Printf("%s: %s\n", k, v)
			if v != "Ok" {
				allAlive = false
			}
		}
		if allAlive {
			fmt.Printf("Starting Server %s\n", DFKParse)
			parse.Start(DFKParse)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VERSION = version

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	//flags and configuration settings. They are global for the application.

	RootCmd.Flags().StringVarP(&DFKParse, "DFK_PARSE", "p", "127.0.0.1:8001", "HTTP listen address")
	RootCmd.Flags().StringVarP(&DFKFetch, "DFK_FETCH", "f", "127.0.0.1:8000", "DFK Fetch service address")

	//set here default type of storage
	RootCmd.Flags().StringVarP(&storageType, "STORAGE_TYPE", "", "Diskv", "Storage backend for intermediary data passed to html parser. Types: S3, Redis, Diskv")
	RootCmd.Flags().Int64VarP(&storageExpires, "STORAGE_EXPIRE", "", 3600, "Default Storage expire value in seconds")
	RootCmd.Flags().StringVarP(&diskvBaseDir, "DISKV_BASE_DIR", "", "diskv", "diskv base directory for storing fetch results")
	RootCmd.Flags().StringVarP(&fetchBucket, "FETCH_BUCKET", "", "fetch-bucket", "S3 bucket name for storing fetch results")

	RootCmd.Flags().StringVarP(&redisHost, "REDIS", "r", "127.0.0.1:6379", "Redis host address")
	RootCmd.Flags().IntVarP(&redisExpire, "REDIS_EXPIRE", "", 3600, "Default Redis expire value in seconds")
	RootCmd.Flags().StringVarP(&redisNetwork, "REDIS_NETWORK", "", "tcp", "Redis Network")
	RootCmd.Flags().StringVarP(&redisPassword, "REDIS_PASSWORD", "", "", "Redis Password")
	RootCmd.Flags().IntVarP(&redisDB, "REDIS_DB", "", 0, "Redis DB")
	RootCmd.Flags().StringVarP(&redisSocketPath, "REDIS_SOCKET_PATH", "", "", "Redis Socket Path")

	viper.AutomaticEnv() // read in environment variables that match
	viper.BindPFlag("DFK_FETCH", RootCmd.Flags().Lookup("DFK_FETCH"))
	viper.BindPFlag("DFK_PARSE", RootCmd.Flags().Lookup("DFK_PARSE"))

	viper.BindPFlag("STORAGE_TYPE", RootCmd.Flags().Lookup("STORAGE_TYPE"))
	viper.BindPFlag("STORAGE_EXPIRE", RootCmd.Flags().Lookup("STORAGE_EXPIRE"))
	viper.BindPFlag("DISKV_BASE_DIR", RootCmd.Flags().Lookup("DISKV_BASE_DIR"))
	viper.BindPFlag("FETCH_BUCKET", RootCmd.Flags().Lookup("FETCH_BUCKET"))
	viper.BindPFlag("REDIS", RootCmd.Flags().Lookup("REDIS"))
	viper.BindPFlag("REDIS_EXPIRE", RootCmd.Flags().Lookup("REDIS_EXPIRE"))
	viper.BindPFlag("REDIS_NETWORK", RootCmd.Flags().Lookup("REDIS_NETWORK"))
	viper.BindPFlag("REDIS_PASSWORD", RootCmd.Flags().Lookup("REDIS_PASSWORD"))
	viper.BindPFlag("REDIS_DB", RootCmd.Flags().Lookup("REDIS_DB"))
	viper.BindPFlag("REDIS_SOCKET_PATH", RootCmd.Flags().Lookup("REDIS_SOCKET_PATH"))
}