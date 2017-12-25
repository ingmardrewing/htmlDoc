package htmlDoc

import (
	"fmt"
	"strconv"

	"github.com/ingmardrewing/fs"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

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
	GetCss() string
	GetRssUrl() string
	GetHomeUrl() string
	GetDisqusShortname() string
	GetMainNavigationLocations() []Location
	GetReadNavigationLocations() []Location
	GetFooterNavigationLocations() []Location
	GetElements() []Element
	AddComponent(c component)
	WriteTo(targetDir string)
	AddPage(p Element)
	AddComponents()
	AddLocations(l []Location)
	SetGlobalFields(twitterHandle, topic, tags, site, cardType, section, fbPage, twitterPage, cssUrl, rssUrl, home, disqusShortname string, mainNavi, footerNavi []Location)
}

/* Global Context */

type GlobalContext struct {
	twitterHandle             string
	contentSection            string
	tags                      string
	siteName                  string
	twitterCardType           string
	ogType                    string
	fbPageUrl                 string
	twitterPageUrl            string
	cssUrl                    string
	rssUrl                    string
	homeUrl                   string
	disqusShortname           string
	mainNavigationLocations   []Location
	footerNavigationLocations []Location
	pages                     []Element
	components                []component
}

func (gc *GlobalContext) SetGlobalFields(
	twitterHandle,
	topic,
	tags,
	site,
	cardType,
	section,
	fbPage,
	twitterPage,
	cssUrl,
	rssUrl,
	home,
	disqusShortname string,
	mainNavi,
	footerNavi []Location) {
	gc.twitterHandle = twitterHandle
	gc.contentSection = topic
	gc.tags = tags
	gc.siteName = site
	gc.twitterCardType = cardType
	gc.ogType = section
	gc.fbPageUrl = fbPage
	gc.twitterPageUrl = twitterPage
	gc.cssUrl = cssUrl
	gc.rssUrl = rssUrl
	gc.homeUrl = home
	gc.disqusShortname = disqusShortname
	gc.mainNavigationLocations = mainNavi
	gc.footerNavigationLocations = footerNavi
}

func (gc *GlobalContext) AddLocations(l []Location) {
	//
}

func (gc *GlobalContext) GetElements() []Element {
	return gc.pages
}

func (gc *GlobalContext) AddPage(p Element) {
	gc.pages = append(gc.pages, p)
}

func (gc *GlobalContext) AddComponent(c component) {
	gc.components = append(gc.components, c)
}

func (gc *GlobalContext) GetHomeUrl() string {
	return gc.homeUrl
}

func (gc *GlobalContext) GetRssUrl() string {
	return gc.rssUrl
}

func (gc *GlobalContext) GetDisqusShortname() string {
	return gc.disqusShortname
}

func (gc *GlobalContext) GetMainNavigationLocations() []Location {
	return gc.mainNavigationLocations
}

func (gc *GlobalContext) GetFooterNavigationLocations() []Location {
	return gc.footerNavigationLocations
}

func (gc *GlobalContext) GetCssUrl() string {
	return gc.cssUrl
}

func (gc *GlobalContext) GetTwitterPage() string {
	return gc.twitterPageUrl
}

func (gc *GlobalContext) GetFBPageUrl() string {
	return gc.fbPageUrl
}

func (gc *GlobalContext) GetOGType() string {
	return gc.ogType
}

func (gc *GlobalContext) GetTwitterCardType() string {
	return gc.twitterCardType
}

func (gc *GlobalContext) GetTwitterHandle() string {
	return gc.twitterHandle
}

func (gc *GlobalContext) GetContentSection() string {
	return gc.contentSection
}

func (gc *GlobalContext) GetContentTags() string {
	return gc.tags
}

func (gc *GlobalContext) GetSiteName() string {
	return gc.siteName
}

func (gc *GlobalContext) GetCss() string {
	css := `
body, p, span {
	margin: 0;
	padding: 0;
	font-family: Arial, Helvetica, sans-serif;
}
a {
	color: grey;
	text-decoration: none;
}
a:hover {
	text-decoration: underline;
}
.wrapperOuter {
	text-align: center;
}

.wrapperInner {
	margin: 0 auto;
	width: 800px;
}
`
	for _, c := range gc.components {
		css += c.GetCss()
	}
	return gc.minifyCss(css)
}

func (gc *GlobalContext) GetJs() string {
	return ""
}

func (gc *GlobalContext) minifyCss(txt string) string {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	s, err := m.String("text/css", txt)
	if err != nil {
		panic(err)
	}
	return s
}

