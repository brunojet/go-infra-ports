package rest

import (
	"context"
)

// restService implements Service[C, R, U] by mapping HTTP REST requests/responses through configurable mappers.
type restService[C, R, U any] struct {
	repo       RestRepository
	upstream   RestUpstreamMapper[C, U]
	downstream RestDownstreamMapper[R]
}

// NewRestService creates a REST service with injectible mappers and repository.
// Pass WithUpstreamMapper or WithDownstreamMapper to override the defaults.
func NewRestService[C, R, U any](
	repo RestRepository,
	opts ...RestServiceOption,
) (Service[C, R, U], error) {
	if repo == nil {
		return nil, errRestServiceRepositoryNil
	}
	o := newRestServiceOptions[C, R, U](opts)
	return &restService[C, R, U]{
		repo:       repo,
		upstream:   o.upstream,
		downstream: o.downstream,
	}, nil
}

// Create implements Service[C, R, U].Create by mapping the service request through mappers.
// nolint:dupl // Similar pattern to Update but calls repo.Create
func (s *restService[C, R, U]) Create(ctx context.Context, request ServiceCreate[C], response *ServiceResponse[R]) error {
	if response == nil {
		return errRestServiceResponseNil
	}
	restReq := RestRequest{}
	if err := s.mapUpstreamContext(request.Context, &restReq.Context); err != nil {
		return errRestServiceUpstreamMappingFailed("Create/Context", err)
	}
	if err := s.upstream.ToUpstreamPost(request.Body, request.Context.Identifiers, &restReq.Body); err != nil {
		return errRestServiceUpstreamMappingFailed("Create", err)
	}
	var restResp RestResponse
	if err := s.repo.Create(ctx, restReq, &restResp); err != nil {
		return err
	}
	if err := s.mapDownstreamContext(restResp.Context, &response.Context); err != nil {
		return err
	}
	return s.mapRestResponseToServiceResponse(restResp, &response.Data, &response.Meta)
}

// List implements Service[C, R, U].List by fetching from repository and mapping response.
// nolint:dupl // Similar upstream request mapping flow to Get/Delete.
func (s *restService[C, R, U]) List(ctx context.Context, reqCtx RequestContext, response *ServiceResponses[R]) error {
	if response == nil {
		return errRestServiceResponseNil
	}
	upsCtx := RequestContext{}
	if err := s.mapUpstreamContext(reqCtx, &upsCtx); err != nil {
		return errRestServiceUpstreamMappingFailed("List/Context", err)
	}
	var restResp RestResponses
	if err := s.repo.List(ctx, upsCtx, &restResp); err != nil {
		return err
	}
	if err := s.mapDownstreamContext(restResp.Context, &response.Context); err != nil {
		return err
	}
	return s.mapRestResponsesToServiceResponses(restResp, &response.Data, &response.Meta)
}

// Get implements Service[C, R, U].Get by fetching from repository.
// nolint:dupl // Similar upstream request mapping flow to List/Delete.
func (s *restService[C, R, U]) Get(ctx context.Context, reqCtx RequestContext, response *ServiceResponse[R]) error {
	if response == nil {
		return errRestServiceResponseNil
	}
	upsCtx := RequestContext{}
	if err := s.mapUpstreamContext(reqCtx, &upsCtx); err != nil {
		return errRestServiceUpstreamMappingFailed("Get/Context", err)
	}
	var restResp RestResponse
	if err := s.repo.Get(ctx, upsCtx, &restResp); err != nil {
		return err
	}
	if err := s.mapDownstreamContext(restResp.Context, &response.Context); err != nil {
		return err
	}
	return s.mapRestResponseToServiceResponse(restResp, &response.Data, &response.Meta)
}

// Update implements Service[C, R, U].Update by mapping update request through mappers.
// nolint:dupl // Similar to Create, but calls repo.Update instead of repo.Create
func (s *restService[C, R, U]) Update(ctx context.Context, request ServiceUpdate[U], response *ServiceResponse[R]) error {
	if response == nil {
		return errRestServiceResponseNil
	}
	restReq := RestRequest{}
	if err := s.mapUpstreamContext(request.Context, &restReq.Context); err != nil {
		return errRestServiceUpstreamMappingFailed("Update/Context", err)
	}
	if err := s.upstream.ToUpstreamPatch(request.Body, request.Context.Identifiers, &restReq.Body); err != nil {
		return errRestServiceUpstreamMappingFailed("Update", err)
	}
	var restResp RestResponse
	if err := s.repo.Update(ctx, restReq, &restResp); err != nil {
		return err
	}
	if err := s.mapDownstreamContext(restResp.Context, &response.Context); err != nil {
		return err
	}
	return s.mapRestResponseToServiceResponse(restResp, &response.Data, &response.Meta)
}

// Save implements Service[C, R, U].Save by mapping the save request through mappers.
// nolint:dupl // Similar to Create flow but calls repo.Save and ToUpstreamPut
func (s *restService[C, R, U]) Save(ctx context.Context, request ServiceSave[C], response *ServiceResponse[R]) error {
	if response == nil {
		return errRestServiceResponseNil
	}
	restReq := RestRequest{}
	if err := s.mapUpstreamContext(request.Context, &restReq.Context); err != nil {
		return errRestServiceUpstreamMappingFailed("Save/Context", err)
	}
	if err := s.upstream.ToUpstreamPut(request.Body, request.Context.Identifiers, &restReq.Body); err != nil {
		return errRestServiceUpstreamMappingFailed("Save", err)
	}
	var restResp RestResponse
	if err := s.repo.Save(ctx, restReq, &restResp); err != nil {
		return err
	}
	if err := s.mapDownstreamContext(restResp.Context, &response.Context); err != nil {
		return err
	}
	return s.mapRestResponseToServiceResponse(restResp, &response.Data, &response.Meta)
}

// Delete implements Service[C, R, U].Delete by calling repository and mapping response.
// nolint:dupl // Similar upstream request mapping flow to List/Get.
func (s *restService[C, R, U]) Delete(ctx context.Context, reqCtx RequestContext, response *ServiceResponse[R]) error {
	if response == nil {
		return errRestServiceResponseNil
	}
	upsCtx := RequestContext{}
	if err := s.mapUpstreamContext(reqCtx, &upsCtx); err != nil {
		return errRestServiceUpstreamMappingFailed("Delete/Context", err)
	}
	var restResp RestResponse
	if err := s.repo.Delete(ctx, upsCtx, &restResp); err != nil {
		return err
	}
	if err := s.mapDownstreamContext(restResp.Context, &response.Context); err != nil {
		return err
	}
	return s.mapRestResponseToServiceResponse(restResp, &response.Data, &response.Meta)
}
