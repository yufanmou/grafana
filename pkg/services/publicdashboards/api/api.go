package api

import (
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/publicdashboards"
	"github.com/grafana/grafana/pkg/services/publicdashboards/internal/tokens"
	. "github.com/grafana/grafana/pkg/services/publicdashboards/models"
	"github.com/grafana/grafana/pkg/web"
)

type Api struct {
	PublicDashboardService publicdashboards.Service
	RouteRegister          routing.RouteRegister
	AccessControl          accesscontrol.AccessControl
	Features               *featuremgmt.FeatureManager
	Log                    log.Logger
}

func ProvideApi(
	pd publicdashboards.Service,
	rr routing.RouteRegister,
	ac accesscontrol.AccessControl,
	features *featuremgmt.FeatureManager,
) *Api {
	api := &Api{
		PublicDashboardService: pd,
		RouteRegister:          rr,
		AccessControl:          ac,
		Features:               features,
		Log:                    log.New("publicdashboards.api"),
	}

	// attach api if PublicDashboards feature flag is enabled
	if features.IsEnabled(featuremgmt.FlagPublicDashboards) {
		api.RegisterAPIEndpoints()
	}

	return api
}

// RegisterAPIEndpoints Registers Endpoints on Grafana Router
func (api *Api) RegisterAPIEndpoints() {
	// Public endpoints
	// Anonymous access to public dashboard route is configured in pkg/api/api.go
	// because it is deeply dependent on the HTTPServer.Index() method and would result in a
	// circular dependency

	api.RouteRegister.Get("/api/public/dashboards/:accessToken", routing.Wrap(api.ViewPublicDashboard))
	api.RouteRegister.Post("/api/public/dashboards/:accessToken/panels/:panelId/query", routing.Wrap(api.QueryPublicDashboard))
	api.RouteRegister.Get("/api/public/dashboards/:accessToken/annotations", routing.Wrap(api.GetAnnotations))

	// Auth endpoints
	auth := accesscontrol.Middleware(api.AccessControl)
	uidScope := dashboards.ScopeDashboardsProvider.GetResourceScopeUID(accesscontrol.Parameter(":dashboardUid"))

	// List public dashboards for org
	api.RouteRegister.Get("/api/dashboards/public-dashboards", middleware.ReqSignedIn, routing.Wrap(api.ListPublicDashboards))

	// Get public dashboard
	api.RouteRegister.Get("/api/dashboards/uid/:dashboardUid/public-dashboards",
		auth(middleware.ReqSignedIn, accesscontrol.EvalPermission(dashboards.ActionDashboardsRead, uidScope)),
		routing.Wrap(api.GetPublicDashboard))

	// Create Public Dashboard
	api.RouteRegister.Post("/api/dashboards/uid/:dashboardUid/public-dashboards",
		auth(middleware.ReqOrgAdmin, accesscontrol.EvalPermission(dashboards.ActionDashboardsPublicWrite, uidScope)),
		routing.Wrap(api.CreatePublicDashboard))

	// Update Public Dashboard
	api.RouteRegister.Put("/api/dashboards/uid/:dashboardUid/public-dashboards/:uid",
		auth(middleware.ReqOrgAdmin, accesscontrol.EvalPermission(dashboards.ActionDashboardsPublicWrite, uidScope)),
		routing.Wrap(api.UpdatePublicDashboard))

	// Delete Public dashboard
	api.RouteRegister.Delete("/api/dashboards/uid/:dashboardUid/public-dashboards/:uid",
		auth(middleware.ReqOrgAdmin, accesscontrol.EvalPermission(dashboards.ActionDashboardsPublicWrite, uidScope)),
		routing.Wrap(api.DeletePublicDashboard))
}

// ListPublicDashboards Gets list of public dashboards by orgId
// GET /api/dashboards/public-dashboards
func (api *Api) ListPublicDashboards(c *models.ReqContext) response.Response {
	resp, err := api.PublicDashboardService.FindAll(c.Req.Context(), c.SignedInUser, c.OrgID)
	if err != nil {
		return response.Err(err)
	}
	return response.JSON(http.StatusOK, resp)
}

