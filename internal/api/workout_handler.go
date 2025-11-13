package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rpstvs/fm-goapp/internal/middleware"
	"github.com/rpstvs/fm-goapp/internal/store"
	"github.com/rpstvs/fm-goapp/internal/utils"
)

type WorkoutHanlder struct {
	workoutStore store.WorkoutStore
	Logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHanlder {
	return &WorkoutHanlder{
		workoutStore: workoutStore,
		Logger:       logger,
	}
}

func (wh *WorkoutHanlder) HandleGetWorkById(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParams(r)

	if err != nil {
		wh.Logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
	}

	workout, err := wh.workoutStore.GetWorkoutById(workoutID)

	if err != nil {
		wh.Logger.Printf("ERROR: GetWorkoutbyId: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHanlder) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		return
	}

	workout.UserID = currentUser.ID
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHanlder) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")

	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutID)

	if err != nil {
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}

	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}

	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}

	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}

	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(int64(currentUser.ID))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		return
	}

	if workoutOwner != currentUser.ID {
		return
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}

func (wh *WorkoutHanlder) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")

	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(int64(currentUser.ID))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		return
	}

	if workoutOwner != currentUser.ID {
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutID)

	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error deleting workout", http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
