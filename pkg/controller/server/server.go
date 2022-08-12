package server

import (
	"encoding/json"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/frontend"
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

	if err := importStaticFiles(r); err != nil {
		panic("failed to import static file" + err.Error())
	}

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

func importStaticFiles(r chi.Router) error {
	assets := frontend.Assets()
	if err := fs.WalkDir(assets, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return goerr.Wrap(err)
		}
		if d.IsDir() {
			return nil
		}

		path := strings.TrimPrefix(filePath, "out")
		if path == "" {
			return nil
		}

		body, err := assets.ReadFile(filePath)
		if err != nil {
			return goerr.Wrap(err)
		}
		mimeType := mime.TypeByExtension(filepath.Ext(filePath))

		utils.Logger.With("path", filePath).With("type", mimeType).Trace("set static file")

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", mimeType)
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(body); err != nil {
				utils.Logger.Err(err).Error("failed to response static file")
			}
		}

		r.Get(path, handler)
		if strings.HasSuffix(path, "/index.html") {
			r.Get(strings.TrimSuffix(path, "index.html"), handler)
		} else if strings.HasSuffix(path, ".html") {
			// e.g. out/verify/email.html
			indexPath := strings.TrimSuffix(path, ".html")
			r.Get(indexPath, handler)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
