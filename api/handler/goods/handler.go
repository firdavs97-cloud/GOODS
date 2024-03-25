package goods

import (
	"goods/api/httputils"
	goodModel "goods/models/good"
	"log"
	"net/http"
	"strconv"
)

type ListResponseMeta struct {
	goodModel.Report
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListResponse struct {
	Meta  ListResponseMeta `json:"meta"`
	Goods []goodModel.Good `json:"goods"`
}

func List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offsetStr := queryParams.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 1
	}

	meta, list, err := goodModel.ListRecordsAndReport(limit, offset)
	if err != nil || meta == nil {
		http.Error(w, "Could not query items", http.StatusNotFound)
		log.Println("db error:", err)
		return
	}

	// Respond with updated resource
	httputils.ResponseBody(w, ListResponse{
		Meta: ListResponseMeta{
			Report: *meta,
			Limit:  limit,
			Offset: offset,
		},
		Goods: list,
	}, http.StatusOK)
}
