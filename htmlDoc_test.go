package htmlDoc

import "testing"

func TestNewPage(t *testing.T) {
	page := NewHtmlDoc()
	actual := page.Render()

	expected := `<!doctype html><html><head></head><body></body></html>`
	if actual != expected {
		t.Fatal("Expected ", expected, " but got ", actual)
	}
}

func TestNodeAddChild(t *testing.T) {
	n := NewNode("nav", "", "class", "mainNavi")
	a := NewNode("a", "test", "href", "test.html")
	n.AddChild(a)

	actual := n.Render()
	expected := `<nav class="mainNavi"><a href="test.html">test</a></nav>`

	if actual != expected {
		t.Fatal("Expected ", expected, " but got ", actual)
	}
}

func TestNewPageWithcontent(t *testing.T) {
	page := NewHtmlDoc()
	page.AddMeta("name", "wurst", "value", "mett")

	a1 := NewNode("a", "1", "href", "page1.html")
	a2 := NewNode("a", "2", "href", "page2.html")

	nav := NewNode("nav", "", "class", "mainNav")
	page.AddBodyNode(nav)
	nav.AddChild(a1)
	nav.AddChild(a2)

	h1 := NewNode("h1", "WTF")

	header := NewNode("header", "")
	page.AddBodyNode(header)
	header.AddChild(h1)

	p := NewNode("p", "Test")

	main := NewNode("main", "")
	page.AddBodyNode(main)
	main.AddChild(p)

	actual := page.Render()
	expected := `<!doctype html><html><head><meta name="wurst" value="mett"/></head><body><nav class="mainNav"><a href="page1.html">1</a><a href="page2.html">2</a></nav><header><h1>WTF</h1></header><main><p>Test</p></main></body></html>`

	if actual != expected {
		t.Fatal("Expected ", expected, " but got ", actual)
	}
}

func TestAddMeta(t *testing.T) {
	page := NewHtmlDoc()
	page.AddMeta("name", "testname", "value", "testvalue")
	actual := page.Render()
	expected := `<!doctype html><html><head><meta name="testname" value="testvalue"/></head><body></body></html>`
	if expected != actual {
		t.Fatal("Expected ", expected, " but got ", actual)
	}
}

func TestAddNestedContentTags(t *testing.T) {
	page := NewHtmlDoc()

	p := NewNode("p", "")
	page.AddBodyNode(p)

	a := NewNode("a", "label", "href", "test")
	p.AddChild(a)

	actual := page.Render()
	expected := `<!doctype html><html><head></head><body><p><a href="test">label</a></p></body></html>`

	if expected != actual {
		t.Fatal("Expected ", expected, " but got ", actual)
	}
}

func TestNodeIsEmpty(t *testing.T) {
	a := NewNode("a", "", "href", "test")

	actual := a.isEmpty()
	expected := true

	if expected != actual {
		t.Fatal("Expected ", expected, " but got ", actual)
	}

	a = NewNode("a", "wurst", "href", "test")

	actual = a.isEmpty()
	expected = false

	if expected != actual {
		t.Fatal("Expected ", expected, " but got ", actual)
	}
}
