/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-05 16:19:49
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-07-16 19:32:18
 * @FilePath: /videoUpload/internal/transport/http/startUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	videoEnpoint "sideTube/videoUpload/internal/videoUpload/endpoint"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Define a mock struct that implements the VideoCommend service interface
type mockVideoService struct{}

func (m *mockVideoService) StartUpload() error {
	// Implement the StartUpload method as needed for testing
	return nil
}

// func TestStartUploadRegister(t *testing.T) {
// 	// Create a mock instance of the VideoCommend service
// 	videoService := &mockVideoService{}

// 	// Create a mock instance of the validator.Validate validator
// 	validator := validator.New()

// 	// Call the function and get the HttpTransport instance
// 	httpTransport := StartUploadRegister(videoService, validator)

// 	// Assert the type of the returned HttpTransport
// 	assert.IsType(t, &httptransport.HttpTransport{}, httpTransport)

// }

func TestDecodeStartUploadRequest(t *testing.T) {
	// Create a sample request body as a JSON string
	requestBody := `{"totalChunks":100}`

	// Create a new HTTP request with the sample request body
	req, err := http.NewRequest("POST", "/your-endpoint", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder to capture the response
	// httptest.NewRecorder()

	// Call the function and pass the request and response recorder
	result, err := decodeStartUploadRequest(req.Context(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the result is of type videoEnpoint.AbortUploadRequest
	strReq, ok := result.(videoEnpoint.StartUploadRequest)
	assert.True(t, ok)

	// Assert the expected values of the decoded request
	assert.Equal(t, int32(100), strReq.TotalChunks)
}

func TestEncodeStartUploadResponse(t *testing.T) {
	// Create a sample response object
	response := videoEnpoint.StartUploadRespond{
		VideoId: "123",
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the function and pass the response recorder
	err := encodeResponse(nil, rr, response)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the HTTP response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body as JSON
	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the expected values in the response body
	// assert.Equal(t, true, responseBody["success"])
	assert.Equal(t, "123", responseBody["videoId"])
}
