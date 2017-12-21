package htmlDoc

import (
	"fmt"
)

/* component */
type component interface {
	AddNode(n *Node)
	visitPage(p Element)
}

type concreteComponent struct {
	nodes         []*Node
	visitFunction func(p Element)
}

func (m *concreteComponent) AddNode(n *Node) {
	m.nodes = append(m.nodes, n)
}

func (m *concreteComponent) AddTag(tagName string, text string, attributes ...string) {
	m.AddNode(NewNode(tagName, text, ToMap(attributes...)))
}

func (m *concreteComponent) AddMeta(metaData ...string) {
	n := NewNode("meta", "", ToMap(metaData...))
	m.AddNode(n)
}

func (m *concreteComponent) SetVisitFunction(f func(p Element)) {
	m.visitFunction = f
}

func (m *concreteComponent) visitPage(p Element) {
	m.visitFunction(p)
}

/* HeaderComponent */

type HeaderComponent struct {
	concreteComponent
}

func (hc *HeaderComponent) visitPage(p Element) {
	hc.concreteComponent.visitPage(p)
	p.addHeaderNodes(hc.concreteComponent.nodes)
}

/* BodyComponent */

type BodyComponent struct {
	concreteComponent
}

func (hc *BodyComponent) visitPage(p Element) {
	hc.concreteComponent.visitPage(p)
	p.addBodyNodes(hc.concreteComponent.nodes)
}

/* fb component */
type FBComponent struct {
	concreteComponent
	context Context
}

func NewFBComponent(context Context) *FBComponent {
	fb := new(FBComponent)
	fb.context = context
	return fb
}

func (fbc *FBComponent) visitPage(p Element) {
	fbc.AddMeta("property", "og:title", "content", p.GetTitle())
	fbc.AddMeta("property", "og:url", "content", p.GetUrl())
	fbc.AddMeta("property", "og:image", "content", p.GetImageUrl())
	fbc.AddMeta("property", "og:description", "content", p.GetDescription())
	fbc.AddMeta("property", "og:site_name", "content", fbc.context.GetSiteName())
	fbc.AddMeta("property", "og:type", "content", fbc.context.GetOGType())
	fbc.AddMeta("property", "article:published_time", "content", p.GetPublishedTime())
	fbc.AddMeta("property", "article:modified_time", "content", p.GetPublishedTime())
	fbc.AddMeta("property", "article:section", "content", fbc.context.GetContentSection())
	fbc.AddMeta("property", "article:tag", "content", fbc.context.GetContentTags())
	p.addHeaderNodes(fbc.concreteComponent.nodes)
}

/* google component */

type GoogleComponent struct {
	concreteComponent
	context Context
}

func NewGoogleComponent(context Context) *GoogleComponent {
	gc := new(GoogleComponent)
	gc.context = context
	return gc
}

func (goo *GoogleComponent) visitPage(p Element) {
	goo.AddMeta("itemprop", "name", "content", p.GetTitle())
	goo.AddMeta("itemprop", "description", "content", p.GetDescription())
	goo.AddMeta("itemprop", "image", "content", p.GetImageUrl())
	p.addHeaderNodes(goo.concreteComponent.nodes)
}

/* twitter component */

type TwitterComponent struct {
	concreteComponent
	context Context
}

func NewTwitterComponent(context Context) *TwitterComponent {
	t := new(TwitterComponent)
	t.context = context
	return t
}

func (tw *TwitterComponent) visitPage(p Element) {
	tw.AddMeta("name", "t:card", "content", tw.context.GetTwitterCardType())
	tw.AddMeta("name", "t:site", "content", tw.context.GetTwitterHandle())
	tw.AddMeta("name", "t:title", "content", p.GetTitle())
	tw.AddMeta("name", "t:text:description", "content", p.GetDescription())
	tw.AddMeta("name", "t:creator", "content", tw.context.GetTwitterHandle())
	tw.AddMeta("name", "t:image", "content", p.GetImageUrl())
	p.addHeaderNodes(tw.concreteComponent.nodes)
}

/* title component */
type TitleComponent struct {
	concreteComponent
}

func NewTitleComponent() *TitleComponent {
	return new(TitleComponent)
}

