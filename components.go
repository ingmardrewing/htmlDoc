package htmlDoc

import (
	"fmt"
	"strings"
)

/* component */
type component interface {
	visitPage(p Element)
	GetCss() string
	GetJs() string
}

type abstractComponent struct{}

func (ac *abstractComponent) GetCss() string {
	return ""
}

func (ac *abstractComponent) GetJs() string {
	return ""
}

/* wrapper */
type wrapper struct{}

func (cw *wrapper) wrap(n *Node, addedclasses ...string) *Node {
	inner := NewNode("div", "", "class", "wrapperInner")
	inner.AddChild(n)
	classes := "wrapperOuter " + strings.Join(addedclasses, " ")
	wrapperNode := NewNode("div", "", "class", classes)
	wrapperNode.AddChild(inner)
	return wrapperNode
}

/* global css component */

func NewGlobalCssComponent() *GlobalCssComponent {
	return new(GlobalCssComponent)
}

type GlobalCssComponent struct {
	abstractComponent
}

func (gcc *GlobalCssComponent) visitPage(p Element) {}

func (gcc *GlobalCssComponent) GetCss() string {
	return `
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
}

/* fb component */
type FBComponent struct {
	abstractComponent
	context Context
}

func NewFBComponent(context Context) *FBComponent {
	fb := new(FBComponent)
	fb.context = context
	return fb
}

func (fbc *FBComponent) visitPage(p Element) {
	m := []*Node{
		NewNode("meta", "", "property", "og:title", "content", p.GetTitle()),
		NewNode("meta", "", "property", "og:url", "content", p.GetPath()),
		NewNode("meta", "", "property", "og:image", "content", p.GetImageUrl()),
		NewNode("meta", "", "property", "og:description", "content", p.GetDescription()),
		NewNode("meta", "", "property", "og:site_name", "content", fbc.context.GetSiteName()),
		NewNode("meta", "", "property", "og:type", "content", fbc.context.GetOGType()),
		NewNode("meta", "", "property", "article:published_time", "content", p.GetPublishedTime()),
		NewNode("meta", "", "property", "article:modified_time", "content", p.GetPublishedTime()),
		NewNode("meta", "", "property", "article:section", "content", fbc.context.GetContentSection()),
		NewNode("meta", "", "property", "article:tag", "content", fbc.context.GetContentTags())}

	p.addHeaderNodes(m)
}

/* google component */

type GoogleComponent struct {
	abstractComponent
	context Context
}

func NewGoogleComponent(context Context) *GoogleComponent {
	gc := new(GoogleComponent)
	gc.context = context
	return gc
}

func (goo *GoogleComponent) visitPage(p Element) {
	m := []*Node{
		NewNode("meta", "", "itemprop", "name", "content", p.GetTitle()),
		NewNode("meta", "", "itemprop", "description", "content", p.GetDescription()),
		NewNode("meta", "", "itemprop", "image", "content", p.GetImageUrl())}
	p.addHeaderNodes(m)
}

/* twitter component */

type TwitterComponent struct {
	abstractComponent
	context Context
}

func NewTwitterComponent(context Context) *TwitterComponent {
	t := new(TwitterComponent)
	t.context = context
	return t
}

func (tw *TwitterComponent) visitPage(p Element) {
	m := []*Node{
		NewNode("meta", "",
			"name", "t:card",
			"content", tw.context.GetTwitterCardType()),
		NewNode("meta", "",
			"name", "t:site",
			"content", tw.context.GetTwitterHandle()),
		NewNode("meta", "",
			"name", "t:title",
			"content", p.GetTitle()),
		NewNode("meta", "",
			"name", "t:text:description",
			"content", p.GetDescription()),
		NewNode("meta", "",
			"name", "t:creator",
			"content", tw.context.GetTwitterHandle()),
		NewNode("meta", "",
			"name", "t:image",
			"content", p.GetImageUrl())}
	p.addHeaderNodes(m)
}

/* title component */
type TitleComponent struct {
	abstractComponent
}

func NewTitleComponent() *TitleComponent {
	return new(TitleComponent)
}

func (tc *TitleComponent) visitPage(p Element) {
	title := NewNode("title", p.GetTitle())
	p.addHeaderNodes([]*Node{title})
}

/* css link component */

type CssLinkComponent struct {
	abstractComponent
	url string
}

func NewCssLinkComponent(url string) *CssLinkComponent {
	clc := new(CssLinkComponent)
	clc.url = "/drewing2018" + url
	return clc
}

func (clc *CssLinkComponent) visitPage(p Element) {
	link := NewNode("link", "", "href", clc.url, "rel", "stylesheet", "type", "text/css")
	p.addHeaderNodes([]*Node{link})
}

/**/

type BlogNaviComponent struct {
	wrapper
	abstractComponent
	context Context
}

func NewBlogNaviContextComponent(c Context) *BlogNaviComponent {
	bnc := new(BlogNaviComponent)
	bnc.context = c
	return bnc
}

/*
func (b *BlogNaviComponent) visitPage(p Element) {
	n := NewNode("div", p.GetContent(), "class", "blognavicomponent")
	wn := b.wrap(n)
	p.addHeaderNodes([]*Node{wn})
}

func (b *BlogNaviComponent) addHeaderNodes(p Element) {
	inx := b.getIndexOfPage(p)
	n := []*Node{}
	firstUrl := b.locations[0].GetPath()
	n = append(n, NewNode("link", "", "rel", "first", "href", firstUrl))
	if inx > 0 {
		prevUrl := b.locations[inx-1].GetPath()
		n = append(n, NewNode("link", "", "rel", "prev", "href", prevUrl))
	}
	if inx < len(b.locations)-1 {
		nextUrl := b.locations[inx+1].GetPath()
		n = append(n, NewNode("link", "", "rel", "next", "href", nextUrl))
	}
	lastUrl := b.locations[len(b.locations)-1].GetPath()
	n = append(n, NewNode("link", "", "rel", "last", "href", lastUrl))
	p.addHeaderNodes(n)
}
*/

func (b *BlogNaviComponent) addPrevious(p Element, n *Node) {
	inx := b.getIndexOfPage(p)
	if inx == 0 {
		span := NewNode("span", "< previous posts")
		n.AddChild(span)
	} else {
		elems := b.context.GetElements()
		pv := elems[inx-1]
		a := NewNode("a", "< previous posts", "href", pv.GetPath(), "rel", "prev")
		n.AddChild(a)
	}
}

func (b *BlogNaviComponent) addNext(p Element, n *Node) {
	inx := b.getIndexOfPage(p)
	if inx == len(b.context.GetElements())-1 {
		span := NewNode("span", "next >")
		n.AddChild(span)
	} else {
		elems := b.context.GetElements()
		nx := elems[inx+1]
		a := NewNode("a", "next >", "href", nx.GetPath(), "rel", "next")
		n.AddChild(a)
	}
}

func (b *BlogNaviComponent) addBodyNodes(p Element) {
	nav := NewNode("nav", "")
	b.addPrevious(p, nav)
	b.addNext(p, nav)

	d := NewNode("div", "", "class", "blognavicomponent")
	d.AddChild(nav)
	d.AddChild(NewNode("div", p.GetContent()))
	d.AddChild(nav)

	wn := b.wrap(d)

	p.addBodyNodes([]*Node{wn})
}

func (b *BlogNaviComponent) visitPage(p Element) {
	if len(b.context.GetElements()) < 3 {
		return
	}
	//b.addHeaderNodes(p)
	b.addBodyNodes(p)

}

func (b *BlogNaviComponent) getIndexOfPage(p Element) int {
	for i, l := range b.context.GetElements() {
		if l.GetPath() == p.GetPath() {
			return i
		}
	}
	return -1
}

func (b *BlogNaviComponent) GetCss() string {
	return `
.blognavicomponent {
	text-align: left;
	padding-top: 200px;
	padding-bottom: 200px;
}
`
}

/* MainNaviComponent */
func NewMainNaviComponent(locations []Location) *MainNaviComponent {
	nc := new(MainNaviComponent)
	nc.locations = locations
	return nc
}

type MainNaviComponent struct {
	wrapper
	locations []Location
	cssClass  string
}

func (nv *MainNaviComponent) visitPage(p Element) {
	nav := NewNode("nav", "",
		"class", "mainnavi")
	url := p.GetPath()
	for _, l := range nv.locations {
		if url == l.GetPath() {
			span := NewNode("span", l.GetTitle(),
				"class", "mainnavi__navelement")
			nav.AddChild(span)
		} else {
			a := NewNode("a", l.GetTitle(),
				"href", l.GetPath(),
				"class", "mainnavi__navelement")
			nav.AddChild(a)
		}
	}
	node := NewNode("div", "", "class", nv.cssClass)
	node.AddChild(nav)
	wn := nv.wrap(node, "mainnavi__wrapper")
	p.addBodyNodes([]*Node{wn})
}

func (mhc *MainNaviComponent) GetJs() string {
	return ""
}

func (mhc *MainNaviComponent) GetCss() string {
	return `
.mainnavi {
	border-top: 1px solid black;
	border-bottom: 2px solid black;
}
.mainnavi__wrapper {
	position: fixed;
	width: 100%;
	top: 80px;
	background-color: white;
}
.mainnavi__navelement {
	display: inline-block;
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	font-weight: 900;
	font-size: 18px;
	line-height: 20px;
	text-transform: uppercase;
	text-decoration: none;
	color: black;
	padding: 10px 20px;
}
.mainnavi__nav {
	border-bottom: 2px solid black;
}
`
}

/* FooterNaviComponent */

func NewFooterNaviComponent(locations []Location) *FooterNaviComponent {
	nc := new(FooterNaviComponent)
	nc.locations = locations
	return nc
}

type FooterNaviComponent struct {
	wrapper
	locations []Location
	cssClass  string
}

func (nv *FooterNaviComponent) visitPage(p Element) {
	nav := NewNode("nav", "",
		"class", "footernavi")
	url := p.GetPath()
	for _, l := range nv.locations {
		if url == l.GetPath() {
			span := NewNode("span", l.GetTitle(),
				"class", "footernavi__navelement")
			nav.AddChild(span)
		} else {
			a := NewNode("a", l.GetTitle(),
				"href", l.GetPath(),
				"class", "footernavi__navelement")
			nav.AddChild(a)
		}
	}
	node := NewNode("div", "", "class", nv.cssClass)
	node.AddChild(nav)
	wn := nv.wrap(node, "footernavi__wrapper")
	p.addBodyNodes([]*Node{wn})
}

func (mhc *FooterNaviComponent) GetJs() string { return "" }

func (mhc *FooterNaviComponent) GetCss() string {
	return `
.footernavi {
	border-top: 1px solid black;
}
.footernavi__wrapper {
	position: fixed;
	width: 100%;
	bottom: 0;
	background-color: white;
}
.footernavi__navelement {
	display: inline-block;
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	font-weight: 900;
	font-size: 16px;
	line-height: 20px;
	text-transform: uppercase;
	text-decoration: none;
	color: black;
	padding: 10px 20px;
}
`
}

/* ReadNaviComponent */

type ReadNaviComponent struct {
	wrapper
	locations []Location
}

func NewReadNaviComponent(locations []Location) *ReadNaviComponent {
	rnc := new(ReadNaviComponent)
	rnc.locations = locations
	return rnc
}

func (rnv *ReadNaviComponent) addFirst(p Element, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == 0 {
		span := NewNode("span", "<< first")
		n.AddChild(span)
	} else {
		f := rnv.locations[0]
		a := NewNode("a", "<< first", "href", f.GetPath(), "rel", "first")
		n.AddChild(a)
	}
}

func (rnv *ReadNaviComponent) addPrevious(p Element, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == 0 {
		span := NewNode("span", "< previous")
		n.AddChild(span)
	} else {
		p := rnv.locations[inx-1]
		a := NewNode("a", "< previous", "href", p.GetPath(), "rel", "prev")
		n.AddChild(a)
	}
}

func (rnv *ReadNaviComponent) addNext(p Element, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == len(rnv.locations)-1 {
		span := NewNode("span", "next >")
		n.AddChild(span)
	} else {
		nx := rnv.locations[inx+1]
		a := NewNode("a", "next >", "href", nx.GetPath(), "rel", "next")
		n.AddChild(a)
	}
}

func (rnv *ReadNaviComponent) addLast(p Element, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == len(rnv.locations)-1 {
		span := NewNode("span", "newest >>")
		n.AddChild(span)
	} else {
		nw := rnv.locations[len(rnv.locations)-1]
		a := NewNode("a", "neweset >>", "href", nw.GetPath(), "rel", "last")
		n.AddChild(a)
	}
}

func (rnv *ReadNaviComponent) addHeaderNodes(p Element) {
	inx := rnv.getIndexOfPage(p)
	n := []*Node{}
	firstUrl := rnv.locations[0].GetPath()
	n = append(n, NewNode("link", "", "rel", "first", "href", firstUrl))
	if inx > 0 {
		prevUrl := rnv.locations[inx-1].GetPath()
		n = append(n, NewNode("link", "", "rel", "prev", "href", prevUrl))
	}
	if inx < len(rnv.locations)-1 {
		nextUrl := rnv.locations[inx+1].GetPath()
		n = append(n, NewNode("link", "", "rel", "next", "href", nextUrl))
	}
	lastUrl := rnv.locations[len(rnv.locations)-1].GetPath()
	n = append(n, NewNode("link", "", "rel", "last", "href", lastUrl))
	p.addHeaderNodes(n)
}

func (rnv *ReadNaviComponent) addBodyNodes(p Element) {
	bodyNav := NewNode("nav", "")
	rnv.addFirst(p, bodyNav)
	rnv.addPrevious(p, bodyNav)
	rnv.addNext(p, bodyNav)
	rnv.addLast(p, bodyNav)
	wn := rnv.wrap(bodyNav)
	p.addBodyNodes([]*Node{wn})
}

func (rnv *ReadNaviComponent) visitPage(p Element) {
	if len(rnv.locations) < 3 {
		return
	}
	rnv.addHeaderNodes(p)
	rnv.addBodyNodes(p)

}

func (rnv *ReadNaviComponent) getIndexOfPage(p Element) int {
	for i, l := range rnv.locations {
		if l.GetPath() == p.GetPath() {
			return i
		}
	}
	return -1
}

/* disqus component */

type DisqusComponent struct {
	wrapper
	shortname    string
	configuredJs string
}

func NewDisqusComponent(shortname string) *DisqusComponent {
	d := new(DisqusComponent)
	d.shortname = shortname
	return d
}

func (dc *DisqusComponent) GetCss() string {
	return `
.diqus,
.diqus p {
	font-family: Arial, Helvetica, sans-serif;
}
`
}

func (dc *DisqusComponent) GetJs() string {
	return dc.configuredJs
}

func (dc *DisqusComponent) visitPage(p Element) {
	dc.configuredJs = fmt.Sprintf(
		`
var disqus_config = function () {
	this.page.title= "%s";
	this.page.url = '%s';
	this.page.identifier =  '%s';
};
(function() {
	var d = document, s = d.createElement('script');
	s.src = 'https://%s.disqus.com/embed.js';
	s.setAttribute('data-timestamp', +new Date());
	(d.head || d.body).appendChild(s);
})();
`, p.GetTitle(), p.GetDomain()+p.GetPath(), p.GetDisqusId(), dc.shortname)
	n := NewNode("div", " ", "id", "disqus_thread", "class", "disqus")
	wn := dc.wrap(n)
	p.addBodyNodes([]*Node{wn})
}

/* main  header component */

type MainHeaderComponent struct {
	wrapper
	context Context
}

func NewMainHeaderComponent(context Context) *MainHeaderComponent {
	mhc := new(MainHeaderComponent)
	mhc.context = context
	return mhc
}

func (mhc *MainHeaderComponent) visitPage(p Element) {
	logo := NewNode("a", "<!-- logo -->",
		"href", mhc.context.GetHomeUrl(),
		"class", "headerbar__logo")
	logocontainer := NewNode("div", "",
		"class", "headerbar__logocontainer")
	logocontainer.AddChild(logo)

	header := NewNode("header", "", "class", "headerbar")
	header.AddChild(logocontainer)

	wn := mhc.wrap(header, "headerbar__wrapper")
	p.addBodyNodes([]*Node{wn})
}

func (mhc *MainHeaderComponent) GetJs() string {
	return ""
}

func (mhc *MainHeaderComponent) GetCss() string {
	return `
.headerbar__wrapper {
	position: fixed;
	width: 100%;
	top: 0;
	background-color: white;
}
.headerbar__logo {
	background-image: url(https://s3.amazonaws.com/drewingdeblog/drewing_de_logo.png);
	background-repeat: no-repeat;
	background-position: center center;
	display: block;
	width: 100%;
	height: 80px;
}
.headerbar__navelement {
	display: inline-block;
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	font-weight: 900;
	font-size: 18px;
	line-height: 20px;
	text-transform: uppercase;
	text-decoration: none;
	color: black;
	padding: 10px 20px;
}
`
}

/* content component */
type ContentComponent struct {
	wrapper
}

func NewContentComponent() *ContentComponent {
	return new(ContentComponent)
}

func (cc *ContentComponent) visitPage(p Element) {
	h1 := NewNode("h1", p.GetTitle(),
		"class", "maincontent__h1")
	h2 := NewNode("h2", p.GetPublishedTime(),
		"class", "maincontent__h2")
	n := NewNode("main", p.GetContent(),
		"class", "maincontent")
	n.AddChild(h1)
	n.AddChild(h2)
	wn := cc.wrap(n)
	p.addBodyNodes([]*Node{wn})
}

func (cc *ContentComponent) GetJs() string { return "" }

func (cc *ContentComponent) GetCss() string {
	return `
.maincontent{
	padding-top: 126px;
	text-align: left;
}
.maincontent__h1,
.maincontent__h2 {
	display: inline-block;
	font-family: Arial Black, Arial, Helvetica, sans-serif;
	text-transform: uppercase;
}
.maincontent__h1 ,
.maincontent__h2 {
	font-size: 18px;
	line-height: 20px;
}
.maincontent__h2 {
	color: grey;
	margin-left: 10px;
}
`
}

/* gallery component */

type GalleryComponent struct {
	wrapper
}

func NewGalleryComponent() *GalleryComponent {
	gc := new(GalleryComponent)
	return gc
}

func (gal *GalleryComponent) visitPage(p Element) {
	inner := NewNode("div", "", "class", "maincontent__inner")
	for i := 0; i < 5; i++ {
		title := NewNode("span", "At The Zoo", "class", "portfoliothumb__title")
		subtitle := NewNode("span", "Digital drawing", "class", "portfoliothumb__details")

		label := NewNode("div", "", "class", "portfoliothumb__label")
		label.AddChild(title)
		label.AddChild(subtitle)

		img := NewNode("img", "", "class", "portfoliothumb__image", "src", "https://s3.amazonaws.com/drewingdeblog/blog/wp-content/uploads/2017/12/02152842/atthezoo-400x400.png")

		div := NewNode("a", "", "class", "portfoliothumb", "href", "https://drewing.de")
		div.AddChild(img)
		div.AddChild(label)

		inner.AddChild(div)
	}

	m := NewNode("main", "", "class", "maincontent")
	m.AddChild(inner)
	wn := gal.wrap(m)
	p.addBodyNodes([]*Node{wn})
}

func (gal *GalleryComponent) getCss() string { return `` }

/* copyright component */
type CopyRightComponent struct {
	wrapper
}

func NewCopyRightComponent() *CopyRightComponent {
	return new(CopyRightComponent)
}

func (crc *CopyRightComponent) visitPage(p Element) {
	n := NewNode("div", `
	<a rel="license" class="copyright__cc" href="https://creativecommons.org/licenses/by-nc-nd/3.0/"></a>
	<p class="copyright__license">(c) 2017 by Ingmar Drewing </p><p class="copyright__license">Except where otherwise noted, content on this site is licensed under a <a rel="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">Create Commons Attribution-NonCommercial-NoDerivs 3.0 Unported (CC BY-NC-ND 3.0) license</a>.</p><p class="copyright__license">Soweit nicht anders explizit ausgewiesen, stehen die Inhalte auf dieser Website unter der <a rel="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">Creative Commons Namensnennung-NichtKommerziell-KeineBearbeitung (CC BY-NC-ND 3.0)</a> Lizenz. Unless otherwise noted the author of the content on this page is <a href="https://plus.google.com/113943655600557711368?rel=author">Ingmar Drewing</a></p>
	`, "class", "copyright")
	wn := crc.wrap(n)
	p.addBodyNodes([]*Node{wn})
}

func (crc *CopyRightComponent) GetCss() string {
	return `
.copyright {
	text-align: left;
	font-family: Arial, Helvetica, sans-serif;
	font-size: 14px;
	color: black;
	padding: 20px 0 70px;
}
.copyright__license {
	margin-top: 20px;
	margin-bottom: 20px;
}
.copyright__cc {
    display: block;
    border-width: 0;
    background-image: url(https://i.creativecommons.org/l/by-nc-nd/3.0/88x31.png);
    width: 88px;
    height: 31px;
    margin-right: 15px;
    margin-bottom: 5px;
}
`
}

func (crc *CopyRightComponent) GetJs() string { return `` }

/* cookie notifier component */

type CookieNotifierComponent struct {
}

func (cnc *CookieNotifierComponent) visitPage(p Element) {}

func (cnc *CookieNotifierComponent) getCss() string { return `` }

func (cnc *CookieNotifierComponent) getJs() string {
	return `
function cli_show_cookiebar(p) {
	var Cookie = {
		set: function(name,value,days) {
			if (days) {
				var date = new Date();
				date.setTime(date.getTime()+(days*24*60*60*1000));
				var expires = "; expires="+date.toGMTString();
			}
			else var expires = "";
			document.cookie = name+"="+value+expires+"; path=/";
		},
		read: function(name) {
			var nameEQ = name + "=";
			var ca = document.cookie.split(';');
			for(var i=0;i < ca.length;i++) {
				var c = ca[i];
				while (c.charAt(0)==' ') {
					c = c.substring(1,c.length);
				}
				if (c.indexOf(nameEQ) === 0) {
					return c.substring(nameEQ.length,c.length);
				}
			}
			return null;
		},
		erase: function(name) {
			this.set(name,"",-1);
		},
		exists: function(name) {
			return (this.read(name) !== null);
		}
	};

	var ACCEPT_COOKIE_NAME = 'viewed_cookie_policy',
		ACCEPT_COOKIE_EXPIRE = 365,
		json_payload = p.settings;

	if (typeof JSON.parse !== "function") {
		console.log("CookieLawInfo requires JSON.parse but your browser doesn't support it");
		return;
	}
	var settings = JSON.parse(json_payload);

	var cached_header = jQuery(settings.notify_div_id),
		cached_showagain_tab = jQuery(settings.showagain_div_id),
		btn_accept = jQuery('#cookie_hdr_accept'),
		btn_decline = jQuery('#cookie_hdr_decline'),
		btn_moreinfo = jQuery('#cookie_hdr_moreinfo'),
		btn_settings = jQuery('#cookie_hdr_settings');

	cached_header.hide();
	if ( !settings.showagain_tab ) {
		cached_showagain_tab.hide();
	}

	var hdr_args = { };

	var showagain_args = { };
	cached_header.css( hdr_args );
	cached_showagain_tab.css( showagain_args );

	if (!Cookie.exists(ACCEPT_COOKIE_NAME)) {
		displayHeader();
	}
	else {
		cached_header.hide();
	}

	if ( settings.show_once_yn ) {
		setTimeout(close_header, settings.show_once);
	}
	function close_header() {
		Cookie.set(ACCEPT_COOKIE_NAME, 'yes', ACCEPT_COOKIE_EXPIRE);
		hideHeader();
	}

	var main_button = jQuery('.cli-plugin-main-button');
	main_button.css( 'color', settings.button_1_link_colour );

	if ( settings.button_1_as_button ) {
		main_button.css('background-color', settings.button_1_button_colour);

		main_button.hover(function() {
			jQuery(this).css('background-color', settings.button_1_button_hover);
		},
		function() {
			jQuery(this).css('background-color', settings.button_1_button_colour);
		});
	}
	var main_link = jQuery('.cli-plugin-main-link');
	main_link.css( 'color', settings.button_2_link_colour );

	if ( settings.button_2_as_button ) {
		main_link.css('background-color', settings.button_2_button_colour);

		main_link.hover(function() {
			jQuery(this).css('background-color', settings.button_2_button_hover);
		},
		function() {
			jQuery(this).css('background-color', settings.button_2_button_colour);
		});
	}

	cached_showagain_tab.click(function(e) {
		e.preventDefault();
		cached_showagain_tab.slideUp(settings.animate_speed_hide, function slideShow() {
			cached_header.slideDown(settings.animate_speed_show);
		});
	});

	jQuery("#cookielawinfo-cookie-delete").click(function() {
		Cookie.erase(ACCEPT_COOKIE_NAME);
		return false;
	});
	jQuery("#cookie_action_close_header").click(function(e) {
		e.preventDefault();
		accept_close();
	});

	function accept_close() {
		Cookie.set(ACCEPT_COOKIE_NAME, 'yes', ACCEPT_COOKIE_EXPIRE);

		if (settings.notify_animate_hide) {
			cached_header.slideUp(settings.animate_speed_hide);
		}
		else {
			cached_header.hide();
		}
		cached_showagain_tab.slideDown(settings.animate_speed_show);
		return false;
	}

	function closeOnScroll() {
		if (window.pageYOffset > 100 && !Cookie.read(ACCEPT_COOKIE_NAME)) {
			accept_close();
			if (settings.scroll_close_reload === true) {
				location.reload();
			}
			window.removeEventListener("scroll", closeOnScroll, false);
		}
	}
	if (settings.scroll_close === true) {
		window.addEventListener("scroll", closeOnScroll, false);
	}

	function displayHeader() {
		if (settings.notify_animate_show) {
			cached_header.slideDown(settings.animate_speed_show);
		}
		else {
			cached_header.show();
		}
		cached_showagain_tab.hide();
	}
	function hideHeader() {
		if (settings.notify_animate_show) {
			cached_showagain_tab.slideDown(settings.animate_speed_show);
		}
		else {
			cached_showagain_tab.show();
		}
		cached_header.slideUp(settings.animate_speed_show);
	}
};

function l1hs(str){if(str.charAt(0)=="#"){str=str.substring(1,str.length);}else{return "#"+str;}return l1hs(str);}

cli_show_cookiebar({
					settings: '{"animate_speed_hide":"500","animate_speed_show":"500","background":"#fff","border":"#444","border_on":true,"button_1_button_colour":"#000","button_1_button_hover":"#000000","button_1_link_colour":"#fff","button_1_as_button":true,"button_2_button_colour":"#333","button_2_button_hover":"#292929","button_2_link_colour":"#444","button_2_as_button":false,"font_family":"inherit","header_fix":false,"notify_animate_hide":true,"notify_animate_show":false,"notify_div_id":"#cookie-law-info-bar","notify_position_horizontal":"right","notify_position_vertical":"bottom","scroll_close":false,"scroll_close_reload":false,"showagain_tab":false,"showagain_background":"#fff","showagain_border":"#000","showagain_div_id":"#cookie-law-info-again","showagain_x_position":"100px","text":"#000","show_once_yn":false,"show_once":"10000"}'
});

`
}
