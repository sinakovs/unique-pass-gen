package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"unique-pass-gen/pkg/generator"
	"unique-pass-gen/pkg/passwordstore"
)

type Handler struct {
	Store passwordstore.PasswordStore
}

func NewHandler(store passwordstore.PasswordStore) *Handler {
	return &Handler{Store: store}
}

var baseTmpl = template.Must(template.ParseFiles("templates/base.page.tmpl"))
var resultTmpl = template.Must(template.ParseFiles("templates/result.page.tmpl"))

func (h *Handler) GetForm(w http.ResponseWriter, r *http.Request) {
	if err := baseTmpl.Execute(w, nil); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func (h *Handler) GeneratePass(w http.ResponseWriter, r *http.Request) {
	opts, err := parseOptions(r)

	generator := generator.NewGenerator(h.Store)

	if err != nil {
		renderError(w, baseTmpl, err.Error())
		return
	}

	password, err := generator.UniquePasswordGenerator(opts)
	if err != nil {
		renderError(w, baseTmpl, err.Error())
		return
	}

	if err := resultTmpl.Execute(w, map[string]string{"Password": password}); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func renderError(w http.ResponseWriter, tmpl *template.Template, msg string) {
	if err := tmpl.Execute(w, map[string]string{"Error": msg}); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func parseOptions(r *http.Request) (generator.Options, error) {
	opts := []generator.Option{}

	if err := r.ParseForm(); err != nil {
		return generator.Options{}, fmt.Errorf("invalid form")
	}

	length, err := strconv.Atoi(r.FormValue("length"))
	if err != nil {
		return generator.Options{}, fmt.Errorf("invalid length")
	}

	opts = append(opts, generator.WithLength(length))

	sets := r.Form["sets"]
	if len(sets) == 0 {
		return generator.Options{}, fmt.Errorf("select at least one set")
	}

	for _, s := range sets {
		switch s {
		case "digits":
			opts = append(opts, generator.WithDigits())
		case "lower":
			opts = append(opts, generator.WithLowerC())
		case "upper":
			opts = append(opts, generator.WithUpperC())
		}
	}

	data := generator.NewOptions(opts...)

	return data, nil
}
