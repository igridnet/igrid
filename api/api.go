package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/igridnet/users/api"
	"github.com/techcraftlabs/base"
	"github.com/techcraftlabs/base/io"
	"net/http"
	"time"
)

type (
	Client struct {
		rv    base.Receiver
		rp    base.Replier
		Users *api.Client
	}

	Empty struct {
	}

	LoginResponse struct {
		Token string `json:"token,omitempty"`
	}
)

func NewClient(client *api.Client) *Client {
	rp := base.NewReplier(io.Stderr, true)
	rv := base.NewReceiver(io.Stderr, true)

	return &Client{
		rv:    rv,
		rp:    rp,
		Users: client,
	}
}

func (c *Client) MakeHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/admins", c.registerAdmin).Methods(http.MethodPost)
	r.HandleFunc("/login", c.adminLogin).Methods(http.MethodGet)
	r.HandleFunc("/nodes",c.addNode).Methods(http.MethodPost)
	r.HandleFunc("/nodes",c.getNodes).Methods(http.MethodGet)
	r.HandleFunc("/nodes/{id}",c.getNodeById).Methods(http.MethodGet)
	r.HandleFunc("/regions",c.getRegions).Methods(http.MethodGet)
	r.HandleFunc("/regions",c.addRegions).Methods(http.MethodPost)
	r.HandleFunc("/regions/{id}",c.getRegionById).Methods(http.MethodGet)


	return r
}

func (c *Client) registerAdmin(writer http.ResponseWriter, request *http.Request) {

}

func (c *Client) adminLogin(writer http.ResponseWriter, request *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	receipt, err := c.rv.Receive(ctx, "login", request, &Empty{})
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	basicAuth := receipt.BasicAuth

	token, err := c.Users.Login(ctx, basicAuth.Username, basicAuth.Password)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	responsePayload := LoginResponse{Token: token}
	response := base.NewResponse(200, responsePayload)
	c.rp.Reply(writer, response)
}

func (c *Client) getNodes(writer http.ResponseWriter, request *http.Request) {

}

func (c *Client) getNodeById(writer http.ResponseWriter, request *http.Request) {

}

func (c *Client) getRegions(writer http.ResponseWriter, request *http.Request) {

}

func (c *Client) getRegionById(writer http.ResponseWriter, request *http.Request) {

}

func (c *Client) addNode(writer http.ResponseWriter, request *http.Request) {
	
}

func (c *Client) addRegions(writer http.ResponseWriter, request *http.Request) {
	
}
