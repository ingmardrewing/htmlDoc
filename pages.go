package htmlDoc

import (
	"log"
	"strconv"
	"strings"

	"github.com/ingmardrewing/fs"
)

type Location interface {
	GetPath() string
	GetDomain() string
	GetTitle() string
	GetThumbnailUrl() string
	GetFsPath() string
}

type Element interface {
	Location
	acceptVisitor(v component)
	addBodyNodes([]*Node)
	addHeaderNodes([]*Node)
	GetPublishedTime() string
	GetDescription() string
	GetContent() string
	GetImageUrl() string
	GetDisqusId() string
	Render() string
	GetFsFilename() string
}

func NewLocation(url, prodDomain, title, thumbnailUrl, fsPath, fsFilename string) *Loc {
	return &Loc{url, prodDomain, title, thumbnailUrl, fsPath, fsFilename}
}

type Loc struct {
	url          string
	prodDomain   string
	title        string
	thumbnailUrl string
	fsPath       string
	fsFilename   string
}

func (l *Loc) GetDomain() string {
	return l.prodDomain
}

func (l *Loc) GetFsPath() string {
	return l.fsPath
}

func (l *Loc) GetFsFilename() string {
	return l.fsFilename
}

func (l *Loc) GetPath() string {
	return l.url
}

func (l *Loc) GetTitle() string {
	return l.title
}

func (l *Loc) GetThumbnailUrl() string {
	return l.thumbnailUrl
}

/* Page */

type Page struct {
	Loc
	doc           *HtmlDoc
	id            string
	Content       string
	Description   string
	ImageUrl      string
	PublishedTime string
	DisqusId      string
}

func NewPage(
	id, title, description, content,
	thumbUrl, imageUrl, prodDomain,
	path, filename, publishedTime,
	disqusId string) *Page {
	p := &Page{
		Loc: Loc{
			title:        title,
			url:          path + filename,
			prodDomain:   prodDomain,
			thumbnailUrl: thumbUrl,
			fsPath:       path,
			fsFilename:   filename},
		id:            id,
		Description:   description,
		Content:       content,
		ImageUrl:      imageUrl,
		PublishedTime: publishedTime,
		DisqusId:      disqusId,
		doc:           NewHtmlDoc()}
	return p
}

func (p *Page) Render() string {
	return p.doc.Render()
}

func (p *Page) GetDisqusId() string {
	return p.DisqusId
}

func (p *Page) GetContent() string {
	return p.Content
}

func (p *Page) GetDescription() string {
	if p.Description != "" {
		return p.Description
	}
	return " "
}

func (p *Page) GetImageUrl() string {
	return p.ImageUrl
}

func (p *Page) GetPublishedTime() string {
	return p.PublishedTime
}

func (p *Page) acceptVisitor(v component) {
	v.visitPage(p)
}

func (p *Page) addHeaderNodes(nodes []*Node) {
	for _, n := range nodes {
		p.doc.AddHeadNode(n)
	}
}

func (p *Page) addBodyNodes(nodes []*Node) {
	for _, n := range nodes {
		p.doc.AddBodyNode(n)
	}
}

/* PageManager */

func NewPageManager() *PageManager {
	return new(PageManager)
}

type PageManager struct {
	Json
	booklistData  *Page
	aboutData     *Page
	imprintData   *Page
	posts         []*Page
	pages         []*Page
	postNaviPages []*Page
}

func (p *PageManager) AddPage(
	id, title, description,
	content, imageUrl, thumbUrl,
	prodDomain, path, filename,
	createDate, disqusId string) {
	page := NewPage(id, title, description, content, imageUrl, thumbUrl, prodDomain, path, filename, createDate, disqusId)
	p.pages = append(p.pages, page)
}

func (p *PageManager) AddPost(
	id, title, description,
	content, imageUrl, thumbUrl,
	prodDomain, path, filename,
	createDate, disqusId string) {
	post := NewPage(id, title, description, content, imageUrl, thumbUrl, prodDomain, path, filename, createDate, disqusId)
	p.posts = append(p.posts, post)
}

func (p *PageManager) GeneratePostNaviPages(atPath string) {
	bundles := GenerateElementBundles(p.posts)
	last := len(bundles) - 1
	for i, b := range bundles {
		ix := strconv.Itoa(i)
		naviPageContent := p.generateNaviPageContent(b)
		filename := "index" + ix + ".html"
		if i == last {
			filename = "index.html"
		}
		pnp := NewPage(ix, "blog navi", "descr ...",
			naviPageContent, "", "", "https://drewing.de",
			atPath, filename, "", "")
		p.AddPostNaviPage(pnp)
	}
}

func (p *PageManager) generateNaviPageContent(bundle *ElementBundle) string {
	n := NewNode("div", "", "class", "blognavipage")
	elems := bundle.GetElements()
	for _, e := range elems {
		ta := e.GetThumbnailUrl()
		if ta == "" {
			ta = e.GetImageUrl()
		}
		a := NewNode("a", " ",
			"href", e.GetPath(),
			"class", "blognavientry__tile")
		span := NewNode("span", " ",
			"style", "background-image: url("+e.GetThumbnailUrl()+")",
			"class", "blognavientry__image")
		h2 := NewNode("h2", e.GetTitle())
		a.AddChild(span)
		a.AddChild(h2)
		n.AddChild(a)
	}
	n.AddChild(NewNode("div", "", "style", "clear: both"))
	return n.Render()
}

func (p *PageManager) GetFooterPages() []Element {
	return []Element{p.imprintData, p.aboutData, p.booklistData}
}

func (p *PageManager) AddPostNaviPage(page *Page) {
	p.postNaviPages = append(p.postNaviPages, page)
}

