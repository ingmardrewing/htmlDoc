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
	Render() string
}

type Location interface {
	GetUrl() string
	GetTitle() string
	GetThumbnailUrl() string
	GetFsPath() string
}

func NewLocation(url, title, thumbnailUrl, fsPath, fsFilename string) *Loc {
	return &Loc{url, title, thumbnailUrl, fsPath, fsFilename}
}

type Loc struct {
	url          string
	title        string
	thumbnailUrl string
	fsPath       string
	fsFilename   string
}

func (l *Loc) GetFsPath() string {
	return l.fsPath
}

func (l *Loc) GetFsFilename() string {
	return l.fsFilename
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
	id            string
	Description   string
	ImageUrl      string
	PublishedTime string
	DisqusId      string
}

func NewPage(
	id, title, description,
	imageUrl, thumbUrl, path, filename,
	publishedTime, disqusId string) *Page {
	p := &Page{
		id:            id,
		Description:   description,
		ImageUrl:      imageUrl,
		PublishedTime: publishedTime,
		DisqusId:      disqusId,
		doc:           NewHtmlDoc()}
	p.Loc.title = title
	p.Loc.url = path + filename
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
	AddComponent(c component)
	Render(p Element) string
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
	bc := &BlogContext{
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
		[]Element{},
		[]component{}}
	bc.AddComponent(NewGoogleComponent(bc))
	bc.AddComponent(NewTwitterComponent(bc))
	bc.AddComponent(NewFBComponent(bc))
	bc.AddComponent(NewCssLinkComponent(bc.GetCssUrl()))
	bc.AddComponent(NewTitleComponent())
	bc.AddComponent(NewMainHeaderComponent(bc))
	bc.AddComponent(NewGalleryComponent())
	bc.AddComponent(NewNaviComponent(bc.GetMainNavigationLocations()))
	bc.AddComponent(NewReadNaviComponent(bc.GetReadNavigationLocations()))
	bc.AddComponent(NewDisqusComponent(bc.GetDisqusShortname()))
	bc.AddComponent(NewCopyRightComponent())
	return bc

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
	components              []component
}

func (bc *BlogContext) AddComponent(c component) {
	bc.components = append(bc.components, c)
}

func (bc *BlogContext) Render(p Element) string {
	for _, c := range bc.components {
		p.acceptVisitor(c)
	}
	return p.Render()
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
