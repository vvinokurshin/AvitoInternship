package router

import (
	"github.com/gorilla/mux"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	historyDelivery "github.com/vvinokurshin/AvitoInternship/internal/history/delivery"
	segmentDelivery "github.com/vvinokurshin/AvitoInternship/internal/segment/delivery"
	userDelivery "github.com/vvinokurshin/AvitoInternship/internal/user/delivery"
	"net/http"
)

func AddRoutes(r *mux.Router, cfg *config.Config, userD userDelivery.DeliveryI, segmentD segmentDelivery.DeliveryI,
	historyD historyDelivery.DeliveryI) {
	// User
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserCreate, userD.CreateUser).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUser, userD.EditUser).Methods(http.MethodPut)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUser, userD.DeleteUser).Methods(http.MethodDelete)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUser, userD.GetUser).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserSegments, segmentD.GetUserSegments).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserEditSegments, segmentD.EditUserSegments).Methods(http.MethodPut)

	// Segment
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSegmentCreate, segmentD.CreateSegment).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSegment, segmentD.DeleteSegment).Methods(http.MethodDelete)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSegment, segmentD.GetSegment).Methods(http.MethodGet)

	// History
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteGetHistory, historyD.GetHistoryCSV).Methods(http.MethodGet)
}
