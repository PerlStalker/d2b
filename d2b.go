package main

import (
	"fmt"
	"log"
	"code.google.com/p/d2b/bookmarks"
	"code.google.com/p/d2b/delicious"
	"code.google.com/p/d2b/blog"
)

func main() {
	var headings = []string{"Colorado", "politics", "technology"}

	log.Printf("testing\n");
	url := "http://feeds.delicious.com/v2/json/musicfreestatic";
	bmarks := delicious.Parse(url);
	//log.Printf("Bookmarks: %+v", bmarks);

	var bookmark_map = build_map(headings, "Everything else", bmarks);
	log.Printf("Bmark map: %+v", bookmark_map)

	var html = blog.Format_html(headings, bookmark_map, true);
	fmt.Println(html)
}

func build_map (
	headings []string,
	default_heading string,
	bmarks []bookmarks.Bookmark) (bmap map[string][]bookmarks.Bookmark) {

	bmap = make(map[string][]bookmarks.Bookmark)

	BMARK: for i:= 0; i < len(bmarks); i++ {
		for j:= 0; j < len(bmarks[i].Tags); j++ {
			// create an array for every tag but only if
			// it doesn't already exist
			if bmap[bmarks[i].Tags[j]] == nil {
				bmap[bmarks[i].Tags[j]] = make([]bookmarks.Bookmark, 1)
			}
		}
		// look at each tag and see if it matches one
		// of the headings. If it does, append the
		// bookmark to that heading's array.
		for k:= 0; k < len(headings); k++ {
			for j:= 0; j < len(bmarks[i].Tags); j++ {
				if bmarks[i].Tags[j] == headings[k] {
					bmap[headings[k]] = append(bmap[headings[k]], bmarks[i])
					continue BMARK
				}
			}
		}
	}

	return bmap
}
