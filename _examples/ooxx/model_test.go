package ooxx

import (
	"testing"
	"encoding/json"
)

func TestOOXXModel_UnmarshalJSON(t *testing.T) {
	stuff :=`
	{
      "comment_ID": "3936846",
      "comment_post_ID": "26402",
      "comment_author": "鱼腩",
      "comment_date": "2018-08-22 08:17:15",
      "comment_date_gmt": "2018-08-22 00:17:15",
      "comment_content": "智能家电够智障的…\nOk Google \n打开电视\nOk Google\n向上换台\nOk Google\n向上换台\nOk Google\n向上换台\nOk Google\n向上换台\nOk Google\n电视遥控器在哪儿？\n<img src=\"http://ww3.sinaimg.cn/mw600/006XNEY7gy1fui5yo95u0j31kw1kw1kx.jpg\" />\n<img src=\"http://ww3.sinaimg.cn/mw600/0073ob6Pgy1fui5yjut79j31kw0vvh0h.jpg\" />",
      "user_id": "0",
      "vote_positive": "23",
      "vote_negative": "16",
      "sub_comment_count": "10",
      "text_content": "智能家电够智障的…\nOk Google \n打开电视\nOk Google\n向上换台\nOk Google\n向上换台\nOk Google\n向上换台\nOk Google\n向上换台\nOk Google\n电视遥控器在哪儿？\n\n",
      "pics": [
        "http://ww3.sinaimg.cn/mw600/006XNEY7gy1fui5yo95u0j31kw1kw1kx.jpg",
        "http://ww3.sinaimg.cn/mw600/0073ob6Pgy1fui5yjut79j31kw0vvh0h.jpg"
      ]
    }`

	var ox OOXXModel
	 if err := json.Unmarshal([]byte(stuff), &ox); err != nil {
	 	t.Error("OOXXModel Unmarshal failed")
	 }
	if len(ox.Pics) != 2 {
		t.Error("OOXXModel pics Unmarshal failed")
	}
}
