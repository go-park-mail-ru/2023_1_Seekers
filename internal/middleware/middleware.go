package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgHttp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/http"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/logger"
	promMetrics "github.com/go-park-mail-ru/2023_1_Seekers/pkg/metrics/prometheus"
	pkgErrors "github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
)

type HttpMiddleware struct {
	cfg    *config.Config
	sUC    auth.UseCaseI
	log    *logger.Logger
	metric *promMetrics.MetricsHttp
}

type GRPCMiddleware struct {
	cfg    *config.Config
	log    *logger.Logger
	metric *promMetrics.MetricsGRPC
}

func NewHttpMiddleware(c *config.Config, sUC auth.UseCaseI, l *logger.Logger, m *promMetrics.MetricsHttp) *HttpMiddleware {
	return &HttpMiddleware{c, sUC, l, m}
}

func NewGRPCMiddleware(c *config.Config, l *logger.Logger, m *promMetrics.MetricsGRPC) *GRPCMiddleware {
	return &GRPCMiddleware{c, l, m}
}

func (m *HttpMiddleware) Cors(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedMethods:   m.cfg.Cors.AllowedMethods,
		AllowedOrigins:   m.cfg.Cors.AllowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   m.cfg.Cors.AllowedHeaders,
	})
	return c.Handler(h)
}

func (m *HttpMiddleware) HandlerLogger(h http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerLogger := m.log.LoggerWithFields(map[string]any{
			"method": r.Method,
			"url":    r.URL.Path,
		})
		handlerLogger.Info("new request")
		r = r.WithContext(context.WithValue(r.Context(), common.ContextHandlerLog, handlerLogger))
		h.ServeHTTP(w, r)
	})
	return handler
}

func (m *HttpMiddleware) CheckAuth(h http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(m.cfg.Sessions.CookieName)
		if err != nil {
			pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
			return
		}
		session, err := m.sUC.GetSession(cookie.Value)
		if err != nil {
			pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), common.ContextUser, session.UID))

		h.ServeHTTP(w, r)
	})
	return handler
}

func (m *HttpMiddleware) CheckCSRF(h http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//cookie, err := r.Cookie(m.cfg.Sessions.CookieName)
		//if err != nil {
		//	pkgHttp.HandleError(w, r, pkgErrors.Wrap(errors.ErrFailedAuth, err.Error()))
		//	return
		//}
		//csrfToken := r.Header.Get(m.cfg.Sessions.CSRFHeader)
		//if csrfToken == "" {
		//	pkgHttp.HandleError(w, r, pkgErrors.WithMessage(errors.ErrWrongCSRF, "token not presented"))
		//	return
		//}
		//
		//err = crypto.CheckCSRF(cookie.Value, csrfToken)
		//if err != nil {
		//	pkgHttp.HandleError(w, r, pkgErrors.Wrap(err, "failed check csrf"))
		//	return
		//}
		h.ServeHTTP(w, r)
	})
	return handler
}

func (m *HttpMiddleware) MetricsHttp(h http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rWriterWithCode := negroni.NewResponseWriter(w)

		start := time.Now()
		h.ServeHTTP(rWriterWithCode, r)

		code := rWriterWithCode.Status()

		m.metric.Timings.WithLabelValues(strconv.Itoa(code), r.URL.String(), r.Method).Observe(time.Since(start).Seconds())
		m.metric.Hits.Inc()

		if code >= http.StatusMultipleChoices {
			m.metric.Errors.WithLabelValues(
				strconv.Itoa(code),
				r.URL.String(),
				r.Method,
			).Inc()
		}
	})
	return handler
}

func (m *GRPCMiddleware) MetricsGRPCUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, uHandler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := uHandler(ctx, req)

	errStatus, _ := status.FromError(err)
	code := errStatus.Code()

	if code != codes.OK {
		m.metric.Errors.WithLabelValues(code.String(), info.FullMethod).Inc()
	}

	m.metric.Timings.WithLabelValues(code.String(), info.FullMethod).Observe(time.Since(start).Seconds())
	m.metric.Hits.Inc()

	return resp, err
}

func (m *GRPCMiddleware) LoggerGRPCUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, uHandler grpc.UnaryHandler) (interface{}, error) {
	// TODO request ID to context and log
	grpcLogger := m.log.LoggerWithFields(logrus.Fields{
		"method": info.FullMethod,
	})

	ctx = context.WithValue(ctx, common.ContextHandlerLog, grpcLogger)
	resp, err := uHandler(ctx, req)

	return resp, err
}