func (gc *GlobalContext) GetReadNavigationLocations() []Location {
	return nil
}

/* Blog Context */

func NewBlogContext() *BlogContext {
	bc := new(BlogContext)
	return bc
}

type BlogContext struct {
	GlobalContext
	readNavigationLocations []Location
}

func (bc *BlogContext) AddComponents() {
	bc.AddComponent(NewGoogleComponent(bc))
	bc.AddComponent(NewTwitterComponent(bc))
	bc.AddComponent(NewFBComponent(bc))
	bc.AddComponent(NewCssLinkComponent(bc.GetCssUrl()))
	bc.AddComponent(NewTitleComponent())
	bc.AddComponent(NewMainHeaderComponent(bc))
	bc.AddComponent(NewMainNaviComponent(bc.GetMainNavigationLocations()))
	bc.AddComponent(NewContentComponent())
	bc.AddComponent(NewDisqusComponent(bc.GetDisqusShortname()))
	bc.AddComponent(NewCopyRightComponent())
	bc.AddComponent(NewFooterNaviComponent(bc.GetFooterNavigationLocations()))
}

func (bc *BlogContext) GetReadNavigationLocations() []Location {
	return bc.readNavigationLocations
}

func (bc *BlogContext) WriteTo(targetDir string) {
	for _, p := range bc.pages {
		for _, c := range bc.components {
			p.acceptVisitor(c)
		}

		path := targetDir + p.GetFsPath()
		filename := p.GetFsFilename()
		html := p.Render() + bc.GetJs()
		fs.WriteStringToFS(path, filename, html)

	}
	fs.WriteStringToFS(targetDir, bc.GetCssUrl(), bc.GetCss())
}

func (bc *BlogContext) minifyHtml(html string) string {
	// TODO: find error - produces defect html
	m := minify.New()
	m.AddFunc("text/html", js.Minify)
	s, err := m.String("text/html", html)
	if err != nil {
		panic(err)
	}
	return s
}

func (bc *BlogContext) GetJs() string {
	jsCode := ""
	for _, c := range bc.components {
		jsCode += c.GetJs()
	}

	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	s, err := m.String("text/javascript", jsCode)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`<script>%s</script>`, s)
}

/* Blog Navi Context */

func NewBlogNaviContext() *BlogNaviContext {
	bn := new(BlogNaviContext)
	return bn
}

type BlogNaviContext struct {
	GlobalContext
	context   *BlogContext
	locations []Location
}

func (bn *BlogNaviContext) AddComponents() {
	bn.AddComponent(NewGoogleComponent(bn))
	bn.AddComponent(NewTwitterComponent(bn))
	bn.AddComponent(NewFBComponent(bn))
	bn.AddComponent(NewCssLinkComponent(bn.GetCssUrl()))
	bn.AddComponent(NewTitleComponent())
	bn.AddComponent(NewMainHeaderComponent(bn))
	bn.AddComponent(NewMainNaviComponent(bn.GetMainNavigationLocations()))
	bn.AddComponent(NewBlogNaviContextComponent())

	bn.AddComponent(NewCopyRightComponent())
	bn.AddComponent(NewFooterNaviComponent(bn.GetFooterNavigationLocations()))
}

func (bn *BlogNaviContext) WriteTo(targetDir string) {
	for _, p := range bn.pages {
		for _, c := range bn.components {
			p.acceptVisitor(c)
		}

		path := targetDir + p.GetFsPath()
		filename := p.GetFsFilename()
		html := p.Render() + bn.GetJs()
		fs.WriteStringToFS(path, filename, html)
	}
}

func (bn *BlogNaviContext) AddComponent(c component) {
	bn.components = append(bn.components, c)
}

func (bn *BlogNaviContext) AddLocations(l []Location) {
	nf := NewNaviPageFactory(l)
	bundles := nf.generateBundles()

	for i, b := range bundles {
		ix := strconv.Itoa(i)
		bn.AddPage(NewPage(ix, "blog navi", "descr ...",
			b.GetContent(), "", "", "https://drewing.de",
			"/", "index"+ix+".html", "", ""))
	}
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

func (l *LocationBundle) GetContent() string {
	html := ""
	for _, l := range l.locations {
		html += fmt.Sprintf("<p>%s</p>", l.GetDomain())
		html += fmt.Sprintf("<p>%s</p>", l.GetPath())
		html += fmt.Sprintf("<p>%s</p>", l.GetTitle())
		html += fmt.Sprintf("<p>%s</p>", l.GetThumbnailUrl())
		html += fmt.Sprintf("<p>%s</p>", l.GetFsPath())
	}
	return html
}
