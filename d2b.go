package main

import (
	"fmt"
	"log"
	"code.google.com/p/d2b/bookmarks"
	"code.google.com/p/d2b/delicious"
	"code.google.com/p/d2b/blog"
	"code.google.com/p/d2b/data"
	"net/http"
	"html/template"
	"appengine"
	"strings"
	"time"
)

func init() {
	http.HandleFunc("/", create_sample_page)
	http.HandleFunc("/admin/", fake_admin)
	http.HandleFunc("/admin/config", handle_admin_config)
	http.HandleFunc("/cron/fetch", handle_cron_fetch)
	http.HandleFunc("/cron/publish", handle_cron_publish)
}

func fake_admin (w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Admin")
}

func handle_admin_config (w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r);

	if r.FormValue("submit") == "Set key" {
		data.SetConfigOption(c, r.FormValue("key"), r.FormValue("value"));
	}
	if r.FormValue("delete") != "" {
		data.DeleteConfigOption(c, r.FormValue("delete"))
		http.Redirect(w, r, "/admin/config", http.StatusFound)
	}

	options := data.GetAllConfigOptions(c);
	templ := template.Must(template.ParseFiles("templates/admin_config.templ"));
	templ.Execute(w, options);
}

func create_sample_page (w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r);

	// get from config
	headings := strings.Split(data.GetConfigOption(c, "headings"), ",")

	log.Printf("testing\n");

	bmarks := data.GetAllBookmarks(c);

	var bookmark_map = build_map(headings, "Everything else", bmarks);

	var html = blog.Format_html(headings, bookmark_map, false);
	fmt.Fprint(w, html)
}

func handle_cron_fetch (w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r);

	url := data.GetConfigOption(c, "bookmarkUrl");
	bookmark_service := data.GetConfigOption(c, "bookmarkService");

	var bmarks []bookmarks.Bookmark;
	if bookmark_service == "delicious" {
		bmarks = delicious.Parse(c, url);
	}
	data.SetConfigOption(c, "lastFetch", time.Now().String())
	data.AddBookmarks(c, bmarks);
}

func handle_cron_publish (w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r);

	headings := strings.Split(data.GetConfigOption(c, "headints"), ",")

	// get last publish date
	var last_published time.Time;

	bmarks := data.GetBookmarks(c, last_published);

	bookmark_map := build_map(headings, "Everything else", bmarks);

	// instead of printing, this will send somewhere based on the
	// configured blogging engine
	var html = blog.Format_html(headings, bookmark_map, false);
	fmt.Fprint(w, html)

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

