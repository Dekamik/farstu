package shared

import "github.com/a-h/templ"

type PageViewModel struct {
	NavItems []NavItemViewModel
}

type NavItemViewModel struct {
	Href     templ.SafeURL
	Icon     string
	IsActive bool
}

var navItems = []NavItemViewModel{
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

func NewPageViewModel(activeHref string) PageViewModel {
	navModels := make([]NavItemViewModel, 0)

	for _, item := range navItems {
		model := NavItemViewModel{
			Href:     item.Href,
			Icon:     item.Icon,
			IsActive: string(item.Href) == activeHref,
		}
		navModels = append(navModels, model)
	}

	return PageViewModel{
		NavItems: navModels,
	}
}
