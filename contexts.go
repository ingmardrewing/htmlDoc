package htmlDoc

import (
	"fmt"
	"log"

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
	Render(targetDir string)
	AddPage(p Element)
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
	rssUrl string,
	homeUrl string,
	disqusShortname string,
	mainNavigationLocations []Location,
	readNavigationLocations []Location,
	footerNavigationLocations []Location) *BlogContext {
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
		rssUrl,
		homeUrl,
		disqusShortname,
		mainNavigationLocations,
		readNavigationLocations,
		footerNavigationLocations,
		[]Element{},
		[]component{}}
	bc.AddComponent(NewGoogleComponent(bc))
	bc.AddComponent(NewTwitterComponent(bc))
	bc.AddComponent(NewFBComponent(bc))
	bc.AddComponent(NewCssLinkComponent(bc.GetCssUrl()))
	bc.AddComponent(NewTitleComponent())
	bc.AddComponent(NewMainHeaderComponent(bc))
	bc.AddComponent(NewMainNaviComponent(mainNavigationLocations))
	//bc.AddComponent(NewGalleryComponent())
	bc.AddComponent(NewMainNaviComponent(bc.GetMainNavigationLocations()))
	bc.AddComponent(NewContentComponent())
	//bc.AddComponent(NewReadNaviComponent(bc.GetReadNavigationLocations()))
	bc.AddComponent(NewDisqusComponent(bc.GetDisqusShortname()))
	bc.AddComponent(NewCopyRightComponent())
	bc.AddComponent(NewFooterNaviComponent(bc.GetFooterNavigationLocations()))
	return bc
}

type BlogContext struct {
	twitterHandle             string
	contentSection            string
	contentTags               string
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
	readNavigationLocations   []Location
	footerNavigationLocations []Location
	pages                     []Element
	components                []component
}

func (bc *BlogContext) GetCss() string {
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
	for _, c := range bc.components {
		css += c.GetCss()
	}
	return bc.minifyCss(css)
}

func (bc *BlogContext) minifyCss(txt string) string {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	s, err := m.String("text/css", txt)
	if err != nil {
		panic(err)
	}
	return s
}

func (bc *BlogContext) AddPage(p Element) {
	bc.pages = append(bc.pages, p)
}

func (bc *BlogContext) AddComponent(c component) {
	bc.components = append(bc.components, c)
}

func (bc *BlogContext) GetHomeUrl() string {
	return bc.homeUrl
}
func (bc *BlogContext) GetRssUrl() string {
	return bc.rssUrl
}

func (bc *BlogContext) Render(targetDir string) {
	for _, p := range bc.pages {
		path := targetDir + p.GetFsPath()
		log.Println("Writing to " + path)
		filename := p.GetFsFilename()
		for _, c := range bc.components {
			p.acceptVisitor(c)
		}
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

func (bc *BlogContext) GetDisqusShortname() string {
	return bc.disqusShortname
}

func (bc *BlogContext) GetMainNavigationLocations() []Location {
	return bc.mainNavigationLocations
}

func (bc *BlogContext) GetReadNavigationLocations() []Location {
	return bc.readNavigationLocations
}

func (bc *BlogContext) GetFooterNavigationLocations() []Location {
	return bc.footerNavigationLocations
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
