package api_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mdanialr/api-mock/api"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var sample = `
endpoint:
  /api/v1:
  /api/v2:
    status: 422
    200:
     sample: '{"status":"ok"}'
  /api/v3:
    200:
     sample: '{"status":"v3"}'
  /api/v4:
    status: 400
    400:
      sample: '{"status":"error"}'
  /api/v5:
    200:
      data: 2
      sample:
        - '{"status":"ok"}'
        - '{"status":"success"}'
  /api/v6:
    200:
      data: 5
      sample:
        - '{"status":"ok"}'
        - '{"status":"2nd success"}'
  /api/v7:
    status: 400
    400:
      sample:
        - '{"status":"error"}'
        - '{"status":"crash"}'
`

func TestMainHandler_Entry(t *testing.T) {
	conf := viper.New()
	conf.SetConfigType("yaml")
	conf.ReadConfig(bytes.NewBufferString(sample))

	testCases := []struct {
		name       string
		endpoint   string
		expectCode int
		expectResp string
	}{
		{
			name:       "Given non existent endpoint should return error 404 and expected json response",
			endpoint:   "/api",
			expectCode: http.StatusNotFound,
			expectResp: `{"message":"oops!! there is nothing here ¯\\_(ツ)_/¯"}`,
		},
		{
			name: "Given endpoint that exist but empty should return error 404 and" +
				" expected json response",
			endpoint:   "/api/v1",
			expectCode: http.StatusNotFound,
			expectResp: `{"message":"oops!! there is nothing here ¯\\_(ツ)_/¯"}`,
		},
		{
			name: "Given endpoint that exist also has status 422 but has no sample for that status" +
				" should return error 400 and expected json response",
			endpoint:   "/api/v2",
			expectCode: http.StatusBadRequest,
			expectResp: `{"message":"there is no sample in the config, make sure you add it and follow the convention"}`,
		},
		{
			name: "Given endpoint that exist and has no status but has sample for 200 should" +
				" return 200 and that sample as response",
			endpoint:   "/api/v3",
			expectCode: http.StatusOK,
			expectResp: `{"status":"v3"}`,
		},
		{
			name: "Given endpoint that exist has status 400 along with the sample should return" +
				" 400 and that sample as response",
			endpoint:   "/api/v4",
			expectCode: http.StatusBadRequest,
			expectResp: `{"status":"error"}`,
		},
		{
			name: "Given 200 as chosen sample has data 2 and more than one samples should return 200" +
				" and the second sample as response",
			endpoint:   "/api/v5",
			expectCode: http.StatusOK,
			expectResp: `{"status":"success"}`,
		},
		{
			name: "Given 200 as chosen sample has data 4 but only has 2 samples should return 200" +
				" and the last sample as response",
			endpoint:   "/api/v6",
			expectCode: http.StatusOK,
			expectResp: `{"status":"2nd success"}`,
		},
		{
			name: "Given 400 as chosen sample has no data but has 2 samples should return 400 and the" +
				" first sample as response",
			endpoint:   "/api/v7",
			expectCode: http.StatusBadRequest,
			expectResp: `{"status":"error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, tc.endpoint, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			h := api.NewMainHandler(conf)

			if assert.NoError(t, h.Entry(ctx)) {
				assert.Equal(t, tc.expectCode, rec.Code)
				assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))

				// in case you wondering why we need to trim the string see https://go.dev/src/encoding/json/stream.go
				// at section -- func (enc *Encoder) Encode(v any) error --
				// they add `\n` as delimiter for good, so we need to trim it before asserting.
				assert.Equal(t, tc.expectResp, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}
