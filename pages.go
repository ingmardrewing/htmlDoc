package htmlDoc

/* location */
type Element interface {
	Location
	acceptVisitor(v component)
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
