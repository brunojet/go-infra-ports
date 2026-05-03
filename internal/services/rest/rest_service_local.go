package rest

import "net/http"

const (
	statusInfo = http.StatusContinue
	statusOK   = http.StatusOK
	status3xx  = http.StatusMultipleChoices
	status4xx  = http.StatusBadRequest
)

type restNon2xxMapper func(statusCode int, upsPayload RestResponseSpec, serviceMeta *ServiceMeta) error

func (s *restService[C, R, U]) mapUpstreamContext(reqCtx RequestContext, upsCtx *RequestContext) error {
	if err := s.upstream.ToUpstreamQuery(reqCtx.Query, upsCtx.Query); err != nil {
		return errRestServiceUpstreamMappingFailed("ToUpstreamQuery", err)
	}
	if err := s.upstream.ToUpstreamHeaders(reqCtx.Headers, upsCtx.Headers); err != nil {
		return errRestServiceUpstreamMappingFailed("ToUpstreamHeaders", err)
	}
	upsCtx.Identifiers = reqCtx.Identifiers
	return nil
}

func (s *restService[C, R, U]) mapDownstreamContext(upsCtx ResponseContext, respCtx *ResponseContext) error {
	if err := s.downstream.ToDownstreamStatusCode(upsCtx.StatusCode, &respCtx.StatusCode); err != nil {
		return errRestServiceDownstreamMappingFailed("ToDownstreamStatusCode", err)
	}
	if respCtx.Headers == nil {
		respCtx.Headers = http.Header{}
	}
	if err := s.downstream.ToDownstreamHeaders(upsCtx.Headers, respCtx.Headers); err != nil {
		return errRestServiceDownstreamMappingFailed("ToDownstreamHeaders", err)
	}
	return nil
}

func (s *restService[C, R, U]) non2xxMapper(statusCode int) (restNon2xxMapper, string, error) {
	switch {
	case statusCode >= statusInfo && statusCode <= statusInfo+99:
		return s.downstream.ToDownstreamInformation, "ToDownstreamInformation", nil
	case statusCode >= status3xx && statusCode <= status3xx+99:
		return s.downstream.ToDownstreamRedirection, "ToDownstreamRedirection", nil
	case statusCode >= status4xx && statusCode <= status4xx+199:
		return s.downstream.ToDownstreamProblem, "ToDownstreamProblem", nil
	default:
		return nil, "", errRestServiceInvalidNon2xxStatus
	}
}

//nolint:gocritic // Value parameter is intentional to avoid pointer-only API for read-only payload selection.
func non2xxPayload(statusCode int, restResp RestResponse) (RestResponseSpec, error) {
	switch {
	case statusCode >= statusInfo && statusCode <= statusInfo+99:
		return restResp.Information, nil
	case statusCode >= status3xx && statusCode <= status3xx+99:
		return restResp.Redirection, nil
	case statusCode >= status4xx && statusCode <= status4xx+199:
		return restResp.Problem, nil
	default:
		return nil, errRestServiceInvalidNon2xxStatus
	}
}

//nolint:gocritic // Value parameter is intentional to keep list mapping path pointer-free.
func non2xxPayloadFromResponses(statusCode int, restResp RestResponses) (RestResponseSpec, error) {
	switch {
	case statusCode >= statusInfo && statusCode <= statusInfo+99:
		return restResp.Information, nil
	case statusCode >= status3xx && statusCode <= status3xx+99:
		return restResp.Redirection, nil
	case statusCode >= status4xx && statusCode <= status4xx+199:
		return restResp.Problem, nil
	default:
		return nil, errRestServiceInvalidNon2xxStatus
	}
}

func (s *restService[C, R, U]) mapNon2xxResponse(statusCode int, payload RestResponseSpec, meta *ServiceMeta) error {
	mapper, mapperName, err := s.non2xxMapper(statusCode)
	if err != nil {
		return err
	}
	if err := mapper(statusCode, payload, meta); err != nil {
		return errRestServiceDownstreamMappingFailed(mapperName, err)
	}

	return nil
}

// mapRestResponseToServiceResponse maps a RestResponse from repository layer to the individual
// response fields: data, meta, and response context.
// It handles 2xx (data), 3xx (redirection), 1xx (information), and 4xx/5xx (problem) responses.
//
//nolint:gocritic // Value parameter is intentional; caller owns a local copy and no mutation is needed.
func (s *restService[C, R, U]) mapRestResponseToServiceResponse(restResp RestResponse, data *R, meta *ServiceMeta) error {
	statusCode := restResp.Context.StatusCode
	if statusCode == http.StatusNoContent {
		return nil
	}
	if statusCode >= statusOK && statusCode <= statusOK+99 {
		if restResp.Data == nil {
			return errRestServiceNilResponseData
		}
		if err := s.downstream.ToDownstreamResponse(restResp.Data, data); err != nil {
			return errRestServiceDownstreamMappingFailed("ToDownstreamResponse", err)
		}
		return nil
	}
	payload, err := non2xxPayload(statusCode, restResp)
	if err != nil {
		return err
	}
	return s.mapNon2xxResponse(statusCode, payload, meta)
}

// mapRestResponsesToServiceResponses maps a RestResponses from repository layer to data slice and meta.
// It handles 2xx (data slice) and non-2xx (information, redirection, problem) responses.
//
//nolint:gocritic // Value parameter is intentional; no mutation and simpler nil semantics for caller.
func (s *restService[C, R, U]) mapRestResponsesToServiceResponses(restResp RestResponses, data *[]R, meta *ServiceMeta) error {
	statusCode := restResp.Context.StatusCode
	if statusCode >= statusOK && statusCode <= statusOK+99 {
		if restResp.Data == nil {
			return errRestServiceNilResponseData
		}
		if data == nil {
			return errRestServiceResponseNil
		}
		*data = make([]R, len(restResp.Data))
		for i, restSpec := range restResp.Data {
			if err := s.downstream.ToDownstreamResponse(restSpec, &(*data)[i]); err != nil {
				return errRestServiceDownstreamMappingFailed("ToDownstreamResponse[slice]", err)
			}
		}
		if err := s.downstream.ToDownstreamResponseMeta(restResp.Context.Meta, meta); err != nil {
			return errRestServiceDownstreamMappingFailed("ToDownstreamResponseMeta", err)
		}
		return nil
	}
	payload, err := non2xxPayloadFromResponses(statusCode, restResp)
	if err != nil {
		return err
	}
	return s.mapNon2xxResponse(statusCode, payload, meta)
}
