package data

import (
	"appengine"
	"appengine/datastore"
	"log"
	"time"
	"code.google.com/p/d2b/bookmarks"
)

func bookmarkKey (c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Bookmark", "bookmarks", 0, nil);
}

func AddBookmark (c appengine.Context, bmark bookmarks.Bookmark) {
	q := datastore.NewQuery("Bookmark").
		Filter("Url = ", bmark.Url);
	results := q.Run(c);

	db_entry := datastore.NewIncompleteKey(c, "Bookmark", bookmarkKey(c));

	var found_bmark bookmarks.Bookmark
	for {
		db_key, err := results.Next(&found_bmark);
		if err == datastore.Done {
			break // the result was not found
		}
		if err != nil {
			log.Fatal("error fetching bookmarks: ", err)
			break
		}
		db_entry = db_key
	}

	_, err := datastore.Put(c, db_entry, &bmark);
	if err != nil {
		log.Fatal("error adding bookmark: ", err)
	}
}

func AddBookmarks (c appengine.Context, bmarks []bookmarks.Bookmark) {
	for i := 0; i < len(bmarks); i++ {
		AddBookmark(c, bmarks[i]);
	}
}

// add time stamp to call?
func DeleteOldBookmarks (c appengine.Context) {
}

func GetBookmarks (c appengine.Context, from_time time.Time) (bmarks []bookmarks.Bookmark) {
	q := datastore.NewQuery("Bookmark").
		Filter("Time > ", from_time).
		Order("Time");

	_, err := q.GetAll(c, &bmarks);

	if err != nil {
		log.Fatal("error fetching options: ", err);
	}

	return bmarks
}

func GetAllBookmarks (c appengine.Context) (bmarks []bookmarks.Bookmark) {
	q := datastore.NewQuery("Bookmark").
		Order("Time");

	_, err := q.GetAll(c, &bmarks);

	if err != nil {
		log.Fatal("error fetching options: ", err);
	}

	return bmarks
}
