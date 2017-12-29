package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/mikaelm1/blog_app/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
