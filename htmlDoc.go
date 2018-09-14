// A wrapper around qoquery and a means to create markup structures
package htmlDoc

import (
	"fmt"
	"log"
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

// AddChild allow to add a child Node to the current Node
func (n *Node) AddChild(node *Node) {
	n.children = append(n.children, node)
}

func (n *Node) isEmpty() bool {
	return len(n.children) == 0 && n.text == ""
}

// NewHtmlDoc createas a pointer to a new HtmlDoc and initializes it
// with an (almost) empty goquery.Document
func NewHtmlDoc() *HtmlDoc {
	p := new(HtmlDoc)
	return p
}

// HtmlDoc sort of wraps a goquery document and allows to add
// Nodes to the head and/or body of the document
type HtmlDoc struct {
	head     []*Node
	body     []*Node
	rootAttr []string
}

// Render renders the HtmlDoc as HTML, including all its nodes
// within the head and body part
func (p *HtmlDoc) Render() string {
	renderer := NewHtmlDocRenderer(p)
	return renderer.render()
}

// Render as AMP
func (p *HtmlDoc) RenderAmp() string {
	renderer := NewAmpDocRenderer(p)
	return renderer.render()
}

func (p *HtmlDoc) AddRootAttr(att ...string) {
	if len(att) == 1 {
		p.rootAttr = append(p.rootAttr, att[0])
	} else if len(att) == 2 {
		attr := fmt.Sprintf(`%s="%s"`, att[0], att[1])
		p.rootAttr = append(p.rootAttr, attr)
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
