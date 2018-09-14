package htmlDoc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type htmlDocRenderer struct {
	doc *HtmlDoc
}

func NewHtmlDocRenderer(d *HtmlDoc) *htmlDocRenderer {
	h := new(htmlDocRenderer)
	h.doc = d
	return h
}

func (h *htmlDocRenderer) render() string {
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(""))
	for _, n := range h.doc.head {
		renderer := NewHtmlNodeRenderer(n)
		dom.Find("head").AppendHtml(renderer.render())
	}
	for _, n := range h.doc.body {
		renderer := NewHtmlNodeRenderer(n)
		dom.Find("body").AppendHtml(renderer.render())
	}

	html, _ := dom.Html()
	parts := strings.Split(html, "<html>")
	dtd := "<!doctype html>"
	return dtd + strings.Join(parts, h.renderRootNode())
}

func (h *htmlDocRenderer) renderRootNode() string {
	attrs := strings.Join(h.doc.rootAttr, " ")
	if len(attrs) > 0 {
		return fmt.Sprintf("<html %s>", attrs)
	}
	return "<html>"
}

type htmlNodeRenderer struct {
	n *Node
}

func NewHtmlNodeRenderer(n *Node) *htmlNodeRenderer {
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
