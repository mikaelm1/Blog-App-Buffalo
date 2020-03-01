package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/mikaelm1/Blog-App-Buffalo/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
