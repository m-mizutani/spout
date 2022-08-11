package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/usecase"
	"github.com/m-mizutani/spout/pkg/utils"
)

type Server struct {
	mux *chi.Mux
}

func New(uc *usecase.Usecase) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	hdlr := handler{
		uc: uc,
	}

	r.Route("/api", func(r chi.Router) {
		r.Route("/logs", func(r chi.Router) {
			r.Get("/", hdlr.serve(getLogs))
		})
	})

	return &Server{
		mux: r,
	}
}

func (x *Server) Listen(addr string) error {
	if err := http.ListenAndServe(addr, x.mux); err != nil {
		return goerr.Wrap(err, "http server error")
	}
	return nil
}

type httpResponse struct {
	code int
	data any
}

type httpErrorMessage struct {
	Error string `json:"error"`
}

type apiHandler func(ctx *model.Context, uc *usecase.Usecase, r *http.Request) (*httpResponse, error)

type handler struct {
	uc *usecase.Usecase
}

func (x *handler) serve(f apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base := r.Context()
		var ctx *model.Context
		if c, ok := base.(*model.Context); ok {
			ctx = c
		} else {
			ctx = model.NewContext(model.WithCtx(base))
		}

		resp, err := f(ctx, x.uc, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := httpErrorMessage{
				Error: err.Error(),
			}
			if err := json.NewEncoder(w).Encode(msg); err != nil {
				utils.Logger.Err(err).Error("marshal http error message")
			}
			return
		}

		w.WriteHeader(resp.code)
		if err := json.NewEncoder(w).Encode(resp.data); err != nil {
			utils.Logger.Err(err).Error("marshal http response message")
		}
	}
}
