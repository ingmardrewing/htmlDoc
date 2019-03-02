package htmlDoc

import "fmt"

type ampDocRenderer struct {
	docRenderer
	doc Doc
}

var ampTemplate = `<!doctype html>
<html amp lang="en">
  <head>
    <meta charset="utf-8">
    <script async src="https://cdn.ampproject.org/v0.js"></script>
	%s
    <title>Hello, AMPs</title>
    <link rel="canonical" href="http://example.ampproject.org/article-metadata.html">
    <meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">
    <script type="application/ld+json">
      {
        "@context": "http://schema.org",
        "@type": "NewsArticle",
        "headline": "Open-source framework for publishing content",
        "datePublished": "2015-10-07T12:02:41Z",
        "image": [
          "logo.jpg"
        ]
      }
    </script>
    <style amp-boilerplate>body{-webkit-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-moz-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-ms-animation:-amp-start 8s steps(1,end) 0s 1 normal both;animation:-amp-start 8s steps(1,end) 0s 1 normal both}@-webkit-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-moz-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-ms-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-o-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}</style><noscript><style amp-boilerplate>body{-webkit-animation:none;-moz-animation:none;-ms-animation:none;animation:none}</style></noscript>
  </head>
  <body>
    %s
  </body>
</html>`

func NewAmpDocRenderer(d *HtmlDoc) *ampDocRenderer {
	a := new(ampDocRenderer)
	a.doc = d
	a.nodeRendererProvider = NewAmpNodeRenderer
	return a
}

func (a *ampDocRenderer) render() string {
	return fmt.Sprintf(ampTemplate,
		a.renderSliceOfNodes(a.doc.headNodes()),
		a.renderSliceOfNodes(a.doc.bodyNodes()))
}

type ampNodeRenderer struct {
	n *Node
}

func NewAmpNodeRenderer(n *Node) nodeRenderer {
	r := new(ampNodeRenderer)
	r.n = n
	return r
}

func (a *ampNodeRenderer) render() string {
	// TODO: Implement
	return "<amp>"
}
