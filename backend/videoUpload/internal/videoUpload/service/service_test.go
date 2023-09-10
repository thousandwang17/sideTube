// /*
//  * @Author: dennyWang thousandwang17@gmail.com
//  * @Date: 2023-01-10 14:50:43
//  * @LastEditors: dennyWang thousandwang17@gmail.com
//  * @LastEditTime: 2023-05-14 14:05:17
//  * @FilePath: /videoUpload/internal/videoUpload/service/service_test.go
//  * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
//  */
// package service_test

// import (
// 	"context"
// 	"errors"
// 	"io"
// 	"sideTube/videoUpload/internal/videoUpload"
// 	"sideTube/videoUpload/internal/videoUpload/service"
// 	"testing"
// )

// type fileRepo struct {
// 	rep videoUpload.VideoRepoMeta
// 	err error
// }

// func newFileRepo(rep videoUpload.VideoRepoMeta, err error) *fileRepo {
// 	return &fileRepo{rep, err}
// }

// func (f fileRepo) CreateMultipartUpload(c context.Context, id string) (videoUpload.VideoRepoMeta, error) {
// 	return f.rep, f.err
// }

// func (f fileRepo) UploadPart(c context.Context, v videoUpload.VideoRepoMeta, data io.ReadSeeker) error {
// 	return f.err
// }

// func (f fileRepo) CompleteMultipartUpload(c context.Context, id string) error {
// 	return f.err
// }

// func (f fileRepo) AbortUpload(c context.Context, id string) error {
// 	return f.err
// }

// type metaRepo struct {
// 	id  string
// 	err error
// }

// func newMetaRepo(id string, err error) *metaRepo {
// 	return &metaRepo{id, err}
// }

// func (m metaRepo) Insert(c context.Context, userId, userName string) (id string, err error) {
// 	return m.id, m.err
// }

// func (m metaRepo) UpdateState(c context.Context, videoId, userId string) error {
// 	return m.err
// }

// func (m metaRepo) UndoUpdateState(c context.Context, videoId, userId string) error {
// 	return m.err
// }

// func (m metaRepo) Remove(c context.Context, videoId, userId string) error {
// 	return m.err
// }

// func (m *metaRepo) setErr(e error) {
// 	m.err = e
// }

// type metaOption func(m *metaRepo)
// type fileoption func(m *fileRepo)

// func mSetErr(err string) metaOption {
// 	return func(m *metaRepo) {
// 		m.setErr(errors.New(err))
// 	}
// }

// func TestStartUpload(t *testing.T) {
// 	expectedId := "aD3tfsa2"

// 	mRepo := newMetaRepo(expectedId, nil)
// 	fRepo := newFileRepo(
// 		videoUpload.VideoRepoMeta{
// 			Id: expectedId,
// 		},
// 		nil,
// 	)

// 	s := service.NewVideoCommend(
// 		mRepo,
// 		fRepo,
// 		nil,
// 	)

// 	// assert := assert.New(t)
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "uid", "Ds13dGaq")
// 	ctx = context.WithValue(ctx, "userName", "test")

// 	for _, v := range []struct {
// 		name  string
// 		want  string
// 		err   error
// 		mprev []metaOption
// 		fprev []fileoption
// 	}{
// 		{"success", expectedId, nil, nil, nil},
// 		{"meta crash", "", errors.New("connection lost"), []metaOption{mSetErr("connection lost")}, nil},
// 	} {
// 		for _, m := range v.mprev {
// 			m(mRepo)
// 		}

// 		for _, f := range v.fprev {
// 			f(fRepo)
// 		}

// 		id, err := s.StartUpload(ctx, 0)

// 		if v.want != id {
// 			t.Errorf("%v : want %s, get %s", v.name, v.want, id)
// 		}
// 		if v.err != nil && v.err.Error() != err.Error() {
// 			t.Errorf("%v : wantErr but got nil", v.name)
// 		}
// 	}
// }

