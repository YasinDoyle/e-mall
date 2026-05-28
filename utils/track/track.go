package track

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func GetDefaultConfig() *config.Configuration {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	return cfg
}

func InitTrack() (opentracing.Tracer, io.Closer) {
	cfg := GetDefaultConfig()
	service := "mall"
	cfg.ServiceName = service
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("Error: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func StartSpan(tracer opentracing.Tracer, name string) opentracing.Span {
	return tracer.StartSpan(name)
}

func WithSpan(ctx context.Context, name string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	return span, ctx
}

func StartRequestSpan(ctx context.Context, name string, header http.Header) (opentracing.Span, context.Context, error) {
	tracer := opentracing.GlobalTracer()
	wireContext, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		return nil, ctx, err
	}

	var span opentracing.Span
	if err == nil {
		span = tracer.StartSpan(name, ext.RPCServerOption(wireContext))
	} else {
		span = tracer.StartSpan(name)
	}

	return span, opentracing.ContextWithSpan(ctx, span), nil
}

func GetCarrier(span opentracing.Span) (opentracing.HTTPHeadersCarrier, error) {
	carrier := opentracing.HTTPHeadersCarrier{}
	err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		return nil, err
	}
	return carrier, nil
}

func GetTextMapCarrier(span opentracing.Span) (opentracing.TextMapCarrier, error) {
	carrier := opentracing.TextMapCarrier{}
	err := span.Tracer().Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		return nil, err
	}
	return carrier, nil
}

func GetParentSpan(spanName string, traceId string, header http.Header) (opentracing.Span, error) {
	if traceId != "" {
		header.Set("uber-trace-id", traceId)
	}

	span, _, err := StartRequestSpan(context.Background(), spanName, header)
	return span, err
}

func TraceIDFromSpanContext(spanCtx any) (string, error) {
	ctx, ok := spanCtx.(jaeger.SpanContext)
	if !ok {
		return "", fmt.Errorf("invalid span context type %T", spanCtx)
	}

	return ctx.TraceID().String(), nil
}
