package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rpstvs/fm-goapp/internal/api"
	"github.com/rpstvs/fm-goapp/internal/store"
	"github.com/rpstvs/fm-goapp/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHanlder
	UserHandler    *api.UserHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFs(pgDB, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//stores
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)

	//handlers
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		DB:             pgDB,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK\n")
	w.WriteHeader(200)

}
