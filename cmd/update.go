/*
Copyright © 2020 Isadora Ribeiro - github.com/isaitgirl

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	// To avoid break in go mod
	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/spf13/cobra"
)

type docs struct {
	Docs        []dataUnit
	PublishedIn string `json:"published_in"`
	UpdatedAt   string `json:"updated_at"`
}

type dataUnit struct {
	State    string
	CityName string `json:"city_name"`
	Date     string
	Cases    float64
	Count    float64
}

type influxDB struct {
	URL string
	Cli *client.Client
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update InfluxDB database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Update called")

		jsonBytes, err := ioutil.ReadFile("./resources/db.json")

		if err != nil {
			panic(err)
		}

		docs := docs{}

		_ = json.Unmarshal(jsonBytes, &docs)

		//TODO: InfluxDB authentication support
		influxHostURL := os.Getenv("INFLUXDB_URL")
		fmt.Println("InfluxDB endpoint URL:", influxHostURL)

		if influxHostURL == "" {
			fmt.Println("You must set INFLUXDB_URL environment variable. Ex: http://influxdb.local:8086/")
			os.Exit(1)
		}

		ic := &influxDB{
			URL: influxHostURL,
		}

		cli, err := client.NewHTTPClient(client.HTTPConfig{
			Addr: ic.URL,
		})

		if err != nil {
			fmt.Println("Error creating InfluxDB Client: ", err.Error())
		}

		points := make([]*client.Point, 0)

		for _, data := range docs.Docs {
			t := time.Now()
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			if strings.HasPrefix(data.Date, "2020") {
				data.Date = data.Date + "T00:00:00Z"
				t, _ = time.Parse(time.RFC3339, data.Date)
			}

			tags := make(map[string]string, 0)
			tags["state"] = data.State
			tags["city"] = strings.Replace(data.CityName, " ", "_", -1)

			fields := make(map[string]interface{}, 0)
			fields["cases"] = data.Cases
			fields["count"] = data.Count

			point, _ := client.NewPoint("data", tags, fields, t)

			points = append(points, point)
		}

		spew.Dump(points)

		bpc := client.BatchPointsConfig{
			Database:  "covid19",
			Precision: "h",
		}
		bp, _ := client.NewBatchPoints(bpc)
		bp.AddPoints(points)

		cli.Write(bp)

		defer cli.Close()

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