func (tc *TitleComponent) visitPage(p Element) {
	tc.AddTag("title", p.GetTitle())
	p.addHeaderNodes(tc.concreteComponent.nodes)
}

/* css link component */

type CssLinkComponent struct {
	concreteComponent
	url string
}

func NewCssLinkComponent(url string) *CssLinkComponent {
	clc := new(CssLinkComponent)
	clc.url = url
	return clc
}

func (clc *CssLinkComponent) visitPage(p Element) {
	clc.AddTag("link", "", "href", clc.url, "rel", "stylesheet", "type", "text/css")
	p.addHeaderNodes(clc.concreteComponent.nodes)
}

/* naviComponent */

type NaviComponent struct {
	concreteComponent
	locations []Location
}

func NewNaviComponent(locations []Location) *NaviComponent {
	nc := new(NaviComponent)
	nc.locations = locations
	return nc
}

func (nv *NaviComponent) visitPage(p Element) {
	node := NewNode("nav", "", map[string]string{})
	url := p.GetUrl()
	for _, l := range nv.locations {
		if url == l.GetUrl() {
			node.AddChild("span", l.GetTitle())
		} else {
			node.AddChild("a", l.GetTitle(), "href", l.GetUrl())
		}
	}
	nv.concreteComponent.AddNode(node)
	p.addBodyNodes(nv.concreteComponent.nodes)
}

func (nv *NaviComponent) AddLocations(locs []Location) {
	for _, l := range locs {
		nv.locations = append(nv.locations, l)
	}
}

/* ReadNaviComponent */

type ReadNaviComponent struct {
	concreteComponent
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
		n.AddChild("span", "<< first")
	} else {
		f := rnv.locations[0]
		n.AddChild("a", "<< first", "href", f.GetUrl(), "rel", "first")
	}
}

func (rnv *ReadNaviComponent) addPrevious(p Element, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == 0 {
		n.AddChild("span", "< previous")
	} else {
		p := rnv.locations[inx-1]
		n.AddChild("a", "< previous", "href", p.GetUrl(), "rel", "prev")
	}
}

func (rnv *ReadNaviComponent) addNext(p Element, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == len(rnv.locations)-1 {
		n.AddChild("span", "next >")
	} else {
		nx := rnv.locations[inx+1]
		n.AddChild("a", "next >", "href", nx.GetUrl(), "rel", "next")
	}
}

func (rnv *ReadNaviComponent) addLast(p Element, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == len(rnv.locations)-1 {
		n.AddChild("span", "newest >>")
	} else {
		nw := rnv.locations[len(rnv.locations)-1]
		n.AddChild("a", "neweset >>", "href", nw.GetUrl(), "rel", "last")
	}
}

func (rnv *ReadNaviComponent) addHeaderNodes(p Element) {
	inx := rnv.getIndexOfPage(p)
	n := []*Node{}
	firstUrl := rnv.locations[0].GetUrl()
	n = append(n, NewNode("link", "", ToMap("rel", "first", "href", firstUrl)))
	if inx > 0 {
		prevUrl := rnv.locations[inx-1].GetUrl()
		pm := ToMap("rel", "prev", "href", prevUrl)
		n = append(n, NewNode("link", "", pm))
	}
	if inx < len(rnv.locations)-1 {
		nextUrl := rnv.locations[inx+1].GetUrl()
		nm := ToMap("rel", "next", "href", nextUrl)
		n = append(n, NewNode("link", "", nm))
	}
	lastUrl := rnv.locations[len(rnv.locations)-1].GetUrl()
	n = append(n, NewNode("link", "", ToMap("rel", "last", "href", lastUrl)))
	p.addHeaderNodes(n)
}

func (rnv *ReadNaviComponent) addBodyNodes(p Element) {
	bodyNav := NewNode("nav", "", map[string]string{})
	rnv.addFirst(p, bodyNav)
	rnv.addPrevious(p, bodyNav)
	rnv.addNext(p, bodyNav)
	rnv.addLast(p, bodyNav)
	rnv.concreteComponent.AddNode(bodyNav)
	p.addBodyNodes(rnv.concreteComponent.nodes)
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
		if l.GetUrl() == p.GetUrl() {
			return i
		}
	}
	return -1
}

