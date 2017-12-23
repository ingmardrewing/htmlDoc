package htmlDoc

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
	GetFooterNavigationLocations() []Location
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
	disqusShortname           string
	mainNavigationLocations   []Location
	readNavigationLocations   []Location
	footerNavigationLocations []Location
	pages                     []Element
	components                []component
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
