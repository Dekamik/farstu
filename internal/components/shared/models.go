package shared

import (
	"github.com/a-h/templ"
)

type PageViewModel struct {
	ActiveHref string
	NavItems   []NavItemViewModel
	PageTTL    int
	Season     string
	TimeOfDay  string
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
}

type NewPageViewModelArgs struct {
	ActiveHref               string
	Season                   string
	SecondsUntilNextSunEvent int
	TimeOfDay                string
}

func NewPageViewModel(args NewPageViewModelArgs) PageViewModel {
	navModels := make([]NavItemViewModel, 0)

	for _, item := range navItems {
		model := NavItemViewModel{
			Href:     item.Href,
			Icon:     item.Icon,
			IsActive: string(item.Href) == args.ActiveHref,
		}
		navModels = append(navModels, model)
	}

	return PageViewModel{
		ActiveHref: args.ActiveHref,
		NavItems:   navModels,
		PageTTL:    args.SecondsUntilNextSunEvent,
		Season:     args.Season,
		TimeOfDay:  args.TimeOfDay,
	}
}

type ModalViewModel struct {
	ID    string
	Title string
}
