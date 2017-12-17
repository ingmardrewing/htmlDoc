package htmlDoc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
	Node
*/
type Node struct {
	tagName           string
	attributeValueMap map[string]string
	children          []*Node
	text              string
}

func NewNode(tagName string, text string, attributeValueMap map[string]string) *Node {
	return &Node{tagName, attributeValueMap, []*Node{}, text}
}

func (n *Node) isEmpty() bool {
	return len(n.children) == 0 && n.text == ""
}

func (n *Node) Render() string {
	if n.isEmpty() {
		return n.renderEmpty()
	}
	return n.renderStuffed()
}

func (n *Node) renderEmpty() string {
	tag := "<" + n.tagName
	keys := getSortedMapKeys(n.attributeValueMap)
	for _, k := range keys {
		tag += fmt.Sprintf(` %s="%s"`, k, n.attributeValueMap[k])
	}
	tag += " />"
	return tag
}

func (n *Node) renderStuffed() string {
	tag := "<" + n.tagName
	keys := getSortedMapKeys(n.attributeValueMap)
	for _, k := range keys {
		tag += fmt.Sprintf(` %s="%s"`, k, n.attributeValueMap[k])
	}
	tag += ">"
	for _, child := range n.children {
		tag += child.Render()
	}
	tag += n.text
	tag += "</" + n.tagName + ">"
	return tag
}

func (n *Node) addChildAsNode(child *Node) {
	n.children = append(n.children, child)
}

/*
	AddChild
*/
func (n *Node) AddChild(tagName string, text string, attributes ...string) *Node {
	child := NewNode(tagName, text, ToMap(attributes...))
	n.children = append(n.children, child)
	return child
}

func getSortedMapKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

/*
	HtmlDoc
*/

type HtmlDoc struct {
	title   string
	head    []*Node
	content []*Node
	dom     *goquery.Document
}

func NewHtmlDoc() *HtmlDoc {
	p := new(HtmlDoc)
	p.dom, _ = goquery.NewDocumentFromReader(strings.NewReader(""))
	return p
}

func (p *HtmlDoc) Render() string {
	dtd := "<!doctype html>"
	p.populateDom()
	html, _ := p.dom.Html()
	return dtd + html
}

func (p *HtmlDoc) populateDom() {
	for _, m := range p.head {
		p.dom.Find("head").AppendHtml(m.Render())
	}
	for _, m := range p.content {
		p.dom.Find("body").AppendHtml(m.Render())
	}
}

func (p *HtmlDoc) AddMeta(attributes ...string) {
	p.AddHeadTag("meta", "", attributes...)
}

func (p *HtmlDoc) AddHeadTag(tagName string, text string, attributes ...string) *Node {
	n := p.CreateLonelyNode(tagName, text, attributes...)
	p.AddHeadNode(n)
	return n
}

func (p *HtmlDoc) AddHeadNode(n *Node) {
	p.head = append(p.head, n)
}

func (p *HtmlDoc) AddBodyNode(n *Node) {
	p.content = append(p.content, n)
}

func (p *HtmlDoc) AddContentTag(tagName string, text string, attributes ...string) *Node {
	n := p.CreateLonelyNode(tagName, text, attributes...)
	p.content = append(p.content, n)
	return n
}

func (p *HtmlDoc) AddContentAsNode(n *Node) {
	p.content = append(p.content, n)
}

func (p *HtmlDoc) CreateLonelyNode(tagName string, text string, attributes ...string) *Node {
	return NewNode(tagName, text, ToMap(attributes...))
}

func (p *HtmlDoc) AddTitle(title string) {
	p.title = title
}

func ToMap(namesAndValues ...string) map[string]string {
	if len(namesAndValues)%2 != 0 {
		panic("Wrong parameter count")
	}
	m := map[string]string{}
	for i := 0; i < len(namesAndValues); i += 2 {
		n := namesAndValues[i]
		v := namesAndValues[i+1]
		m[n] = v
	}
	return m
}

/* location */

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

func (p *Page) GetUrl() string {
	return p.url
}

func (p *Page) GetTitle() string {
	return p.title
}

func (p *Page) acceptVisitor(v component) {
	v.visit(p)
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

/* component */

type component interface {
	AddNode(n *Node)
	visit(p *Page)
}

type concreteComponent struct {
	nodes         []*Node
	visitFunction func(p *Page)
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

func (m *concreteComponent) SetVisitFunction(f func(p *Page)) {
	m.visitFunction = f
}

func (m *concreteComponent) visit(p *Page) {
	m.visitFunction(p)
}

/* HeaderComponent */

type HeaderComponent struct {
	concreteComponent
}

func (hc *HeaderComponent) visit(p *Page) {
	hc.concreteComponent.visit(p)
	p.addHeaderNodes(hc.concreteComponent.nodes)
}

/* BodyComponent */

type BodyComponent struct {
	concreteComponent
}

func (hc *BodyComponent) visit(p *Page) {
	hc.concreteComponent.visit(p)
	p.addBodyNodes(hc.concreteComponent.nodes)
}

/* naviComponent */

type NaviComponent struct {
	concreteComponent
	locations []Location
}

func (nv *NaviComponent) visit(p *Page) {
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

func (rnv *ReadNaviComponent) addFirst(p *Page, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == 0 {
		n.AddChild("span", "<< first")
	} else {
		f := rnv.locations[0]
		n.AddChild("a", "<< first", "href", f.GetUrl())
	}
}

func (rnv *ReadNaviComponent) addPrevious(p *Page, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == 0 {
		n.AddChild("span", "< previous")
	} else {
		p := rnv.locations[inx-1]
		n.AddChild("a", "< previous", "href", p.GetUrl())
	}
}

func (rnv *ReadNaviComponent) addNext(p *Page, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == len(rnv.locations)-1 {
		n.AddChild("span", "next >")
	} else {
		nx := rnv.locations[inx+1]
		n.AddChild("a", "next >", "href", nx.GetUrl())
	}
}

func (rnv *ReadNaviComponent) addLast(p *Page, n *Node) {
	inx := rnv.getIndexOfPage(p)
	if inx == len(rnv.locations)-1 {
		n.AddChild("span", "newest >>")
	} else {
		nw := rnv.locations[len(rnv.locations)-1]
		n.AddChild("a", "neweset >>", "href", nw.GetUrl())
	}
}

func (rnv *ReadNaviComponent) visit(p *Page) {
	if len(rnv.locations) < 3 {
		return
	}
	node := NewNode("nav", "", map[string]string{})
	rnv.addFirst(p, node)
	rnv.addPrevious(p, node)
	rnv.addNext(p, node)
	rnv.addLast(p, node)
	rnv.concreteComponent.AddNode(node)
	p.addBodyNodes(rnv.concreteComponent.nodes)
}

func (rnv *ReadNaviComponent) getIndexOfPage(p *Page) int {
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
