package ooxx

import (
	"encoding/json"
	"errors"
	"strings"
)

type OOXXResult struct {
	Status      string      `json:"status"`
	CurrentPage int32       `json:"current_page"`
	PageCount   int32       `json:"page_count"`
	Count       int32       `json:"count"`
	Comments    []OOXXModel `json:"comments"`
}

type OOXXModel struct {
	CommentID      string   `json:"comment_ID"`
	CommentPostID  string   `json:"comment_post_ID"`
	CommentDate    string   `json:"comment_date"`
	CommentContent string   `json:"comment_content"`
	Pics           []string `json:"pics"`
	PicsStr        string   `json:"-"`
}

func (ox *OOXXModel) UnmarshalJSON(b []byte) error {
	if ox == nil {
		return errors.New("OOXXModel: UnmarshalJSON on nil pointer")
	}

	var stuff map[string]interface{}

	//return nil
	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}
	for k, v := range stuff {
		if k == "comment_ID" {
			ox.CommentID = v.(string)
		}
		if k == "comment_post_ID" {
			ox.CommentPostID = v.(string)
		}
		if k == "comment_date" {
			ox.CommentDate = v.(string)
		}
		if k == "comment_content" {
			ox.CommentContent = v.(string)
		}
		if k == "pics" {
			pics := v.([]interface{})
			var ps []string
			for _, pic := range pics {
				ps = append(ps, pic.(string))
			}
			ox.Pics = ps
			ox.PicsStr = strings.Join(ox.Pics, ";")
		}
	}
	return nil
}
