package htmlDoc

/* component */

type visitor interface {
	visitPage(p *Page)
}

type component interface {
	AddNode(n *Node)
	visitPage(p *Page)
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
