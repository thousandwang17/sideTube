/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-14 14:56:34
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-17 16:31:29
 * @FilePath: /insertToSearchEngine/internal/common/elastic/elastic.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

var elastcClient *elasticsearch.Client

func init() {
	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//

	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", os.Getenv("ELASTIC_HOST"), os.Getenv("ELASTIC_PORT")),
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	var r map[string]interface{}
	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	elastcClient = es

	indexName := os.Getenv("ELASTIC_INDEX")

	// Check if index exists
	existsRes, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		log.Fatal(err)
	}
	defer existsRes.Body.Close()

	if existsRes.StatusCode == 404 { // Index not found
		// Define mapping for index
		mapping := `
		{
			"mappings": {
				"dynamic": "false",
				"properties": {
					"title": {
						"type": "text"
					},
					"duration": {
						"type": "float"
					},
					"uploadTime": {
						"type": "date",
						"format": "yyyy-MM-dd HH:mm:ss"
					},
					"permission": {
						"type": "byte"
					}
				}
			}
		}`

		// Create index with mapping
		createRes, err := es.Indices.Create(
			indexName,
			es.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			log.Println(err)
		}
		defer createRes.Body.Close()

		// Check response status and do something with the response body
	}
}

func GetElasticClient() *elasticsearch.Client {
	return elastcClient
}
