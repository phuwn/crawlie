package keyword

import (
	"io"
	"path"
	"strings"

	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/util"
)

// ListByUser - List user's keyword
func ListByUser(c echo.Context) error {
	var (
		uid = util.GetUserIDFromCtx(c)
		tx  = util.GetTxFromCtx(c)
	)

	u, err := server.Get().Store().Keyword.ListByUser(tx, uid)
	if err != nil {
		return err
	}
	return c.JSON(200, u)
}

// Get - Get keyword detail by name
func Get(c echo.Context) error {
	var (
		name = c.Param("name")
		tx   = util.GetTxFromCtx(c)
	)

	u, err := server.Get().Store().Keyword.Get(tx, name)
	if err != nil {
		return util.JsonError(c, 404, "keyword not found")
	}
	return c.JSON(200, u)
}

// UploadFile - Upload csv file contains list of keywords
func UploadFile(c echo.Context) error {
	var (
		uid = util.GetUserIDFromCtx(c)
		tx  = util.GetTxFromCtx(c)
	)

	file, err := c.FormFile("file")
	if err != nil {
		return util.JsonError(c, 400, "file CSV `file` not found")
	}

	if strings.ToLower(path.Ext(file.Filename)) != ".csv" {
		return util.JsonError(c, 400, "input file `file` should be in csv format")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	b, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
	if len(lines) > 100 {
		return util.JsonError(c, 400, "maximum number of keywords is 100")
	}

	var (
		keywords     = make([]*model.Keyword, len(lines))
		userKeywords = make([]*model.UserKeyword, len(lines))
	)

	for i, line := range lines {
		keywordName := strings.TrimSpace(line)
		keywords[i] = &model.Keyword{Name: keywordName, Status: model.KeywordNeedCrawl}
		userKeywords[i] = &model.UserKeyword{UserID: uid, Keyword: keywordName, FileName: file.Filename}
	}

	srv := server.Get()
	err = srv.Store().Keyword.BulkInsert(tx, keywords)
	if err != nil {
		return err
	}

	err = srv.Store().UserKeyword.BulkInsert(tx, userKeywords)
	if err != nil {
		return err
	}

	return c.JSONBlob(200, []byte(`{"message":"ok"}`))
}
