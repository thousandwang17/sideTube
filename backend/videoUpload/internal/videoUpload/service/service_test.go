/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-10 14:50:43
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-02-15 15:22:14
 * @FilePath: /videoUpload/internal/videoUpload/service/service_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service_test

import (
	"context"
	"errors"
	"io"
	"sideTube/videoUpload/internal/videoUpload"
	"sideTube/videoUpload/internal/videoUpload/service"
	"testing"
)

type fileRepo struct {
	rep videoUpload.VideoRepoMeta
	err error
}

func newFileRepo(rep videoUpload.VideoRepoMeta, err error) *fileRepo {
	return &fileRepo{rep, err}
}

func (f fileRepo) CreateMultipartUpload(c context.Context, id string) (videoUpload.VideoRepoMeta, error) {
	return f.rep, f.err
}

func (f fileRepo) UploadPart(c context.Context, v videoUpload.VideoRepoMeta, data io.ReadSeeker) error {
	return f.err
}

func (f fileRepo) CompleteMultipartUpload(c context.Context, id string) error {
	return f.err
}

func (f fileRepo) AbortUpload(c context.Context, id string) error {
	return f.err
}

type metaRepo struct {
	id  string
	err error
}

func newMetaRepo(id string, err error) *metaRepo {
	return &metaRepo{id, err}
}

func (m metaRepo) Insert(c context.Context, userId string) (id string, err error) {
	return m.id, m.err
}

func (m metaRepo) UpdateState(c context.Context, videoId, userId string) error {
	return m.err
}

func (m metaRepo) UndoUpdateState(c context.Context, videoId, userId string) error {
	return m.err
}

func (m metaRepo) Remove(c context.Context, videoId, userId string) error {
	return m.err
}

func (m *metaRepo) setErr(e error) {
	m.err = e
}

type metaOption func(m *metaRepo)
type fileoption func(m *fileRepo)

func mSetErr(err string) metaOption {
	return func(m *metaRepo) {
		m.setErr(errors.New(err))
	}
}

func TestStartUpload(t *testing.T) {
	expectedId := "aD3tfsa2"
	mRepo := newMetaRepo(expectedId, nil)
	fRepo := newFileRepo(
		videoUpload.VideoRepoMeta{
			Id: expectedId,
		},
		nil,
	)

	s := service.NewVideoCommend(
		mRepo,
		fRepo,
		nil,
	)

	// assert := assert.New(t)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "userId", "Ds13dGaq")

	for _, v := range []struct {
		name  string
		want  string
		err   error
		mprev []metaOption
		fprev []fileoption
	}{
		{"success", expectedId, nil, nil, nil},
		{"meta crash", "", errors.New("connection lost"), []metaOption{mSetErr("connection lost")}, nil},
	} {
		for _, m := range v.mprev {
			m(mRepo)
		}

		for _, f := range v.fprev {
			f(fRepo)
		}

		id, err := s.StartUpload(ctx, 0)

		if v.want != id {
			t.Errorf("%v : want %s, get %s", v.name, v.want, id)
		}
		if v.err != nil && v.err.Error() != err.Error() {
			t.Errorf("%v : wantErr but got nil", v.name)
		}
	}
}
