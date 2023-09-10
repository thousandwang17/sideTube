/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-07-06 15:05:21
 * @FilePath: /search/internal/search/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"context"
	"errors"
	"sideTube/search/internal/search"
)

var (
	ErrVaild = errors.New("args vaild failed")
)

type service struct {
	searchRepo search.SearchRepository
	metaRepo   search.MetaRepository
}

type SearchCommend interface {
	Serach(c context.Context, query string, skip int64, length int64) (data []search.VideoMeta, err error)
}

func NewsearchCommend(searchdb search.SearchRepository, metaDB search.MetaRepository) SearchCommend {
	return &service{
		searchRepo: searchdb,
		metaRepo:   metaDB,
	}
}

func (v service) Serach(c context.Context, query string, skip int64, length int64) (data []search.VideoMeta, err error) {

	ids, err := v.searchRepo.SearchVideos(c, query, length, skip)

	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []search.VideoMeta{}, nil
	}

	data, err = v.metaRepo.GetPublicVideosByIds(c, ids)

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []search.VideoMeta{}, nil
	}

	return data, nil
}
