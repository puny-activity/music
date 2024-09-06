package httprouter

import (
	"bytes"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/httpresp"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

type Wrapper struct {
	writer       *httpresp.Writer
	errorStorage *errs.Storage
	log          *zerolog.Logger
}

func NewWrapper(w *httpresp.Writer, errorStorage *errs.Storage, log *zerolog.Logger) *Wrapper {
	return &Wrapper{
		writer:       w,
		errorStorage: errorStorage,
		log:          log,
	}
}

func (w *Wrapper) Wrap(controllerFunction func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()

		ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

		var requestBody bytes.Buffer
		if request.Body != nil {
			_, err := io.Copy(&requestBody, request.Body)
			if err != nil {
				w.log.Error().Err(err).Msg("failed to read request body")
				return
			}
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.log.Error().Err(err).Msg("failed to close request body")
			}
		}(request.Body)

		clonedRequest := request.Clone(request.Context())
		clonedRequest.Body = io.NopCloser(bytes.NewReader(requestBody.Bytes()))

		ctx := request.Context()
		clonedRequest = clonedRequest.WithContext(ctx)

		controllerError := controllerFunction(ww, clonedRequest)
		if controllerError != nil {
			statusCode := 400
			if errors.Is(controllerError, errs.Unexpected) {
				statusCode = 500
			}
			err := w.writer.Write(
				ww,
				statusCode,
				httpresp.Error{
					Code: w.errorStorage.Code(controllerError),
				})
			if err != nil {
				w.log.Debug().Err(err).Msg("failed to write error response")
			}
		}

		requestDuration := time.Since(startTime)

		if ww.Status() >= http.StatusInternalServerError {
			event := w.log.Error().Err(controllerError)
			logBase(event, request, requestDuration, ww)
			event.Msgf("request handled with unexpected error: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		} else if ww.Status() >= http.StatusBadRequest {
			event := w.log.Warn().Err(controllerError)
			logBase(event, request, requestDuration, ww)
			event.Msgf("request handled with error: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		} else if ww.Status() >= http.StatusOK {
			event := w.log.Info()
			logBase(event, request, requestDuration, ww)
			event.Msgf("request handled succussfully: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		}
	}
}

func logBase(event *zerolog.Event, request *http.Request, duration time.Duration, ww middleware.WrapResponseWriter) {
	event.Str("method", request.Method).
		Str("path", request.URL.Path).
		Str("duration", duration.String()).
		Int("status", ww.Status()).
		Str("user_agent", request.UserAgent()).
		Str("source_ip", request.RemoteAddr)
}
