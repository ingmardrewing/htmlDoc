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
	AddComponent(c component)
}

type Location interface {
	GetUrl() string
	GetTitle() string
}

func NewLocation(url, title string) *Loc {
	return &Loc{url, title}
}

type Loc struct {
	url   string
	title string
}

func (l *Loc) GetUrl() string {
	return l.url
}

func (l *Loc) GetTitle() string {
	return l.title
}

/* Page */

type Page struct {
	Loc
	doc           *HtmlDoc
	components    []component
	Description   string
	ImageUrl      string
	PublishedTime string
}

func NewPage(title, description, url, imageUrl, publishedTime string) *Page {
	p := &Page{
		components:    []component{},
		Description:   description,
		ImageUrl:      imageUrl,
		PublishedTime: publishedTime,
		doc:           NewHtmlDoc()}
	p.Loc.title = title
	p.Loc.url = url
	return p
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
