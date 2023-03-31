/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2022-12-29 17:06:28
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-19 17:07:13
 * @FilePath: /search/internal/search/repository.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package search

import (
	"context"
)

type SearchRepository interface {
	SearchVideos(c context.Context, query string, size int64, from int64) (ids []string, err error)
}

type MetaRepository interface {
	GetPublicVideosByIds(ctx context.Context, ids []string) ([]VideoMeta, error)
}
