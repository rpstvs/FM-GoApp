package routes

import (
	"github.com/go-chi/chi"
	"github.com/rpstvs/fm-goapp/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkById)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)

	return r
}
