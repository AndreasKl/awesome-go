package frontend

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/rs/zerolog/log"
)

//go:embed templates/login.gohtml
var loginTemplate string

type LoginController struct {
	template *template.Template
}

func NewLoginController() (*LoginController, error) {
	t, err := template.New("login").Parse(loginTemplate)
	return &LoginController{template: t}, err
}

func (l *LoginController) RenderForm(resp http.ResponseWriter, req *http.Request) {
	err := l.template.Execute(resp, nil)
	if err != nil {
		log.Err(err).Msg("Unable to render login form template.")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (l *LoginController) HandleForm(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Err(err).Msg("Unable to parse form data.")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Info().Msg(fmt.Sprintf("username: %+v", req.Form["username"]))
	log.Info().Msg(fmt.Sprintf("password: %+v", req.Form["password"]))

	// TODO: Validate credentials, build a cookie session, redirect...
}
