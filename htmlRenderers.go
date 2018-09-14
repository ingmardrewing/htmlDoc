package htmlDoc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type nodeRenderer interface {
	render() string
}

type htmlDocRenderer struct {
	docRenderer
	doc Doc
}

func NewHtmlDocRenderer(d Doc) *htmlDocRenderer {
	h := new(htmlDocRenderer)
	h.doc = d
	h.nodeRendererProvider = NewHtmlNodeRenderer
	return h
}

func (h *htmlDocRenderer) render() string {
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(""))
	dom.Find("head").AppendHtml(h.renderSliceOfNodes(h.doc.headNodes()))
	dom.Find("body").AppendHtml(h.renderSliceOfNodes(h.doc.bodyNodes()))

	html, _ := dom.Html()
	parts := strings.Split(html, "<html>")
	dtd := "<!doctype html>"
	return dtd + strings.Join(parts, h.renderRootNode())
}

func (h *htmlDocRenderer) renderRootNode() string {
	attrs := strings.Join(h.doc.rootAttributes(), " ")
	if len(attrs) > 0 {
		return fmt.Sprintf("<html %s>", attrs)
	}
	return "<html>"
}

//
type htmlNodeRenderer struct {
	n *Node
}

func NewHtmlNodeRenderer(n *Node) nodeRenderer {
	h := new(htmlNodeRenderer)
	h.n = n
	return h
}

func (h *htmlNodeRenderer) render() string {
	if h.n.isEmpty() {
		return fmt.Sprintf("<%s%s />", h.n.tagName, h.renderAttributes())
	}
	return h.renderStuffed()
}

func (h *htmlNodeRenderer) renderStuffed() string {
	start := fmt.Sprintf("<%s%s>", h.n.tagName, h.renderAttributes())
	end := fmt.Sprintf("</%s>", h.n.tagName)
	return start + h.renderChildren() + h.n.text + end
}

func (h *htmlNodeRenderer) renderAttributes() string {
	attr := []string{}
	for k, v := range h.n.attrMap {
		attr = append(attr, fmt.Sprintf(` %s="%s"`, k, v))
	}
	sort.Strings(attr)
	return strings.Join(attr, "")
}

func (h *htmlNodeRenderer) renderChildren() string {
	html := ""
	for _, child := range h.n.children {
		childRenderer := NewHtmlNodeRenderer(child)
		html += childRenderer.render()
	}
	return html
}
