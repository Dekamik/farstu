package shared

import "github.com/a-h/templ"

type PageModel struct {
	NavItems []NavItemModel
}

type NavItemModel struct {
	Href     templ.SafeURL
	Icon     string
	IsActive bool
}
