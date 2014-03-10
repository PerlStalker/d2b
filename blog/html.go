package blog

import (
	"html/template"
	"log"
	"code.google.com/p/d2b/bookmarks"
	//"os"
	//"fmt"
	"bytes"
	"strings"
)

func Format_html (headings []string, bmap map[string][]bookmarks.Bookmark, body_only bool) (html string)  {
	const text_head = `<head><meta http-equiv="Content-Type" content="text/html; charset=utf-8" /></head><body>`;
	templ_head := template.Must(template.New("head").Parse(text_head));

	const text_footer = `</body>`;
	templ_footer := template.Must(template.New("footer").Parse(text_footer));


	const text_heading = `<h2>{{.}}</h2>`;
	templ_heading := template.Must(template.New("heading").Parse(text_heading));

	const text_bookmark_list = `<ul>
{{range . }}{{if .Url }}<li><a href="{{.Url}}">{{.Title}}</a>
{{if .Note}} - {{.Note}}{{end}}</li>{{end}}{{end}}
</ul>`;
	templ_bookmark_list := template.Must(template.New("bookmark_list").Parse(text_bookmark_list));

	//const text_bookmark = `<dt><a href="{{.Url}}">{{.Title}}</a></dt><dd>{{.Note}}</dd>`;
	//templ_bookmark := template.New("bookmark".Parse(text_bookmark);

	var buffer bytes.Buffer;

	if ! body_only {
		//templ_head.Execute(os.Stdout, nil);
		templ_head.Execute(&buffer, nil);
	}
	for i := 0; i < len(headings); i++ {
		// skip a heading if there are not bookmarks
		if len(bmap[headings[i]]) == 0 {
			continue
		}
		// print heading
		log.Printf("Doing heading: %s\n", headings[i]);

		// title case the first character of the heading
		first := headings[i][0:1];
		rest  := headings[i][1:len(headings[i])];

		//templ_heading.Execute(os.Stdout, headings[i]);
		templ_heading.Execute(&buffer, strings.ToTitle(first)+rest);
		// start bookmark wrapper
		//templ_bookmark_list.Execute(os.Stdout, bmap[headings[i]]);
		templ_bookmark_list.Execute(&buffer, bmap[headings[i]]);
	}
	if ! body_only {
		//templ_footer.Execute(os.Stdout, nil);
		templ_footer.Execute(&buffer, nil);
	}

	html = buffer.String()
	
	return html
}
