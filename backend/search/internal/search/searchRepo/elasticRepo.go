/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-31 16:07:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-20 15:39:55
 * @FilePath: /search/internal/search/metaRepository/ElasticRepo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package searchRepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sideTube/search/internal/search"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	ErrSearchVideo = errors.New("search enginde error")
)

type ElasticlRepo struct {
	db *elasticsearch.Client
}

// NOTE : this search list just get videos that public and  been encoded order by time
// in production should use Collaborative filtering algorithms to build search list ,
// Apache Mahout, Apache Spark and  TensorFlow are good choose to impelemnet searchation system
func NewElasticRepo(db *elasticsearch.Client) search.SearchRepository {
	return &ElasticlRepo{
		db: db,
	}
}

func (m ElasticlRepo) SearchVideos(c context.Context, query string, size int64, from int64) ([]string, error) {

	searchBody := `{
		"_source": ["video_id"],
		"query": {
			"multi_match": {
				"query": "%s",
				"type": "best_fields",
				"fields": ["title"]
			}
		},
		"size": %d,
		"from": %d
	}`

	searchBody = fmt.Sprintf(searchBody, query, size, from)

	req := esapi.SearchRequest{
		Index: []string{"videos"},
		Body:  strings.NewReader(searchBody),
	}

	// Perform the search request with the client.
	res, err := req.Do(c, m.db)
	if err != nil {
		log.Printf("Error performing search: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Search request failed with status code %d  ", res.StatusCode)
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading search response body: %s", err)
			return nil, err
		}
		log.Printf(" reading search response body: %s", string(bodyBytes))
		return nil, ErrSearchVideo
	}

	var response struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Printf("Error decode response body: %s", err)
		return nil, nil
	}

	// Parse the search results into a slice of VideoMeta objects.
	var videos_id []string

	for _, hit := range response.Hits.Hits {

		var source struct {
			VideoID string `json:"video_id"`
		}

		if err := json.Unmarshal(hit.Source, &source); err != nil {
			return nil, err
		}

		videos_id = append(videos_id, source.VideoID)
	}

	return videos_id, nil
}
