package api

import (
	"github.com/gorilla/mux"
	"github.com/igridnet/users"
	"net/http"
)

type (
	Client struct {
		Users *users.Client
	}
)

func (c *Client)MakeHandler()http.Handler{
	r := mux.NewRouter()
	r.HandleFunc("/admins",c.registerAdmin).Methods(http.MethodPost)
	r.HandleFunc("/login",c.adminLogin).Methods(http.MethodGet)
	return r
}

func (c *Client) registerAdmin(writer http.ResponseWriter, request *http.Request) {

}

func (c *Client) adminLogin(writer http.ResponseWriter, request *http.Request) {

}


