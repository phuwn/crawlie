package keyword

import (
	"io"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/response"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/util"
)

// ListByUser - List user's keyword
func ListByUser(c echo.Context) error {
	var (
		uid    = util.GetUserIDFromCtx(c)
		tx     = util.GetTxFromCtx(c)
		limit  = 50
		offset = 0
		err    error
		search *string
	)

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return util.JsonError(c, 400, "limit should be number")
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return util.JsonError(c, 400, "offset should be number")
		}
	}

	if searchStr := c.QueryParam("q"); searchStr != "" {
		search = &searchStr
	}

	keywords, count, err := server.Get().Store().Keyword.ListByUser(
		tx,
		uid,
		limit,
		offset,
		search,
	)
	if err != nil {
		return err
	}
	return c.JSON(200, &response.ListKeywordsResponse{Data: keywords, Count: count})
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
		keywords     = make([]*model.Keyword, 0)
		userKeywords = make([]*model.UserKeyword, 0)
		keywordMap   = make(map[string]bool)
	)

	for _, line := range lines {
		keywordName := strings.TrimSpace(line)
		if keywordName == "" {
			continue
		}
		if keywordMap[keywordName] {
			continue
		}
		keywords = append(keywords, &model.Keyword{Name: keywordName, Status: model.KeywordNeedCrawl})
		userKeywords = append(userKeywords, &model.UserKeyword{UserID: uid, Keyword: keywordName, FileName: file.Filename})
		keywordMap[keywordName] = true
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
