package htmlDoc

/* component */

type visitor interface {
	visitPage(p Element)
}

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

/* naviComponent */

type NaviComponent struct {
	concreteComponent
	locations []Location
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

/* disqus */

type disqusComponent struct {
	concreteComponent
}