// GetPublicDashboard Gets public dashboard for dashboard
// GET /api/dashboards/uid/:dashboardUid/public-dashboards
func (api *Api) GetPublicDashboard(c *models.ReqContext) response.Response {
	// exit if we don't have a valid dashboardUid
	dashboardUid := web.Params(c.Req)[":dashboardUid"]
	if !tokens.IsValidShortUID(dashboardUid) {
		return response.Err(ErrPublicDashboardIdentifierNotSet.Errorf("GetPublicDashboard: no dashboard Uid for public dashboard specified"))
	}

	pd, err := api.PublicDashboardService.FindByDashboardUid(c.Req.Context(), c.OrgID, dashboardUid)
	if err != nil {
		return response.Err(err)
	}

	if pd == nil {
		response.Err(ErrPublicDashboardNotFound.Errorf("GetPublicDashboard: public dashboard not found"))
	}

	return response.JSON(http.StatusOK, pd)
}

// CreatePublicDashboard Sets public dashboard for dashboard
// POST /api/dashboards/uid/:dashboardUid/public-dashboards
func (api *Api) CreatePublicDashboard(c *models.ReqContext) response.Response {
	// exit if we don't have a valid dashboardUid
	dashboardUid := web.Params(c.Req)[":dashboardUid"]
	if !tokens.IsValidShortUID(dashboardUid) {
		return response.Err(ErrInvalidUid.Errorf("CreatePublicDashboard: invalid Uid %s", dashboardUid))
	}

	pd := &PublicDashboard{}
	if err := web.Bind(c.Req, pd); err != nil {
		return response.Err(ErrBadRequest.Errorf("CreatePublicDashboard: bad request data %v", err))
	}

	// Always set the orgID and userID from the session
	pd.OrgId = c.OrgID
	dto := SavePublicDashboardDTO{
		UserId:          c.UserID,
		OrgId:           c.OrgID,
		DashboardUid:    dashboardUid,
		PublicDashboard: pd,
	}

	//Create the public dashboard
	pd, err := api.PublicDashboardService.Create(c.Req.Context(), c.SignedInUser, &dto)
	if err != nil {
		return response.Err(err)
	}

	return response.JSON(http.StatusOK, pd)
}

// UpdatePublicDashboard Sets public dashboard for dashboard
// PUT /api/dashboards/uid/:dashboardUid/public-dashboards/:uid
func (api *Api) UpdatePublicDashboard(c *models.ReqContext) response.Response {
	// exit if we don't have a valid dashboardUid
	dashboardUid := web.Params(c.Req)[":dashboardUid"]
	if !tokens.IsValidShortUID(dashboardUid) {
		return response.Err(ErrInvalidUid.Errorf("UpdatePublicDashboard: invalid dashboard Uid %s", dashboardUid))
	}

	uid := web.Params(c.Req)[":uid"]
	if !tokens.IsValidShortUID(uid) {
		return response.Err(ErrInvalidUid.Errorf("UpdatePublicDashboard: invalid Uid %s", uid))
	}

	pd := &PublicDashboard{}
	if err := web.Bind(c.Req, pd); err != nil {
		return response.Err(ErrBadRequest.Errorf("UpdatePublicDashboard: bad request data %v", err))
	}

	// Always set the orgID and userID from the session
	pd.OrgId = c.OrgID
	pd.Uid = uid
	dto := SavePublicDashboardDTO{
		UserId:          c.UserID,
		OrgId:           c.OrgID,
		DashboardUid:    dashboardUid,
		PublicDashboard: pd,
	}

	// Update the public dashboard
	pd, err := api.PublicDashboardService.Update(c.Req.Context(), c.SignedInUser, &dto)
	if err != nil {
		return response.Err(err)
	}

	return response.JSON(http.StatusOK, pd)
}

// Delete a public dashboard
// DELETE /api/dashboards/uid/:dashboardUid/public-dashboards/:uid
func (api *Api) DeletePublicDashboard(c *models.ReqContext) response.Response {
	uid := web.Params(c.Req)[":uid"]
	if !tokens.IsValidShortUID(uid) {
		return response.Err(ErrInvalidUid.Errorf("UpdatePublicDashboard: invalid Uid %s", uid))
	}

	err := api.PublicDashboardService.Delete(c.Req.Context(), c.OrgID, uid)
	if err != nil {
		return response.Err(err)
	}

	return response.JSON(http.StatusOK, nil)
}

// Copied from pkg/api/metrics.go
func toJsonStreamingResponse(features *featuremgmt.FeatureManager, qdr *backend.QueryDataResponse) response.Response {
	statusWhenError := http.StatusBadRequest
	if features.IsEnabled(featuremgmt.FlagDatasourceQueryMultiStatus) {
		statusWhenError = http.StatusMultiStatus
	}

	statusCode := http.StatusOK
	for _, res := range qdr.Responses {
		if res.Error != nil {
			statusCode = statusWhenError
		}
	}

	return response.JSONStreaming(statusCode, qdr)
}
