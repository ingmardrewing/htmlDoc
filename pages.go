package htmlDoc

/* location */
type Element interface {
	Location
	acceptVisitor(v visitor)
	addBodyNodes([]*Node)
	addHeaderNodes([]*Node)
	GetPublishedTime() string
	GetDescription() string
	GetImageUrl() string
	GetDisqusId() string
	AddComponent(c component)
}

type Location interface {
	GetUrl() string
	GetTitle() string
	GetThumbnailUrl() string
}

func NewLocation(url, title, thumbnailUrl string) *Loc {
	return &Loc{url, title, thumbnailUrl}
}

type Loc struct {
	url          string
	title        string
	thumbnailUrl string
}

func (l *Loc) GetUrl() string {
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
	components    []component
	Description   string
	ImageUrl      string
	PublishedTime string
	DisqusId      string
}

func NewPage(title, description, url, imageUrl, publishedTime, thumbnailUrl, disqusId string) *Page {
	p := &Page{
		components:    []component{},
		Description:   description,
		ImageUrl:      imageUrl,
		PublishedTime: publishedTime,
		DisqusId:      disqusId,
		doc:           NewHtmlDoc()}
	p.Loc.title = title
	p.Loc.url = url
	p.Loc.thumbnailUrl = thumbnailUrl
	return p
}

func (p *Page) GetDisqusId() string {
	return p.DisqusId
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

func (p *Page) acceptVisitor(v visitor) {
	v.visitPage(p)
}

func (p *Page) Render() string {
	for _, c := range p.components {
		p.acceptVisitor(c)
	}
	return p.doc.Render()
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

func (p *Page) AddComponent(c component) {
	p.components = append(p.components, c)
}

/* context */

type Context interface {
	GetTwitterHandle() string
	GetContentSection() string
	GetContentTags() string
	GetSiteName() string
	GetTwitterCardType() string
	GetOGType() string
	GetFBPageUrl() string
	GetTwitterPage() string
	GetCssUrl() string
	GetDisqusShortname() string
	GetMainNavigationLocations() []Location
	GetReadNavigationLocations() []Location
	AddPage(p Element)
	PreparePages()
}

func NewBlogContext(twitterHandle string,
	contentSection string,
	contentTags string,
	siteName string,
	twitterCardType string,
	ogType string,
	fbPageUrl string,
	twitterPageUrl string,
	cssUrl string,
	disqusShortname string,
	mainNavigationLocations []Location,
	readNavigationLocations []Location) *BlogContext {
	return &BlogContext{
		twitterHandle,
		contentSection,
		contentTags,
		siteName,
		twitterCardType,
		ogType,
		fbPageUrl,
		twitterPageUrl,
		cssUrl,
		disqusShortname,
		mainNavigationLocations,
		readNavigationLocations,
		[]Element{}}
}

type BlogContext struct {
	twitterHandle           string
	contentSection          string
	contentTags             string
	siteName                string
	twitterCardType         string
	ogType                  string
	fbPageUrl               string
	twitterPageUrl          string
	cssUrl                  string
	disqusShortname         string
	mainNavigationLocations []Location
	readNavigationLocations []Location
	pages                   []Element
}

func (bc *BlogContext) PreparePages() {
	for _, p := range bc.pages {
		bc.renderPage(p)
	}
}

func (bc *BlogContext) renderPage(p Element) {
	p.AddComponent(NewGoogleComponent(bc))
	p.AddComponent(NewTwitterComponent(bc))
	p.AddComponent(NewFBComponent(bc))
	p.AddComponent(NewCssLinkComponent(bc.GetCssUrl()))
	p.AddComponent(NewTitleComponent(p.GetTitle()))
	p.AddComponent(NewMainHeaderComponent(bc))
	p.AddComponent(NewGalleryComponent())
	p.AddComponent(NewNaviComponent(bc.GetMainNavigationLocations()))
	p.AddComponent(NewReadNaviComponent(bc.GetReadNavigationLocations()))
	p.AddComponent(NewDisqusComponent(bc.GetDisqusShortname(), p.GetDisqusId()))
}

func (bc *BlogContext) AddPage(p Element) {
	bc.pages = append(bc.pages, p)
}

func (bc *BlogContext) GetDisqusShortname() string {
	return bc.disqusShortname
}

func (bc *BlogContext) GetMainNavigationLocations() []Location {
	return bc.mainNavigationLocations
}

func (bc *BlogContext) GetReadNavigationLocations() []Location {
	return bc.readNavigationLocations
}

func (bc *BlogContext) GetCssUrl() string {
	return bc.cssUrl
}
func (bc *BlogContext) GetTwitterPage() string {
	return bc.twitterPageUrl
}

func (bc *BlogContext) GetFBPageUrl() string {
	return bc.fbPageUrl
}

func (bc *BlogContext) GetOGType() string {
	return bc.ogType
}

func (bc *BlogContext) GetTwitterCardType() string {
	return bc.twitterCardType
}

func (bc *BlogContext) GetTwitterHandle() string {
	return bc.twitterHandle
}

func (bc *BlogContext) GetContentSection() string {
	return bc.contentSection
}

func (bc *BlogContext) GetContentTags() string {
	return bc.contentTags
}

func (bc *BlogContext) GetSiteName() string {
	return bc.siteName
}
