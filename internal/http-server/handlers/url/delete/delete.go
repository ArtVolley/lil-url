package delete

import (
	"errors"
	resp "lil-url/internal/lib/api/response"
	"lil-url/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	LilUrl string `json:"lilUrl"`
}

type UrlDeleter interface {
	DeleteUrl(lil string) error
}

func New(log *slog.Logger, urlDeleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.redirect.New"

		log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		lilUrl := chi.URLParam(r, "lilUrl")
		if lilUrl == "" {
			log.Info("lilUrl is empty")

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err := urlDeleter.DeleteUrl(lilUrl)
		if errors.Is(err, storage.ErrLilNotFound) {
			log.Info("url not found", slog.String("lilUrl", lilUrl))

			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", slog.String("lilUrl", lilUrl))

			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("deleted lilUrl", slog.String("lilUrl", lilUrl))

		w.WriteHeader(http.StatusNoContent)
	}
}
