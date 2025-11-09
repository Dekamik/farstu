package shared

type PageViewModel struct {
	ActiveHref string
	NavItems   []NavItemViewModel
	PageTTL    int
	Season     string
	TimeOfDay  string
}

type NavItemViewModel struct {
	Href     string
	Icon     string
	IsActive bool
	Badge    NavItemBadgeViewModel
}

type NavItemBadgeViewModel struct {
	Enabled bool
	Color   string
	Text    string
}

type NewPageViewModelArgs struct {
	ActiveHref               string
	NavItems                 []NavItemViewModel
	Season                   string
	SecondsUntilNextSunEvent int
	TimeOfDay                string
}

func NewPageViewModel(args NewPageViewModelArgs) PageViewModel {
	navModels := make([]NavItemViewModel, 0)

	for _, item := range args.NavItems {
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
