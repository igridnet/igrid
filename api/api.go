package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/igridnet/users/api"
	"github.com/igridnet/users/models"
	"github.com/techcraftlabs/base"
	"github.com/techcraftlabs/base/io"
	"net/http"
	"strings"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req := new(models.AdminRegReq)

	_, err := c.rv.Receive(ctx,"admin login",request,req)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	invalid := (strings.TrimSpace(req.Email) == "") || (strings.TrimSpace(req.Password) == "") || (strings.TrimSpace(req.Name) == "")

	if invalid{
		http.Error(writer, "bad request specify email, name and password", http.StatusBadRequest)
		return
	}
	admin, err := c.Users.Register(ctx,*req)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	headersOption := base.WithResponseHeaders(map[string]string{
		"Content-Type":"application/json",
	})
	response := base.NewResponse(200,admin,headersOption)
	c.rp.Reply(writer,response)
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
	nodes, err := c.Users.ListNodes(context.Background())
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}

	headersOption := base.WithResponseHeaders(map[string]string{
		"Content-Type":"application/json",
	})
	response := base.NewResponse(200,nodes,headersOption)
	c.rp.Reply(writer,response)
}

func (c *Client) getNodeById(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	node, err := c.Users.GetNode(context.Background(),id)
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}
	headersOption := base.WithResponseHeaders(map[string]string{
		"Content-Type":"application/json",
	})
	response := base.NewResponse(200,node,headersOption)
	c.rp.Reply(writer,response)
}

func (c *Client) getRegions(writer http.ResponseWriter, request *http.Request) {
	regions, err := c.Users.ListRegions(context.Background())
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}

	headersOption := base.WithResponseHeaders(map[string]string{
		"Content-Type":"application/json",
	})
	response := base.NewResponse(200,regions,headersOption)
	c.rp.Reply(writer,response)
}

func (c *Client) getRegionById(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	region, err := c.Users.GetRegion(context.Background(),id)
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}
	headersOption := base.WithResponseHeaders(map[string]string{
		"Content-Type":"application/json",
	})
	response := base.NewResponse(200,region,headersOption)
	c.rp.Reply(writer,response)
}

func (c *Client) addNode(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(),time.Minute)
	defer cancel()
	req := new(models.NodeRegReq)
	_, err := c.rv.Receive(ctx, "add node",request,req)
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}
	addedNode, err := c.Users.AddNode(ctx, *req)
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}
	headersOption := base.WithResponseHeaders(map[string]string{
		"Content-Type":"application/json",
	})
	response := base.NewResponse(200,addedNode,headersOption)
	c.rp.Reply(writer,response)
}

func (c *Client) addRegions(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(),time.Minute)
	defer cancel()
	req := new(models.RegionRegReq)
	_, err := c.rv.Receive(ctx, "add node",request,req)
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}
	addedRegion,err := c.Users.AddRegion(ctx, *req)
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}
	headersOption := base.WithResponseHeaders(map[string]string{
		"Content-Type":"application/json",
	})
	response := base.NewResponse(200,addedRegion,headersOption)
	c.rp.Reply(writer,response)
}