func (rnv *ReadNaviComponent) AddLocations(locs []Location) {
	for _, l := range locs {
		rnv.locations = append(rnv.locations, l)
	}
}

/* disqus component */

type DisqusComponent struct {
	concreteComponent
	shortname string
}

func NewDisqusComponent(shortname string) *DisqusComponent {
	d := new(DisqusComponent)
	d.shortname = shortname
	return d
}

var disqusJS = `
<div id="disqus_thread"></div>
<script>
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
</script>
`

func (dc *DisqusComponent) visitPage(p Element) {
	js := fmt.Sprintf(disqusJS, p.GetTitle(), p.GetUrl(), p.GetDisqusId(), dc.shortname)
	n := NewNode("script", js, ToMap("language", "javascript", "type", "text/javascript"))
	p.addBodyNodes([]*Node{n})
}

/* main  header component */

type MainHeaderComponent struct {
	concreteComponent
	context Context
}

func NewMainHeaderComponent(context Context) *MainHeaderComponent {
	mhc := new(MainHeaderComponent)
	mhc.context = context
	return mhc
}

func (mhc *MainHeaderComponent) visitPage(p Element) {
	header := NewNode("header", "", ToMap("class", "headerbar"))
	inner := header.AddChild("div", "", "class", "headerbar__inner")
	nav := inner.AddChild("nav", "", "class", "headerbar__nav")
	nav.AddChild("a", "twitter", "href", mhc.context.GetTwitterPage(), "class", "headerbar__navelement")
	nav.AddChild("a", "facebook", "href", mhc.context.GetFBPageUrl(), "class", "headerbar__navelement")
	p.addBodyNodes([]*Node{header})
}

/* gallery component */

type GalleryComponent struct {
	concreteComponent
}

func NewGalleryComponent() *GalleryComponent {
	gc := new(GalleryComponent)
	return gc
}

func (gal *GalleryComponent) visitPage(p Element) {
	m := NewNode("main", "", ToMap("class", "maincontent"))
	inner := m.AddChild("div", "", "class", "maincontent__inner")
	for i := 0; i < 5; i++ {
		div := inner.AddChild("a", "", "class", "portfoliothumb", "href", "https://drewing.de")
		div.AddChild("img", "", "class", "portfoliothumb__image", "src", "https://s3.amazonaws.com/drewingdeblog/blog/wp-content/uploads/2017/12/02152842/atthezoo-400x400.png")
		label := div.AddChild("div", "", "class", "portfoliothumb__label")
		label.AddChild("span", "At The Zoo", "class", "portfoliothumb__title")
		label.AddChild("span", "Digital drawing", "class", "portfoliothumb__details")
	}
	p.addBodyNodes([]*Node{m})
}

/* copyright component */
type CopyRightComponent struct {
	concreteComponent
}

func NewCopyRightComponent() *CopyRightComponent {
	return new(CopyRightComponent)
}

func (crc *CopyRightComponent) visitPage(p Element) {
	n := NewNode("div", `
	<a rel="license" class="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">&nbsp;</a>
	<p id="license">© 2017 Ingmar Drewing <br>
Except where otherwise noted, content on this site is licensed under a <a rel="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">Create Commons Attribution-NonCommercial-NoDerivs 3.0 Unported (CC BY-NC-ND 3.0) license</a>.<br>
        Soweit nicht anders explizit ausgewiesen, stehen die Inhalte auf dieser Website unter der <a rel="license" href="https://creativecommons.org/licenses/by-nc-nd/3.0/">Creative Commons Namensnennung-NichtKommerziell-KeineBearbeitung (CC BY-NC-ND 3.0)</a> Lizenz. Unless otherwise noted the author of the content on this page is <a href="https://plus.google.com/113943655600557711368?rel=author">Ingmar Drewing</a>
    </p>
	`, ToMap())
	p.addBodyNodes([]*Node{n})
}

/* cookie notifier component */

type CookieNotifierComponent struct {
}

func (cnc *CookieNotifierComponent) visitPage(p Element) {
	n := NewNode("script", cookiebar, ToMap("language", "javascript", "type", "text/javascript"))
	p.addBodyNodes([]*Node{n})
}

var cookiebar = `
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
