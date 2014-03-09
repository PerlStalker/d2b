package bookmarks

import (
	"time"
)

type Bookmark struct {
	Tags  []string;
	Title string;
	Time  time.Time; // delicious parses as RFC3339
	Url   string;
	Note  string;
	Author string;

	Approved bool;
	ApprovedBy string;
	ModeratorsNote string;
}
