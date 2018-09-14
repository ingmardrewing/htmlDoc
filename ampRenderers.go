package htmlDoc

type ampDocRenderer struct {
	doc *HtmlDoc
}

func NewAmpDocRenderer(d *HtmlDoc) *ampDocRenderer {
	a := new(ampDocRenderer)
	a.doc = d
	return a
}

// TODO: Implement
func (a *ampDocRenderer) render() string {
	return ""
}
