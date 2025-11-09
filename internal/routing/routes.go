package routing

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Dekamik/farstu/internal/components/deviations"
	"github.com/Dekamik/farstu/internal/components/index"
	"github.com/Dekamik/farstu/internal/components/settings"
	"github.com/Dekamik/farstu/internal/components/shared"
	"github.com/Dekamik/farstu/internal/components/sl"
	"github.com/Dekamik/farstu/internal/components/yr"
	"github.com/Dekamik/farstu/internal/config"
)

var navItems = []shared.NavItemViewModel{
	{
		Href: "/",
		Icon: "bi-house-door-fill",
	},
	{
		Href: "/deviations",
		Icon: "bi-cone-striped",
	},
}

type Services struct {
	AppConfig *config.AppConfig
	SL        *sl.SLService
	YR        *yr.YRService
}

func Routes(services Services) {
	slog.Debug("Initializing routes")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := index.Index{}
		shared.ExecuteLayout(w, "internal/components/index/index.html", "/", data)
	})

	http.HandleFunc("/deviations", func(w http.ResponseWriter, r *http.Request) {
		data := deviations.Deviations{}
		shared.ExecuteLayout(w, "internal/components/deviations/deviations.html", "deviations", data)
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		data := settings.Settings{}
		shared.ExecuteLayout(w, "internal/components/settings/settings.html", "settings", data)
	})

	//route("/", func() templ.Component {
	//	seasonAndTimeOfDay, err := yr.GetSeasonAndTimeOfDay(*services.AppConfig)
	//	if err != nil {
	//		slog.Error("an unhandled error occurred", "err", err)
	//		// Try again in an hour
	//		nextRetry := int(time.Now().Local().Add(time.Duration(1) * time.Hour).Sub(time.Now().Local()).Seconds())
	//		seasonAndTimeOfDay = &yr.SeasonAndTimeOfDay{
	//			Season:                   "summer",
	//			SecondsUntilNextSunEvent: nextRetry,
	//			TimeOfDay:                "day",
	//		}
	//	}

	//	args := shared.NewPageViewModelArgs{
	//		ActiveHref:               "/",
	//		NavItems:                 navItems,
	//		SecondsUntilNextSunEvent: seasonAndTimeOfDay.SecondsUntilNextSunEvent,
	//		Season:                   seasonAndTimeOfDay.Season,
	//		TimeOfDay:                seasonAndTimeOfDay.TimeOfDay,
	//	}

	//	yrNowViewModel, yrForecastViewModel := (*services.YR).GetViewModels()
	//	pageViewModel := shared.NewPageViewModel(args)

	//	model := index.IndexViewModel{
	//		Config:     *services.AppConfig,
	//		Departures: (*services.SL).GetDeparturesViewModel(),
	//		Page:       pageViewModel,
	//		Time:       clock.NewViewModel(),
	//		YRNow:      yrNowViewModel,
	//		YRForecast: yrForecastViewModel,
	//	}

	//	return index.IndexView(model)
	//})

	//route("/deviations", func() templ.Component {
	//	seasonAndTimeOfDay, err := yr.GetSeasonAndTimeOfDay(*services.AppConfig)
	//	if err != nil {
	//		slog.Error("an unhandled error occurred", "err", err)
	//		// Try again in an hour
	//		nextRetry := int(time.Now().Local().Add(time.Duration(1) * time.Hour).Sub(time.Now().Local()).Seconds())
	//		seasonAndTimeOfDay = &yr.SeasonAndTimeOfDay{
	//			Season:                   "summer",
	//			SecondsUntilNextSunEvent: nextRetry,
	//			TimeOfDay:                "day",
	//		}
	//	}

	//	args := shared.NewPageViewModelArgs{
	//		ActiveHref:               "/deviations",
	//		NavItems:                 navItems,
	//		SecondsUntilNextSunEvent: seasonAndTimeOfDay.SecondsUntilNextSunEvent,
	//		Season:                   seasonAndTimeOfDay.Season,
	//		TimeOfDay:                seasonAndTimeOfDay.TimeOfDay,
	//	}

	//	deviationsViewModel := (*services.SL).GetDeviationsViewModel()

	//	model := sl.DeviationsViewModel{
	//		Deviations: deviationsViewModel.Deviations,
	//		Page:       shared.NewPageViewModel(args),
	//	}

	//	return sl.DeviationsView(model)
	//})
}

func HTMXRoutes(services Services) {
	http.HandleFunc("/htmx/time", func(w http.ResponseWriter, r *http.Request) {
		data := shared.GetTime()
		tmpl := template.Must(template.ParseFiles("internal/components/shared/layout/clock.html"))
		tmpl.Execute(w, data)
	})

	//route("/htmx/departures", func() templ.Component {
	//	return sl.DeparturesView((*services.SL).GetDeparturesViewModel())
	//})

	//route("/htmx/departures/next", func() templ.Component {
	//	return sl.SLDeparturesNextView((*services.SL).GetDeparturesViewModel())
	//})

	//route("/htmx/departures/list", func() templ.Component {
	//	return sl.SLDeparturesListView((*services.SL).GetDeparturesViewModel())
	//})

	//route("/htmx/time", func() templ.Component {
	//	return clock.View(clock.NewViewModel())
	//})

	//route("/htmx/yrnow", func() templ.Component {
	//	yrNowViewModel, _ := (*services.YR).GetViewModels()
	//	return yr.YRNowView(yrNowViewModel)
	//})

	//route("/htmx/yrforecast", func() templ.Component {
	//	_, yrForecastViewModel := (*services.YR).GetViewModels()
	//	return yr.YRForecastView(yrForecastViewModel)
	//})
}
