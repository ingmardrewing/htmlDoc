package htmlDoc

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ingmardrewing/aws"
	"github.com/ingmardrewing/img"

	"gopkg.in/russross/blackfriday.v2"
)

const (
	jsonTemplate = `{
	"thumbImg":"%s",
	"postImg":"%s",
	"filename":"%s",
	"post":{
		"post_id":"000",
		"date":"%s",
		"url":"%s",
		"title":"%s",
		"excerpt":"%s",
		"content":"%s",
		"custom_fields":{
			"dsq_thread_id":["%s"]
		}
	}
}
`
)

func NewPageJsonFactory(awsbucket, blogUrl,
	sourceimagepath, markdownfilepath string) *pageJsonFactory {
	if !strings.HasSuffix(blogUrl, "/") {
		blogUrl += "/"
	}
	p := new(pageJsonFactory)
	p.awsbucket = awsbucket
	p.blogUrl = blogUrl
	p.sourceimagepath = sourceimagepath
	p.markdownfilepath = markdownfilepath
	return p
}

type pageJsonFactory struct {
	awsbucket         string
	blogUrl           string
	sourceimagepath   string
	markdownfilepath  string
	uploadimgagepaths []string
	awsimageurls      []string
	imagesizes        []int
}

func (p *pageJsonFactory) AddImageSize(size int) {
	p.imagesizes = append(p.imagesizes, size)
}

func (p *pageJsonFactory) treatImages() {
	imgdir := GetPathWithoutFilename(p.sourceimagepath)
	i := img.NewImg(p.sourceimagepath, imgdir)
	paths := i.PrepareResizeTo(p.imagesizes...)
	p.uploadimgagepaths = append(paths, p.sourceimagepath)
	i.Execute()
}

func (p *pageJsonFactory) uploadImages() {
	for _, filepath := range p.uploadimgagepaths {
		filename := GetFilenameFromPath(filepath)
		key := p.getS3Key(filename)
		url := aws.UploadFile(filepath, p.awsbucket, key)
		p.awsimageurls = append(p.awsimageurls, url)
	}
}

func (p *pageJsonFactory) getS3Key(filename string) string {
	return "blog/" + p.generateDatePath() + filename
}

func (p *pageJsonFactory) generateBlogUrl(title string) string {
	return p.blogUrl + p.generateDatePath() + title
}

func (p *pageJsonFactory) generateDatePath() string {
	now := time.Now()
	return fmt.Sprintf("%d/%d/%d/", now.Year(), now.Month(), now.Day())
}

func (p *pageJsonFactory) generateContentFromMarkdown(input string) string {
	bytes := []byte(input)
	htmlBytes := blackfriday.Run(bytes, blackfriday.WithNoExtensions())
	return string(htmlBytes)

}

func GetPathWithoutFilename(path string) string {
	if _, err := os.Stat(path); err == nil {
		parts := strings.Split(path, "/")
		newpath := strings.Join(parts[:len(parts)-1], "/")
		return newpath + "/"
	}
	log.Fatalln("Not a valid path", path)
	return ""
}

func GetFilenameFromPath(path string) string {
	if _, err := os.Stat(path); err == nil {
		parts := strings.Split(path, "/")
		return parts[len(parts)-1]
	}
	log.Fatalln("Not a valid path", path)
	return ""
}
