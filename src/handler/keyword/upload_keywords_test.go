package keyword

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/store"
	keywordMock "github.com/phuwn/crawlie/src/store/keyword/mocks"
	userKeywordMock "github.com/phuwn/crawlie/src/store/user_keyword/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type uploadFileTestCase struct {
	name     string
	file     []byte
	filename string
	wantErr  error
	code     int
	response string
}

func maximumNumberOfKeywordTestcase() *uploadFileTestCase {
	var fileData string
	for i := 0; i <= 100; i++ {
		fileData += fmt.Sprintf("%d\n", i)
	}
	return &uploadFileTestCase{
		name:     "exceed maximum number of keywords",
		filename: "test.csv",
		file:     []byte(fileData),
		wantErr:  nil,
		code:     400,
		response: `{"error":"maximum number of keywords is 100"}`,
	}
}

func TestUploadFile(t *testing.T) {
	testcases := []*uploadFileTestCase{
		{
			name:     "no file provided",
			wantErr:  nil,
			code:     400,
			response: `{"error":"file CSV ` + "`file`" + ` not found"}`,
		},
		{
			name:     "file provided not csv",
			file:     []byte{},
			filename: "test.txt",
			wantErr:  nil,
			code:     400,
			response: `{"error":"input file ` + "`file`" + ` should be in csv format"}`,
		},
		maximumNumberOfKeywordTestcase(),
		{
			name:     "happy case",
			file:     []byte("keyword 1\nkeyword 2\nkeyword 3"),
			filename: "test.csv",
			wantErr:  nil,
			code:     200,
			response: `{"message":"ok"}`,
		},
	}

	keywordStore := keywordMock.NewStore(t)
	keywordStore.Mock.On("BulkInsert", mock.Anything, mock.Anything).Return(nil)
	userKeywordStore := userKeywordMock.NewStore(t)
	userKeywordStore.Mock.On("BulkInsert", mock.Anything, mock.Anything).Return(nil)
	server.SetupTest(nil, nil, &store.Store{
		Keyword:     keywordStore,
		UserKeyword: userKeywordStore,
	}, nil)

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			body := new(bytes.Buffer)
			req := httptest.NewRequest(http.MethodPost, "/v1/keywords", body)
			if tt.file != nil {
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", tt.filename)
				part.Write(tt.file)
				writer.Close()
				req.Header.Set("Content-Type", writer.FormDataContentType())
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("tx", &gorm.DB{})
			err := UploadFile(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.code, rec.Code)
			assert.Equal(t, tt.response, rec.Body.String())
		})
	}
}
