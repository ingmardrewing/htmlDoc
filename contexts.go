package htmlDoc

import (
	"fmt"

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
	AddComponent(c component)
	WriteTo(targetDir string)
	AddPage(p Element)
}

/* Abstract Context */

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

func NewBlogContext(twitterHandle, topic, tags, site, cardType, section, fbPage, twitterPage, cssUrl, rssUrl, home, disqusShortname string, mainNavi, footerNavi []Location) *BlogContext {
	bc := new(BlogContext)

	bc.GlobalContext.twitterHandle = twitterHandle
	bc.GlobalContext.contentSection = topic
	bc.GlobalContext.tags = tags
	bc.GlobalContext.siteName = site
	bc.GlobalContext.twitterCardType = cardType
	bc.GlobalContext.ogType = section
	bc.GlobalContext.fbPageUrl = fbPage
	bc.GlobalContext.twitterPageUrl = twitterPage
	bc.GlobalContext.cssUrl = cssUrl
	bc.GlobalContext.rssUrl = rssUrl
	bc.GlobalContext.homeUrl = home
	bc.GlobalContext.disqusShortname = disqusShortname
	bc.GlobalContext.mainNavigationLocations = mainNavi
	bc.GlobalContext.footerNavigationLocations = footerNavi

	bc.addComponents()
	return bc
}

type BlogContext struct {
	GlobalContext
	readNavigationLocations []Location
}

func (bc *BlogContext) addComponents() {
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
	bn.AddComponent(NewGoogleComponent(bn))
	bn.AddComponent(NewTwitterComponent(bn))
	bn.AddComponent(NewFBComponent(bn))
	bn.AddComponent(NewCssLinkComponent(bn.GetCssUrl()))
	bn.AddComponent(NewTitleComponent())
	bn.AddComponent(NewMainHeaderComponent(bn))
	bn.AddComponent(NewMainNaviComponent(bn.GetMainNavigationLocations()))

	//bn.AddComponent(NewBlogNaviContextComponent(bn))

	bn.AddComponent(NewCopyRightComponent())
	bn.AddComponent(NewFooterNaviComponent(bn.GetFooterNavigationLocations()))
	return bn
}

type BlogNaviContext struct {
	GlobalContext
	context   *BlogContext
	locations []Location
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
	fs.WriteStringToFS(targetDir, bn.GetCssUrl(), bn.GetCss())
}

func (bn *BlogNaviContext) AddComponent(c component) {
	bn.GlobalContext.components = append(bn.components, c)
}

func (bn *BlogNaviContext) AddLocation(p Location) bool {
	bn.locations = append(bn.locations, p)
	return len(bn.locations) == 10
}

func (bn *BlogNaviContext) Render(id string) string {
	p := NewPage(id, "Blog", "Blog",
		"", "", "", "",
		"/blog", "/index"+id+".html",
		"", "")
	for _, c := range bn.components {
		p.acceptVisitor(c)
	}
	return p.Render()
}
