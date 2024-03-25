package good

import (
	"goods/api/httputils"
	goodModel "goods/models/good"
	"log"
	"net/http"
	"strconv"
)

type CreateRequest struct {
	Name string `json:"name" validate:"required"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()
	projectIdStr := queryParams.Get("projectId")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		http.Error(w, "Invalid projectId", http.StatusBadRequest)
		return
	}

	// Parse JSON body
	var body CreateRequest
	httputils.ParseBody(w, r, &body)

	// Create resource using body data and projectId
	good, err := goodModel.New(body.Name, projectId)
	if err != nil {
		http.Error(w, "Failed to create resource", http.StatusInternalServerError)
		log.Println("Error creating resource:", err)
		return
	}

	// Respond with created resource
	httputils.ResponseBody(w, good, http.StatusCreated)
}

type UpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func Update(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()
	idStr := queryParams.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	projectIdStr := queryParams.Get("projectId")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		http.Error(w, "Invalid projectId", http.StatusBadRequest)
		return
	}

	// Parse JSON body
	var body UpdateRequest
	httputils.ParseBody(w, r, &body)

	old, err := goodModel.Get(id, projectId)
	if err != nil {
		http.Error(w, "Could not find item", http.StatusNotFound)
		log.Println("not found:", err)
		return
	}
	good := *old

	good.Name = body.Name
	if body.Description != "" {
		good.Description = body.Description
	}

	err = good.Save(old)
	if err != nil {
		http.Error(w, "Could not update item", http.StatusInternalServerError)
		log.Println("db error:", err)
		return
	}

	// Respond with updated resource
	httputils.ResponseBody(w, good, http.StatusOK)
}

type RemoveResponse struct {
	ID         int  `json:"id"`
	CampaignId int  `json:"campaignId"`
	Removed    bool `json:"removed"`
}

func Remove(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()
	idStr := queryParams.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	projectIdStr := queryParams.Get("projectId")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		http.Error(w, "Invalid projectId", http.StatusBadRequest)
		return
	}

	old, err := goodModel.Get(id, projectId)
	if err != nil {
		http.Error(w, "Could not find item", http.StatusNotFound)
		log.Println("not found:", err)
		return
	}

	err = goodModel.Remove(id, projectId, old)
	if err != nil {
		http.Error(w, "Could not remove item", http.StatusInternalServerError)
		log.Println("db error:", err)
		return
	}

	// Respond with updated resource
	httputils.ResponseBody(w, RemoveResponse{
		ID:         id,
		CampaignId: projectId,
		Removed:    true,
	}, http.StatusOK)
}

type ReprioritiizeRequest struct {
	NewPriority int `json:"newPriority" validate:"required"`
}

func Reprioritiize(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()
	idStr := queryParams.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	projectIdStr := queryParams.Get("projectId")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		http.Error(w, "Invalid projectId", http.StatusBadRequest)
		return
	}

	var body ReprioritiizeRequest
	httputils.ParseBody(w, r, &body)

	g, err := goodModel.Get(id, projectId)
	if err != nil {
		http.Error(w, "Could not find item", http.StatusNotFound)
		log.Println("not found:", err)
		return
	}

	res, err := g.Reprioritiize(body.NewPriority)
	if err != nil {
		http.Error(w, "failed to change db", http.StatusInternalServerError)
		log.Println("db error:", err)
		return
	}

	// Respond with updated resource
	httputils.ResponseBody(w, res, http.StatusOK)
}
