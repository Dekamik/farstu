package routing

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Dekamik/farstu/internal/config"
	"github.com/Dekamik/farstu/internal/routes/deviations"
	"github.com/Dekamik/farstu/internal/routes/index"
	"github.com/Dekamik/farstu/internal/routes/index/components/sl"
	"github.com/Dekamik/farstu/internal/routes/index/components/yr"
	"github.com/Dekamik/farstu/internal/routes/settings"
	"github.com/Dekamik/farstu/internal/routes/shared"
)

type Services struct {
	AppConfig *config.AppConfig
	SL        *sl.SLService
	YR        *yr.YRService
}

func Routes(services Services) {
	slog.Debug("Initializing routes")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := index.Index{
			Departures: (*services.SL).GetDepartures(),
			Forecast: (*services.YR).GetForecast(*services.AppConfig),
		}
		templates := []string{
			"internal/routes/index/index.html",
			"internal/routes/index/components/sl/sl.html",
			"internal/routes/index/components/yr/yr.html",
		}
		shared.ExecuteLayout(w, "/", data, templates...)
	})

	http.HandleFunc("/deviations", func(w http.ResponseWriter, r *http.Request) {
		data := deviations.Deviations{}
		shared.ExecuteLayout(w, "deviations", data, "internal/routes/deviations/deviations.html")
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		data := settings.Settings{}
		shared.ExecuteLayout(w, "settings", data, "internal/routes/settings/settings.html")
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
		tmpl := template.Must(template.ParseFiles("internal/routes/shared/layout/clock.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			slog.Error(err.Error())
		}
	})

	http.HandleFunc("/htmx/sl", func(w http.ResponseWriter, r *http.Request) {
		data := index.Index{
			Departures: (*services.SL).GetDepartures(),
		}
		tmpl := template.Must(template.ParseFiles("internal/routes/index/components/sl/sl.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			slog.Error(err.Error())
		}
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
