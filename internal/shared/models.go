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

var navItems = []NavItemModel{
	{
		Href: "/",
		Icon: "bi-house-door-fill",
	},
	{
		Href: "/weather",
		Icon: "bi-cloud-sun-fill",
	},
	{
		Href: "/disruptions",
		Icon: "bi-exclamation-triangle-fill",
	},
}

func NewPageModel(activeHref string) PageModel {
	navModels := make([]NavItemModel, 0)

	for _, item := range navItems {
		model := NavItemModel{
			Href:     item.Href,
			Icon:     item.Icon,
			IsActive: string(item.Href) == activeHref,
		}
		navModels = append(navModels, model)
	}

	return PageModel{
		NavItems: navModels,
	}
}
