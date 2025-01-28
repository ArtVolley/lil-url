package save

import (
	"errors"
	resp "lil-url/internal/lib/api/response"
	"lil-url/internal/lib/logger/sl"
	"lil-url/internal/lib/random"
	"lil-url/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Url    string `jsonP:"url" validate:"required,url"`
	LilUrl string `json:"lilUrl,omitempty"`
}

type Response struct {
	resp.Response
	LilUrl string `json:"lilUrl,omitempty"`
}

const urlLilLength = 5

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type UrlSaver interface {
	SaveUrl(urlToSave, lil string) (int64, error)
}

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.save.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err = validator.New().Struct(req)
		if err != nil {
			validatorErrs := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationErrors(validatorErrs))

			return
		}

		lilUrl := req.LilUrl
		if lilUrl == "" {
			lilUrl = random.NewRandomString(urlLilLength)
		}

		id, err := urlSaver.SaveUrl(req.Url, req.LilUrl)
		if errors.Is(err, storage.ErrUrlExists) {

			log.Info("url already exists", slog.String("url", req.Url))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Info("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id), slog.String("lilUrl", lilUrl))

		var res Response
		res.Response = resp.Ok()
		res.LilUrl = lilUrl
		render.JSON(w, r, res)
	}
}
