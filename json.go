package htmlDoc

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/ingmardrewing/fs"
)

type Json struct{}

func (j *Json) Read(value []byte, keys ...string) string {
	v, err := jsonparser.GetString(value, keys...)
	if err != nil {
		return ""
	}
	return v
}

func (j *Json) readPosts(value []byte, addPost func(id, title, description, content, imageUrl, thumbUrl, path, prodDomain, filename, createDate, disqusId string)) {
	postsDir, _ := jsonparser.GetString(value, "postsDir")
	fmt.Println("reading directory for json files ", postsDir)
	files := fs.ReadDirEntriesEndingWith(postsDir, "json")
	for _, f := range files {
		bytes := fs.ReadByteArrayFromFile(postsDir + "/" + f)
		id := j.Read(bytes, "post", "post_id")
		title := j.Read(bytes, "post", "title")
		thumbUrl := j.Read(bytes, "postThumb")
		description := j.Read(bytes, "post", "excerpt")
		disqusId := j.Read(bytes, "post", "custom_fields", "dsq_thread_id", "[0]")
		createDate := j.Read(bytes, "post", "date")
		content := j.Read(bytes, "post", "content")
		path := j.Read(bytes, "post", "url")
		filename := "index.html"

		imageUrl := "https://www.drewing.de/blog/wp-content/themes/drewing2012/silhouette_ingmar_drewing.png"
		if thumbUrl == "" {
			thumbUrl = "https://www.drewing.de/blog/wp-content/themes/drewing2012/silhouette_ingmar_drewing.png"
		}

		prodDomain := "https://drewing.de"

		addPost(id, title, description, content, imageUrl, thumbUrl, prodDomain, path, filename, createDate, disqusId)
	}
}

func (j *Json) readPage(fpath, filename string, addPage func(id, title, description, content, imageUrl, thumbUrl, path, prodDomain, filename, createDate, disqusId string)) {
	bytes := fs.ReadByteArrayFromFile(fpath)
	id := j.Read(bytes, "page", "post_id")
	title := j.Read(bytes, "page", "title")
	thumbUrl := j.Read(bytes, "postThumb")
	description := j.Read(bytes, "page", "excerpt")
	disqusId := j.Read(bytes, "page", "custom_fields", "dsq_thread_id", "[0]")
	createDate := j.Read(bytes, "page", "date")
	content := j.Read(bytes, "page", "content")
	path := j.Read(bytes, "page", "url")

	imageUrl := "https://www.drewing.de/blog/wp-content/themes/drewing2012/silhouette_ingmar_drewing.png"
	if thumbUrl == "" {
		thumbUrl = "https://www.drewing.de/blog/wp-content/themes/drewing2012/silhouette_ingmar_drewing.png"
	}

	prodDomain := "https://drewing.de"

	addPage(id, title, description, content, imageUrl, thumbUrl, prodDomain, path, filename, createDate, disqusId)
}
