package actions

import (
	"html/template"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/microcosm-cc/bluemonday"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public/assets")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			"csrf": func() template.HTML {
				return template.HTML("<input name=\"authenticity_token\" value=\"<%= authenticity_token %>\" type=\"hidden\">")
			},
			"markdown2": func(content string) template.HTML {
				output := blackfriday.Run([]byte(content), blackfriday.WithExtensions(blackfriday.CommonExtensions))
				html := bluemonday.UGCPolicy().SanitizeBytes(output)
				return template.HTML(html)
			},
		},
	})
}
