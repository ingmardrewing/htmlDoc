package htmlDoc

import "strconv"

/* location */
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

type Location interface {
	GetPath() string
	GetDomain() string
	GetTitle() string
	GetThumbnailUrl() string
	GetFsPath() string
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

/* BlogPage */

func NewBlogNaviPage(
	id, title, description, content,
	imageUrl, thumbUrl, prodDomain,
	path, filename, publishedTime,
	disqusId string,
	bundle *LocationBundle) *BlogNavPage {
	p := &BlogNavPage{
		Page: Page{
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
			doc:           NewHtmlDoc()},
		bundle: bundle}
	return p
}

type BlogNavPage struct {
	Page
	bundle *LocationBundle
}

func (b *BlogNavPage) SetBundle(bundle *LocationBundle) {
	b.bundle = bundle
}

func (b *BlogNavPage) GetBundle() *LocationBundle {
	return b.bundle
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
		id:            id,
		Description:   description,
		Content:       content,
		ImageUrl:      imageUrl,
		PublishedTime: publishedTime,
		DisqusId:      disqusId,
		doc:           NewHtmlDoc()}
	p.Loc.title = title
	p.Loc.url = path + filename
	p.Loc.prodDomain = prodDomain
	p.Loc.thumbnailUrl = thumbUrl
	p.Loc.fsPath = path
	p.Loc.fsFilename = filename
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
	post := NewPage(id, title, description, content, imageUrl, thumbUrl, prodDomain, path, filename, createDate, disqusId)
	p.posts = append(p.posts, post)
}

func (p *PageManager) GeneratePostNaviPages() {
	// TODO: ugly ...
	loc := []Location{}
	for _, post := range p.posts {
		loc = append(loc, post)
	}
	n := NewNaviPageFactory(loc)
	bundles := n.generateBundles()

	for i, b := range bundles {
		ix := strconv.Itoa(i)
		naviPageContent := p.generateNaviPageContent(b)
		pnp := NewPage(ix, "blog navi", "descr ...",
			naviPageContent, "", "", "https://drewing.de",
			"/", "index"+ix+".html", "", "")
		p.AddPostNaviPage(pnp)
	}
}

func (p *PageManager) generateNaviPageContent(bundle *LocationBundle) string {
	n := NewNode("div", "", "class", "blognavientry")
	for _, l := range bundle.GetLocations() {
		n.AddChild(NewNode("p", l.GetPath()))
		n.AddChild(NewNode("p", l.GetDomain()))
		n.AddChild(NewNode("p", l.GetTitle()))
		n.AddChild(NewNode("p", l.GetThumbnailUrl()))
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

func NewNaviPageFactory(l []Location) *NaviPageFactory {
	np := new(NaviPageFactory)
	np.locations = l
	return np
}

type NaviPageFactory struct {
	locations []Location
}

func (n NaviPageFactory) generateBundles() []*LocationBundle {
	b := NewLocationBundle()
	bundles := []*LocationBundle{}
	for _, l := range n.locations {
		b.AddLocation(l)
		if b.full() {
			bundles = append(bundles, b)
			b = NewLocationBundle()
		}
	}
	return bundles
}

func NewLocationBundle() *LocationBundle {
	return new(LocationBundle)
}

type LocationBundle struct {
	locations []Location
}

func (l *LocationBundle) AddLocation(p Location) {
	l.locations = append(l.locations, p)
}

func (l *LocationBundle) full() bool {
	return len(l.locations) >= 10
}

func (l *LocationBundle) GetLocations() []Location {
	return l.locations
}
