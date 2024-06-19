package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SetNameInterface interface{}

func (appHandler *AppHandler) SetName(mux chi.Router, db SetNameInterface) {
	mux.Post("/name", func(w http.ResponseWriter, r *http.Request) {

	})
}

type FirstCompanyInterface interface{}

func (appHandler *AppHandler) FirstCompany(mux chi.Router, db FirstCompanyInterface) {
	mux.Post("/company", func(w http.ResponseWriter, r *http.Request) {

	})
}

type EndOfOnboardingInterface interface{}

func (appHandler *AppHandler) EndOfOnboarding(mux chi.Router, db EndOfOnboardingInterface) {
	mux.Post("/end", func(w http.ResponseWriter, r *http.Request) {

	})
}
