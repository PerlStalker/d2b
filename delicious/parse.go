package delicious;

import (
	"encoding/json"
	"log"
	"io/ioutil"
	//"net/http"
	"time"
	"code.google.com/p/d2b/bookmarks"
	//"fmt"
	"appengine"
	"appengine/urlfetch"
)


type DeliciousBookmark struct {
	a string
	d string
	n string
	u string
	t []string
	dt string
}

func Parse (c appengine.Context, url string) (marks []bookmarks.Bookmark) {
	//resp, err := http.Get(url); // net/http
	client := urlfetch.Client(c);
	resp, err := client.Get(url);

	if err != nil {
		log.Fatal(err);
	}

	body, err := ioutil.ReadAll(resp.Body);
	defer resp.Body.Close();

	if err != nil {
		log.Fatal("error fetching bookmarks:", err)
	}

	//log.Printf("Body: %s\n", body);

	var json_data interface{};
	err = json.Unmarshal(body, &json_data);
	if err != nil {
		log.Fatal("json parse error: ", err);
	}
	//log.Printf("JSON: %#v\n", json_data);

	json_map := json_data.([]interface{})

	marks = make([]bookmarks.Bookmark, len(json_map));
	for i := 0; i < len(json_map); i++ {
		bookmark_map := json_map[i].(map[string]interface{});
		//log.Printf("BM: %+v\n", bookmark_map);

		var bm bookmarks.Bookmark;
		bm.Title = bookmark_map["d"].(string);
		bm.Time, err = time.Parse(time.RFC3339, bookmark_map["dt"].(string));
		bm.Url = bookmark_map["u"].(string);
		bm.Note = bookmark_map["n"].(string);

		bm.Author = bookmark_map["a"].(string);

		tags_map := bookmark_map["t"].([]interface{});
		bm.Tags = make([]string, len(tags_map));
		for i, u := range tags_map {
			bm.Tags[i] = u.(string);
		}

		//log.Printf("BM struct: %+v\n", bm);
		marks[i] = bm;
	}

	return marks
}