package service

import (
	"context"
	"io"
	"sideTube/videoUpload/internal/videoUpload"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMetaRepository implements the MetaRepository interface for testing purposes
type MockMetaRepository struct {
	mock.Mock
}

func (m *MockMetaRepository) Insert(ctx context.Context, userId, name string) (string, error) {
	args := m.Called(ctx, userId, name)
	return args.String(0), args.Error(1)
}

func (m *MockMetaRepository) UpdateState(ctx context.Context, videoId, userId string) error {
	args := m.Called(ctx, videoId, userId)
	return args.Error(0)
}

func (m *MockMetaRepository) UndoUpdateState(ctx context.Context, videoId, userId string) error {
	args := m.Called(ctx, videoId, userId)
	return args.Error(0)
}

func (m *MockMetaRepository) Remove(ctx context.Context, videoId, userId string) error {
	args := m.Called(ctx, videoId, userId)
	return args.Error(0)
}

// MockVideoRepository implements the VideoRepository interface for testing purposes
type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) CreateMultipartUpload(ctx context.Context, id string) (videoUpload.VideoRepoMeta, error) {
	args := m.Called(ctx, id)

	// Assuming your VideoRepoMeta struct has the fields Id, UploadId, UserId, PartID, and State
	return videoUpload.VideoRepoMeta{
		Id:       args.String(0),
		UploadId: args.String(1),
		UserId:   args.String(2),
		PartID:   int64(args.Int(3)),
		State:    int8(args.Int(4)),
	}, args.Error(5)
}

func (m *MockVideoRepository) UploadPart(ctx context.Context, v videoUpload.VideoRepoMeta, file io.ReadSeeker) error {
	args := m.Called(ctx, v, file)
	return args.Error(1)
}

func (m *MockVideoRepository) CompleteMultipartUpload(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(1)
}

func (m *MockVideoRepository) AbortUpload(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(1)
}

// MockMessageQueue implements the MessagQqueue interface for testing purposes
type MockMessageQueue struct {
	mock.Mock
}

func (m *MockMessageQueue) NotifyMergeVideo(ctx context.Context, videoId, userId string) error {
	args := m.Called(ctx, videoId, userId)
	return args.Error(0)
}

func TestStartUpload(t *testing.T) {
	// Create mock instances for repositories and queue
	mockMetaRepo := &MockMetaRepository{}
	mockVideoRepo := &MockVideoRepository{}
	mockQueueRepo := &MockMessageQueue{}

	// Create the service using the mock instances
	service := NewVideoCommend(mockMetaRepo, mockVideoRepo, mockQueueRepo)

	// Create a context with mocked values
	ctx := context.WithValue(context.Background(), "uid", "mockedUserID")
	ctx = context.WithValue(ctx, "userName", "mockedUserName")

	// Set up mock behavior for MetaRepository Insert
	mockMetaRepo.On("Insert", ctx, "mockedUserID", "mockedUserName").Return("mockedVideoID", nil)

	// Set up mock behavior for VideoRepository CreateMultipartUpload
	mockVideoRepo.On("CreateMultipartUpload", ctx, "mockedVideoID").Return(videoUpload.VideoRepoMeta{
		Id:       "mockedVideoID",
		UploadId: "mockedUploadID",
		UserId:   "mockedUserID",
		PartID:   123, // Replace with the desired value for PartID
		State:    0,   // Replace with the desired value for State
	}, nil)

	// Call the StartUpload function
	videoID, err := service.StartUpload(ctx, 5)

	// Assert that the function returned the expected values
	assert.NoError(t, err, "StartUpload should not return an error")
	assert.Equal(t, "mockedVideoID", videoID, "StartUpload should return the expected videoID")

	// Assert that the mock methods were called with the correct arguments
	mockMetaRepo.AssertExpectations(t)
	mockVideoRepo.AssertExpectations(t)
	mockQueueRepo.AssertExpectations(t)
}
