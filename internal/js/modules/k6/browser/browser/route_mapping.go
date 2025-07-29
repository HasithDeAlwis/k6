package browser

import (
	"context"

	"github.com/grafana/sobek"

	"go.k6.io/k6/internal/js/modules/k6/browser/common"
	"go.k6.io/k6/internal/js/modules/k6/browser/k6ext"
)

// mapRoute to the JS module.
func mapRoute(vu moduleVU, route *common.Route) mapping {
	return mapping{
		"abort": func(reason string) *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, route.Abort(reason)
			})
		},
		"continue": func(opts sobek.Value) *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				copts := parseContinueOptions(vu.Context(), opts)
				return nil, route.Continue(copts)
			})
		},
		"fulfill": func(opts sobek.Value) *sobek.Promise {
			fopts := parseFulfillOptions(vu.Context(), opts)
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, route.Fulfill(fopts)
			})
		},
		"request": func() mapping {
			return mapRequest(vu, route.Request())
		},
	}
}

func parseContinueOptions(ctx context.Context, opts sobek.Value) *common.ContinueOptions {
	if !sobekValueExists(opts) {
		return nil
	}

	rt := k6ext.Runtime(ctx)
	copts := &common.ContinueOptions{}

	obj := opts.ToObject(rt)
	for _, k := range obj.Keys() {
		switch k {
		case "headers":
			copts.Headers = parseHeaders(obj.Get(k).ToObject(rt))
		case "method":
			copts.Method = obj.Get(k).String()
		case "postData":
			copts.PostData = obj.Get(k).String()
		case "url":
			copts.URL = obj.Get(k).String()
		}
	}

	return copts
}

func parseFulfillOptions(ctx context.Context, opts sobek.Value) common.FulfillOptions {
	fopts := common.FulfillOptions{}
	if !sobekValueExists(opts) {
		return fopts
	}

	rt := k6ext.Runtime(ctx)
	obj := opts.ToObject(rt)
	for _, k := range obj.Keys() {
		switch k {
		case "body":
			fopts.Body = obj.Get(k).String()
		case "contentType":
			fopts.ContentType = obj.Get(k).String()
		case "headers":
			fopts.Headers = parseHeaders(obj.Get(k).ToObject(rt))
		case "status":
			fopts.Status = obj.Get(k).ToInteger()
		}
	}

	return fopts
}

func parseHeaders(headers *sobek.Object) []common.HTTPHeader {
	headersKeys := headers.Keys()
	result := make([]common.HTTPHeader, len(headersKeys))
	for i, hk := range headersKeys {
		result[i] = common.HTTPHeader{
			Name:  hk,
			Value: headers.Get(hk).String(),
		}
	}
	return result
}
