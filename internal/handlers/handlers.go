package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stephenmontague/go-bnb/internal/driver"
	"github.com/stephenmontague/go-bnb/internal/helpers"
	"github.com/stephenmontague/go-bnb/internal/repository"
	"github.com/stephenmontague/go-bnb/internal/repository/dbrepo"
	"net/http"

	"github.com/stephenmontague/go-bnb/internal/config"
	"github.com/stephenmontague/go-bnb/internal/forms"
	"github.com/stephenmontague/go-bnb/internal/models"
	"github.com/stephenmontague/go-bnb/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository {
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Next 2 lines of code stores the remote IP in the session with a key to look it up
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation (w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string] interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation (w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation {
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Phone: r.Form.Get("phone"),
		Email: r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Warlocks renders the warlocks room page
func (m *Repository) Warlocks (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "warlocks.page.tmpl", &models.TemplateData{})
}

// Warriors renders the warriors room page
func (m *Repository) Warriors (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "warriors.page.tmpl", &models.TemplateData{})
}

// Availability renders the availability page
func (m *Repository) Availability (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the availability page
func (m *Repository) PostAvailability (w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK bool `json: "ok"`
	Message string `json: "message"`
}

// AvailabilityJSON handles the reqeust for availability and send JSON response
func (m *Repository) AvailabilityJSON (w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse {
		OK: true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}


func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}