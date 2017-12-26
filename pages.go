package htmlDoc

import (
	"fmt"
	"strconv"
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
	imageUrl, thumbUrl, prodDomain,
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
	return p.Description
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
	fmt.Println(thumbUrl)
	fmt.Println(imageUrl)
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
	n := NewNode("div", "", "class", "blognavientry")
	for _, e := range bundle.GetElements() {
		h2 := NewNode("h2", e.GetTitle())
		h3 := NewNode("h2", e.GetDescription())
		div := NewNode("div", "", "style", "background-image: url("+e.GetThumbnailUrl()+")")
		div.AddChild(h2)
		div.AddChild(h3)

		a := NewNode("a", "", "href", e.GetDomain()+e.GetPath())
		a.AddChild(div)

		n.AddChild(a)

	}
	return n.Render()
}

func (p *PageManager) AddPostNaviPage(page *Page) {
	p.postNaviPages = append(p.postNaviPages, page)
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

func GenerateElementBundles(pages []*Page) []*ElementBundle {
	b := NewElementBundle()
	bundles := []*ElementBundle{}
	for _, p := range pages {
		b.AddElement(p)
		if b.full() {
			bundles = append(bundles, b)
			b = NewElementBundle()
		}
	}
	return bundles
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
	length := len(l.elements)
	reversed := []Element{}
	for i := length - 1; i > 0; i-- {
		reversed = append(reversed, l.elements[i])
	}
	return reversed
}
