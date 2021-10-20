package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/igridnet/igrid/internal"
	"github.com/igridnet/users"
	"net/http"
	"time"
)

type (
	Client struct {
		rv internal.Receiver
		rp internal.Replier
		Users *users.Client
	}

	Empty struct {

	}

	LoginResponse struct {
		Token string `json:"token,omitempty"`
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

	ctx, cancel := context.WithTimeout(context.Background(),time.Minute)
	defer cancel()
	receipt, err := c.rv.Receive(ctx,"login", request,&Empty{})
	if err != nil {
		http.Error(writer,err.Error(),500)
		return
	}
	basicAuth := receipt.BasicAuth

	token, err := c.Users.Login(ctx, basicAuth.Username, basicAuth.Password)
	if err != nil {
		http.Error(writer,err.Error(),500)
		return
	}

	responsePayload := LoginResponse{Token: token}
	response := internal.NewResponse(200,responsePayload)
	c.rp.Reply(writer, response)
}


