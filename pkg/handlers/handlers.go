package handlers

import (
	"net/http"

	"github.com/Gideon-isa/bookings/pkg/config"
	"github.com/Gideon-isa/bookings/pkg/models"
	"github.com/Gideon-isa/bookings/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Respository

// Repository is the repository type
type Respository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Respository {
	return &Respository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Respository) {
	Repo = r
}

// About is the about handler
func (m *Respository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page
func (m *Respository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	//store the remoteIp in the stringMap
	stringMap["remote_ip"] = remoteIP

	// Send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
