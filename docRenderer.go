package htmlDoc

type docRenderer struct {
	nodeRendererProvider func(*Node) nodeRenderer
}

func (h *docRenderer) renderSliceOfNodes(son []*Node) string {
	html := ""
	for _, n := range son {
		renderer := h.nodeRendererProvider(n)
		html += renderer.render()
	}
	return html
}
