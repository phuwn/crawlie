package response

import "github.com/phuwn/crawlie/src/model"

type ListKeywordsResponse struct {
	Data  []*model.Keyword `json:"data"`
	Count int64            `json:"count"`
}
