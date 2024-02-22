package health

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/config"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"time"
)

var healthmessage string
var healthy bool
var lastChecked time.Time
var period int

type Msg struct {
	Message   string `json:"message"`
	LastCheck string `json:"lastCheck,omitempty"`
}
type CheckConfig struct {
	Period int
}

// Routes getting all routes for the health endpoint
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/healthiness", GetHealthyEndpoint)
	router.Get("/readiness", GetReadinessEndpoint)
	return router
}

func GetHealthyEndpoint(response http.ResponseWriter, req *http.Request) {
	render.Status(req, http.StatusOK)
	render.JSON(response, req, Msg{
		Message: fmt.Sprintf("service up and running"),
	})
}
func GetReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
	t := time.Now()
	if t.Sub(lastChecked) > (time.Second * time.Duration(2*period)) {
		healthy = false
		healthmessage = "Healthcheck not running"
	}
	if healthy {
		render.Status(req, http.StatusOK)
		render.JSON(response, req, Msg{
			Message:   "service started",
			LastCheck: lastChecked.String(),
		})
	} else {
		render.Status(req, http.StatusServiceUnavailable)
		render.JSON(response, req, Msg{
			Message:   fmt.Sprintf("service is unavailable: %s", healthmessage),
			LastCheck: lastChecked.String(),
		})
	}
}
func doCheck(repository product.Repository) {
	err := repository.Ping()
	if err != nil {
		healthmessage = "failed to ping product repository"
		healthy = false

	} else {
		healthmessage = ""
		healthy = true
	}
	lastChecked = time.Now()
}

func InitHealthSystem(health config.Health, repository product.Repository) {
	period = health.Period
	healthmessage = "service starting"
	healthy = false
	doCheck(repository)
	go func() {
		background := time.NewTicker(time.Second * time.Duration(period))
		for range background.C {
			doCheck(repository)
		}
	}()
}
