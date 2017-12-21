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

func NewNode(tagName string, text string, attributes ...string) *Node {
	return &Node{tagName, ToMap(attributes...), []*Node{}, text}
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

func (n *Node) AddChild(node *Node) {
	n.children = append(n.children, node)
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
	m := NewNode("meta", "", attributes...)
	p.AddHeadNode(m)
}

func (p *HtmlDoc) AddHeadNode(n *Node) {
	p.head = append(p.head, n)
}

func (p *HtmlDoc) AddBodyNode(n *Node) {
	p.content = append(p.content, n)
}

/* utils */

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

func getSortedMapKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
