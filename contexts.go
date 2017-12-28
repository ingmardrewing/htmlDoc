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
	GetElements() []Element
	SetElements([]Element)
	AddComponent(c component)
	GetComponents() []component
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
	fmt.Println("YYYYY")
	fmt.Println(footerNavi[0].GetPath())
	fmt.Println("YYYYY")
	gc.footerNavigationLocations = footerNavi
}

func (gc *GlobalContext) WriteTo(targetDir string) {
	for _, p := range gc.pages {
		for _, c := range gc.components {
			p.acceptVisitor(c)
		}

		path := targetDir + p.GetFsPath()
		filename := p.GetFsFilename()
		html := p.Render() + gc.GetJs()
		fs.WriteStringToFS(path, filename, html)
	}
}

func (gc *GlobalContext) AddLocations(l []Location) {
	//
}

func (gc *GlobalContext) SetElements(pages []Element) {
	gc.pages = pages
}

func (gc *GlobalContext) GetComponents() []component {
	return gc.components
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
	css := ""
	for _, c := range gc.components {
		css += c.GetCss()
	}
	return gc.minifyCss(css)
}

func (gc *GlobalContext) GetJs() string {
	jsCode := ""
	for _, c := range gc.components {
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
	bc.AddComponent(NewGeneralMetaComponent())
	bc.AddComponent(NewFaviconComponent())
	bc.AddComponent(NewGlobalCssComponent())
	bc.AddComponent(NewGoogleComponent(bc))
	bc.AddComponent(NewTwitterComponent(bc))
	bc.AddComponent(NewFBComponent(bc))
	bc.AddComponent(NewCssLinkComponent(bc.GetCssUrl()))
	bc.AddComponent(NewTitleComponent())
	bc.AddComponent(NewContentComponent())
	bc.AddComponent(NewMainHeaderComponent(bc))
	bc.AddComponent(NewMainNaviComponent(bc.GetMainNavigationLocations()))

	bc.AddComponent(NewDisqusComponent(bc.GetDisqusShortname()))

	bc.AddComponent(NewCopyRightComponent())
	bc.AddComponent(NewFooterNaviComponent(bc.GetFooterNavigationLocations()))
}

/* Footer Context */

func NewFooterContext() *FooterContext {
	return new(FooterContext)
}

type FooterContext struct {
	GlobalContext
}

func (bc *FooterContext) AddComponents() {
	bc.AddComponent(NewGeneralMetaComponent())
	bc.AddComponent(NewFaviconComponent())
	bc.AddComponent(NewGlobalCssComponent())
	bc.AddComponent(NewGoogleComponent(bc))
	bc.AddComponent(NewTwitterComponent(bc))
	bc.AddComponent(NewFBComponent(bc))
	bc.AddComponent(NewCssLinkComponent(bc.GetCssUrl()))
	bc.AddComponent(NewTitleComponent())

	bc.AddComponent(NewContentComponent())
	bc.AddComponent(NewMainHeaderComponent(bc))
	bc.AddComponent(NewMainNaviComponent(bc.GetMainNavigationLocations()))
	bc.AddComponent(NewCopyRightComponent())
	fnl := bc.GetFooterNavigationLocations()
	fmt.Println("XXXX")
	fmt.Println(fnl[0].GetPath())
	fmt.Println(fnl[1].GetPath())
	fmt.Println("XXXX")
	bc.AddComponent(NewFooterNaviComponent(fnl))
}

/* Blog Navi Context */

func NewBlogNaviContext() *BlogNaviContext {
	bn := new(BlogNaviContext)
	return bn
}

type BlogNaviContext struct {
	GlobalContext
	context *BlogContext
}

func (bn *BlogNaviContext) AddComponents() {
	// header
	bn.AddComponent(NewGeneralMetaComponent())
	bn.AddComponent(NewFaviconComponent())
	bn.AddComponent(NewGoogleComponent(bn))
	bn.AddComponent(NewTwitterComponent(bn))
	bn.AddComponent(NewFBComponent(bn))
	bn.AddComponent(NewCssLinkComponent(bn.GetCssUrl()))
	bn.AddComponent(NewTitleComponent())

	// body
	bn.AddComponent(NewBlogNaviContextComponent(bn))
	bn.AddComponent(NewMainHeaderComponent(bn))
	bn.AddComponent(NewMainNaviComponent(bn.GetMainNavigationLocations()))
	bn.AddComponent(NewCopyRightComponent())
	bn.AddComponent(NewFooterNaviComponent(bn.GetFooterNavigationLocations()))
}