func (p *PageManager) AddPostFromJsonData(v []byte) {
	id := p.Read(v, "post", "post_id")
	title := p.Read(v, "post", "title")

	thumbUrl := p.Read(v, "thumbImg")
	imageUrl := p.Read(v, "postImg")

	description := p.Read(v, "post", "excerpt")
	disqusId := p.Read(v, "post", "custom_fields", "dsq_thread_id", "[0]")
	createDate := p.Read(v, "post", "date")
	content := p.Read(v, "post", "content")
	rawUrl := p.Read(v, "post", "url")
	path := p.extractPathFromUrl(rawUrl)

	filename := "index.html"
	prodDomain := "https://drewing.de"

	post := NewPage(id, title, description, content,
		imageUrl, thumbUrl, prodDomain,
		path, filename, createDate, disqusId)
	p.posts = append(p.posts, post)
}

func (p *PageManager) ReadPageFromJsonData(v []byte, filename string) *Page {
	id := p.Read(v, "page", "post_id")
	title := p.Read(v, "page", "title")

	thumbUrl := p.Read(v, "thumbImg")
	imageUrl := p.Read(v, "postImg")

	description := p.Read(v, "page", "excerpt")
	disqusId := p.Read(v, "page", "custom_fields", "dsq_thread_id", "[0]")
	createDate := p.Read(v, "page", "date")
	content := p.Read(v, "page", "content")
	rawUrl := p.Read(v, "page", "url")
	path := p.extractPathFromUrl(rawUrl)

	if imageUrl == "" {
		imageUrl = "https://www.drewing.de/blog/wp-content/themes/drewing2012/silhouette_ingmar_drewing.png"
	}
	if thumbUrl == "" {
		thumbUrl = "https://www.drewing.de/blog/wp-content/themes/drewing2012/silhouette_ingmar_drewing.png"
	}

	prodDomain := "https://drewing.de"

	return NewPage(id, title, description, content,
		imageUrl, thumbUrl, prodDomain,
		path, filename, createDate, disqusId)
}

func (p *PageManager) GetPosts() []Element {
	return p.convertToElements(p.posts)
}

func (p *PageManager) GetPostNaviPages() []Element {
	return p.convertToElements(p.postNaviPages)
}

func (p *PageManager) GetPages() []Element {
	return p.convertToElements(p.pages)
}

func (p *PageManager) convertToElements(pages []*Page) []Element {
	elements := []Element{}
	for _, page := range pages {
		elements = append(elements, page)
	}
	return elements
}

func (p *PageManager) extractPathFromUrl(raw string) string {
	s := strings.Split(raw, "/")
	if len(s) < 3 {
		log.Fatalf("url too short in %s: %s\n", raw)
	}
	return "/" + strings.Join(s[3:], "/")
}

func (p *PageManager) GetMainNaviLocations(config []byte) []Location {
	blog := NewLocation(
		"/blog/index.html",
		"",
		"Blog",
		"",
		"",
		"")

	fb := NewLocation(
		p.Read(config, "context", "fbPage"),
		"",
		"Facebook",
		"",
		"",
		"")

	twitter := NewLocation(
		p.Read(config, "context", "twitterPage"),
		"",
		"Twitter",
		"",
		"",
		"")
	// TODO add rss feed
	return []Location{blog, fb, twitter}
}

func (p *PageManager) GetFooterNaviLocations(config []byte) []Location {
	fbs := p.Read(config, "facebookShare")
	fbShare := NewLocation(fbs, "", "Share on Facebook", "", "", "")

	taf := p.Read(config, "tellAFriend")
	tellAFriend := NewLocation(taf, "", "Tell a friend", "", "", "")

	return []Location{fbShare, tellAFriend, p.imprintData, p.aboutData, p.booklistData}
}

func (p *PageManager) ReadFooterData(v []byte) {
	imprintPath := p.Read(v, "imprint")
	imprintBytes := fs.ReadByteArrayFromFile(imprintPath)
	p.imprintData = p.ReadPageFromJsonData(imprintBytes, "/imprint.html")

	booklistPath := p.Read(v, "booklist")
	booklistBytes := fs.ReadByteArrayFromFile(booklistPath)
	p.booklistData = p.ReadPageFromJsonData(booklistBytes, "/booklist.html")

	aboutPath := p.Read(v, "about")
	aboutBytes := fs.ReadByteArrayFromFile(aboutPath)
	p.aboutData = p.ReadPageFromJsonData(aboutBytes, "/about.html")
}

func GenerateElementBundles(pages []*Page) []*ElementBundle {
	length := len(pages)
	reversed := []*Page{}
	for i := length - 1; i >= 0; i-- {
		reversed = append(reversed, pages[i])
	}

	b := NewElementBundle()
	bundles := []*ElementBundle{}
	for _, p := range reversed {
		b.AddElement(p)
		if b.full() {
			bundles = append(bundles, b)
			b = NewElementBundle()
		}
	}
	if !b.full() {
		bundles = append(bundles, b)
	}

	length = len(bundles)
	revbundles := []*ElementBundle{}
	for i := length - 1; i >= 0; i-- {
		revbundles = append(revbundles, bundles[i])
	}
	return revbundles
}

func NewElementBundle() *ElementBundle {
	return new(ElementBundle)
}

type ElementBundle struct {
	elements []Element
}

func (l *ElementBundle) AddElement(e Element) {
	l.elements = append(l.elements, e)
}

func (l *ElementBundle) full() bool {
	return len(l.elements) >= 10
}

func (l *ElementBundle) GetElements() []Element {
	return l.elements
}
