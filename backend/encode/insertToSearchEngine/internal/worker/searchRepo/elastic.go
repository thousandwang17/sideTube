/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-22 20:59:24
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-20 14:37:17
 * @FilePath: /generateMPD/internal/worker/metaRepo/mongo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package searchRepo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"toSearchEngine/internal/worker"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type elasticRepo struct {
	db *elasticsearch.Client
}

var (
	ErrUpsertVideo = errors.New("invalid upsert video document")
	ErrDB          = errors.New("DB error")
	ErrConnection  = errors.New("failed to connect network")
)

func NewElasticRepo(db *elasticsearch.Client) worker.SearchRepo {
	return &elasticRepo{
		db: db,
	}
}

// update video status to successed ,
func (e elasticRepo) InsertMultiVideoMeta(c context.Context, info []worker.VideoMeta) error {

	var bulkBody bytes.Buffer

	for i := range info {
		data, err := json.Marshal(info[i])
		if err != nil {
			log.Printf("Error marshaling document: %s", err)
			continue
		}
		updateReq := fmt.Sprintf(`{ "update" : {"_id" : "%s" , "_index" : "%s"}}`, info[i].Id, os.Getenv("ELASTIC_INDEX"))

		body := []byte(fmt.Sprintf(`{ "doc" : %s, "doc_as_upsert" : true }`, data))

		bulkBody.WriteString(updateReq)
		bulkBody.WriteByte('\n')
		bulkBody.Write(body)
		bulkBody.WriteByte('\n')
	}

	fmt.Print(bulkBody.String())
	// Set up the request object.
	req := esapi.BulkRequest{
		Body:    &bulkBody,
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(c, e.db)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return ErrConnection
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	if res.StatusCode == 500 {
		log.Println("Error db 500 ", res.Body)
		return ErrDB
	}

	log.Println(res.StatusCode, res.Body)
	return ErrUpsertVideo
}
