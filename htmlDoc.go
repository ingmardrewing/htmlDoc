// A wrapper around qoquery and a means to create markup structures
package htmlDoc

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Create a new Node and takes its name, a possible text contained by
// the Node and a variadic list of attributes as parameters. The attributes
// are supposed to be name value pairs and thus their number need to be
// divisble by two.
func NewNode(tagName string, text string, attributes ...string) *Node {
	if len(attributes)%2 != 0 {
		log.Fatalln("Wrong attribute count at NewNode")
	}
	m := map[string]string{}
	for i := 0; i < len(attributes); i += 2 {
		n := attributes[i]
		v := attributes[i+1]
		m[n] = v
	}
	return &Node{tagName, m, []*Node{}, text}
}

// The Node is a programmatic description of HTML dom nodes
type Node struct {
	tagName  string
	attrMap  map[string]string
	children []*Node
	text     string
}

// Render transforms the node and all its children into HTML
// and returns it as string
func (n *Node) Render() string {
	if n.isEmpty() {
		return n.renderEmpty()
	}
	return n.renderStuffed()
}

// AddChild allow to add a child Node to the current Node
func (n *Node) AddChild(node *Node) {
	n.children = append(n.children, node)
}

func (n *Node) isEmpty() bool {
	return len(n.children) == 0 && n.text == ""
}

func (n *Node) renderAttributes() string {
	attr := []string{}
	for k, v := range n.attrMap {
		attr = append(attr, fmt.Sprintf(`%s="%s"`, k, v))
	}
	sort.Strings(attr)
	return strings.Join(attr, " ")
}

func (n *Node) renderEmpty() string {
	return fmt.Sprintf("<%s %s />", n.tagName, n.renderAttributes())
}

func (n *Node) renderStuffed() string {
	start := fmt.Sprintf("<%s %s>", n.tagName, n.renderAttributes())
	end := fmt.Sprintf("</%s>", n.tagName)
	return start + n.renderChildren() + n.text + end
}

func (n *Node) renderChildren() string {
	html := ""
	for _, child := range n.children {
		html += child.Render()
	}
	return html
}

// NewHtmlDoc createas a pointer to a new HtmlDoc and initializes it
// with an (almost) empty goquery.Document
func NewHtmlDoc() *HtmlDoc {
	p := new(HtmlDoc)
	p.dom, _ = goquery.NewDocumentFromReader(strings.NewReader(""))
	return p
}

// HtmlDoc sort of wraps a goquery document and allows to add
// Nodes to the head and/or body of the document
type HtmlDoc struct {
	head     []*Node
	body     []*Node
	rootAttr map[string]string
	dom      *goquery.Document
}

// Render renders the HtmlDoc as HTML, including all its nodes
// within the head and body part
func (p *HtmlDoc) Render() string {
	dtd := "<!doctype html>"
	p.populateDom()
	html, _ := p.dom.Html()
	parts := strings.Split(html, "<html>")
	return dtd + strings.Join(parts, p.renderRootNode())
}

func (p *HtmlDoc) AddRootAttr(name, value string) {
	p.rootAttr[name] = value
}

func (p *HtmlDoc) renderRootNode() string {
	attrTxts := []string{}
	for k, v := range p.rootAttr {
		att := fmt.Sprintf(`%s="%s"`, k, v)
		attrTxts = append(attrTxts, att)
	}
	attrs := strings.Join(attrTxts, " ")
	if len(attrs) > 0 {
		return fmt.Sprintf("<html %s>", attrs)
	}
	return "<html>"
}

func (p *HtmlDoc) populateDom() {
	for _, m := range p.head {
		p.dom.Find("head").AppendHtml(m.Render())
	}
	for _, m := range p.body {
		p.dom.Find("body").AppendHtml(m.Render())
	}
}

// Add a Node which is going to end up in the head of the HTML Document
func (p *HtmlDoc) AddHeadNode(n *Node) {
	p.head = append(p.head, n)
}

// Add a Node which is going to end up in the body of the HTML Document
func (p *HtmlDoc) AddBodyNode(n *Node) {
	p.body = append(p.body, n)
}
