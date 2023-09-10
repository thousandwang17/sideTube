/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-06 15:28:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-07-06 21:21:50
 * @FilePath: /videoUpload/internal/videoUpload/endpoint/abortUpload.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	httptransport "sideTube/videoUpload/internal/common/simpleKit/httpTransport"
	videoEnpoint "sideTube/videoUpload/internal/videoUpload/endpoint"
	"sideTube/videoUpload/internal/videoUpload/service"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

func TestAbortUploadRegister(t *testing.T) {
	type args struct {
		svc service.VideoCommend
		v   *validator.Validate
	}
	tests := []struct {
		name string
		args args
		want *httptransport.HttpTransport
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbortUploadRegister(tt.args.svc, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AbortUploadRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeAbortUploadRequest(t *testing.T) {
	// Create a sample request body as a JSON string
	requestBody := `{"video_id":"123"}`

	// Create a new HTTP request with the sample request body
	req, err := http.NewRequest("POST", "/your-endpoint", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder to capture the response
	// httptest.NewRecorder()

	// Call the function and pass the request and response recorder
	result, err := decodeAbortUploadRequest(req.Context(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the result is of type videoEnpoint.AbortUploadRequest
	abortReq, ok := result.(videoEnpoint.AbortUploadRequest)
	assert.True(t, ok)

	// Assert the expected values of the decoded request
	assert.Equal(t, "123", abortReq.VideoId)
}

func TestEncodeAbortUploadResponse(t *testing.T) {
	// Create a sample response object
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "Upload aborted successfully",
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the function and pass the response recorder
	err := encodeAbortUploadResponse(nil, rr, response)
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
	assert.Equal(t, true, responseBody["success"])
	assert.Equal(t, "Upload aborted successfully", responseBody["message"])
}
