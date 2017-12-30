package htmlDoc

import (
	"fmt"
	"testing"
	"time"
)

func TestGetPathWithoutFilename(t *testing.T) {
	path := "/Users/drewing/Desktop/drewing2018/add/atthezoo.png"
	actual := GetPathWithoutFilename(path)
	expected := "/Users/drewing/Desktop/drewing2018/add/"
	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestGetFilenameFromPath(t *testing.T) {
	path := "/Users/drewing/Desktop/drewing2018/add/atthezoo.png"
	actual := GetFilenameFromPath(path)
	expected := "atthezoo.png"
	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestTreatImages(t *testing.T) {
	b := NewPageJsonFactory(
		"", "",
		"/Users/drewing/Desktop/drewing2018/add/atthezoo.png",
		"/Users/drewing/Desktop/drewing2018/add/test.md")
	b.AddImageSize(800)
	b.AddImageSize(390)
	b.treatImages()
}

func TestUploadImages(t *testing.T) {
	b := NewPageJsonFactory(
		"drewingde", "",
		"/Users/drewing/Desktop/drewing2018/add/atthezoo.png",
		"/Users/drewing/Desktop/drewing2018/add/test.md")
	b.AddImageSize(800)
	b.AddImageSize(390)
	b.treatImages()
	b.uploadImages()
}

func TestS3KeyGenerationFromDate(t *testing.T) {
	p := NewPageJsonFactory("", "https://drewing.de/", "", "")
	actual := p.generateDatePath()
	now := time.Now()
	expected := fmt.Sprintf("%d/%d/%d/", now.Year(), now.Month(), now.Day())

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestGenerateContentFromMarkdown(t *testing.T) {
	input := `test`
	p := NewPageJsonFactory("", "https://drewing.de/", "", "")
	actual := p.generateContentFromMarkdown(input)
	expected := "<p>test</p>\n"

	if actual != expected {
		t.Error("Expected", expected, "but got", ">"+actual+"<")
	}
}

func TestGenerateBlogUrl(t *testing.T) {
	now := time.Now()
	d := "https://drewing.de/blog"
	k := fmt.Sprintf("%d/%d/%d/", now.Year(), now.Month(), now.Day())
	title := "just-a-test/"

	p := NewPageJsonFactory("", d, "", "")
	actual := p.generateBlogUrl(title)
	expected := d + "/" + k + title

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}
