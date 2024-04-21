package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	v3client.Client
	mock.Mock
}

func (mc *MockClient) GetEndpoint(ctx context.Context, ident masherytypes.ServiceEndpointIdentifier) (masherytypes.Endpoint, bool, error) {
	args := mc.Called(ctx, ident)
	return args.Get(0).(masherytypes.Endpoint), args.Bool(1), args.Error(2)
}

func (mc *MockClient) ListMembersFiltered(ctx context.Context, params map[string]string) ([]masherytypes.Member, error) {
	args := mc.Called(ctx, params)
	return args.Get(0).([]masherytypes.Member), args.Error(1)
}

func (mc *MockClient) ListApplicationsFiltered(ctx context.Context, params map[string]string) ([]masherytypes.Application, error) {
	args := mc.Called(ctx, params)
	return args.Get(0).([]masherytypes.Application), args.Error(1)
}

func (mc *MockClient) ListRolesFiltered(ctx context.Context, params map[string]string) ([]masherytypes.Role, error) {
	args := mc.Called(ctx, params)
	return args.Get(0).([]masherytypes.Role), args.Error(1)
}

func (mc *MockClient) ListEmailTemplateSetsFiltered(ctx context.Context, params map[string]string) ([]masherytypes.EmailTemplateSet, error) {
	args := mc.Called(ctx, params)
	return args.Get(0).([]masherytypes.EmailTemplateSet), args.Error(1)
}

func (mc *MockClient) CreateService(ctx context.Context, service masherytypes.Service) (masherytypes.Service, error) {
	return servicePointerWithError(mc.Called(ctx, service))
}

func servicePointerWithError(args mock.Arguments) (masherytypes.Service, error) {
	var rv = args.Get(0).(masherytypes.Service)

	return rv, args.Error(1)
}

func (mc *MockClient) SetServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier, roles []masherytypes.RolePermission) error {
	args := mc.Called(ctx, id, roles)
	return args.Error(0)
}

func (mc *MockClient) DeleteService(ctx context.Context, serviceId masherytypes.ServiceIdentifier) error {
	args := mc.Called(ctx, serviceId)
	return args.Error(0)
}

func (mc *MockClient) CountEndpointsOf(ctx context.Context, serviceId masherytypes.ServiceIdentifier) (int64, error) {
	args := mc.Called(ctx, serviceId)
	return args.Get(0).(int64), args.Error(1)
}

func (mc *MockClient) GetService(ctx context.Context, id masherytypes.ServiceIdentifier) (masherytypes.Service, bool, error) {
	args := mc.Called(ctx, id)
	var rv = args.Get(0).(masherytypes.Service)

	return rv, args.Bool(1), args.Error(2)
}

func (mc *MockClient) GetServiceRoles(ctx context.Context, serviceId masherytypes.ServiceIdentifier) ([]masherytypes.RolePermission, bool, error) {
	args := mc.Called(ctx, serviceId)
	return args.Get(0).([]masherytypes.RolePermission), args.Bool(1), args.Error(2)
}

func (mc *MockClient) DeleteServiceRoles(ctx context.Context, id masherytypes.ServiceIdentifier) error {
	args := mc.Called(ctx, id)
	return args.Error(0)
}

func (mc *MockClient) UpdateService(ctx context.Context, service masherytypes.Service) (masherytypes.Service, error) {
	return servicePointerWithError(mc.Called(ctx, service))
}

func (mc *MockClient) CreateServiceOAuthSecurityProfile(ctx context.Context, id masherytypes.ServiceIdentifier, service masherytypes.MasheryOAuth) (masherytypes.MasheryOAuth, error) {
	args := mc.Called(ctx, id, service)
	var rv = args.Get(0).(masherytypes.MasheryOAuth)

	return rv, args.Error(1)
}

func (mc *MockClient) GetServiceCache(ctx context.Context, id masherytypes.ServiceIdentifier) (masherytypes.ServiceCache, bool, error) {
	args := mc.Called(ctx, id)
	var rv = args.Get(0).(masherytypes.ServiceCache)

	return rv, args.Bool(1), args.Error(2)
}

func serviceCacheAndErrorFrom(args mock.Arguments) (masherytypes.ServiceCache, error) {
	var rv = args.Get(0).(masherytypes.ServiceCache)

	return rv, args.Error(1)
}

func (mc *MockClient) CreateServiceCache(ctx context.Context, id masherytypes.ServiceIdentifier, service masherytypes.ServiceCache) (masherytypes.ServiceCache, error) {
	args := mc.Called(ctx, id, service)
	return serviceCacheAndErrorFrom(args)
}

func (mc *MockClient) CreateEndpoint(ctx context.Context, serviceId masherytypes.ServiceIdentifier, endp masherytypes.Endpoint) (masherytypes.Endpoint, error) {
	args := mc.Called(ctx, serviceId, endp)
	var rv = args.Get(0).(masherytypes.Endpoint)

	return rv, args.Error(1)
}

func (mc *MockClient) CreatePackage(ctx context.Context, pack masherytypes.Package) (masherytypes.Package, error) {
	args := mc.Called(ctx, pack)
	var rv = args.Get(0).(masherytypes.Package)

	return rv, args.Error(1)
}
func (mc *MockClient) CreatePlan(ctx context.Context, packageId masherytypes.PackageIdentifier, plan masherytypes.Plan) (masherytypes.Plan, error) {
	args := mc.Called(ctx, packageId, plan)
	var rv masherytypes.Plan = args.Get(0).(masherytypes.Plan)

	return rv, args.Error(1)
}

func (mc *MockClient) CreatePlanService(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier) (masherytypes.AddressableV3Object, error) {
	args := mc.Called(ctx, planService)
	var rv = args.Get(0).(masherytypes.AddressableV3Object)

	return rv, args.Error(1)
}

func (mc *MockClient) CheckPlanServiceExists(ctx context.Context, planService masherytypes.PackagePlanServiceIdentifier) (bool, error) {
	args := mc.Called(ctx, planService)
	return args.Bool(0), args.Error(1)
}
func (mc *MockClient) CreatePlanEndpoint(ctx context.Context, planEndp masherytypes.PackagePlanServiceEndpointIdentifier) (masherytypes.AddressableV3Object, error) {
	args := mc.Called(ctx, planEndp)
	var rv = args.Get(0).(masherytypes.AddressableV3Object)

	return rv, args.Error(1)
}

func (mc *MockClient) CheckPlanEndpointExists(ctx context.Context, planEndp masherytypes.PackagePlanServiceEndpointIdentifier) (bool, error) {
	args := mc.Called(ctx, planEndp)
	return args.Bool(0), args.Error(1)
}

func (mc *MockClient) CreateEndpointMethod(ctx context.Context, ident masherytypes.ServiceEndpointIdentifier, methodUpsert masherytypes.ServiceEndpointMethod) (masherytypes.ServiceEndpointMethod, error) {
	args := mc.Called(ctx, ident, methodUpsert)
	var rv = args.Get(0).(masherytypes.ServiceEndpointMethod)

	return rv, args.Error(1)
}

func (mc *MockClient) CreateEndpointMethodFilter(ctx context.Context, ident masherytypes.ServiceEndpointMethodIdentifier, filterUpsert masherytypes.ServiceEndpointMethodFilter) (masherytypes.ServiceEndpointMethodFilter, error) {
	args := mc.Called(ctx, ident, filterUpsert)
	var rv = args.Get(0).(masherytypes.ServiceEndpointMethodFilter)

	return rv, args.Error(1)
}

func (mc *MockClient) CreatePackagePlanMethod(ctx context.Context, id masherytypes.PackagePlanServiceEndpointMethodIdentifier) (masherytypes.PackagePlanServiceEndpointMethod, error) {
	args := mc.Called(ctx, id)
	var rv = args.Get(0).(masherytypes.PackagePlanServiceEndpointMethod)

	return rv, args.Error(1)
}

func (mc *MockClient) CreatePackagePlanMethodFilter(ctx context.Context, id masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier) (masherytypes.PackagePlanServiceEndpointMethodFilter, error) {
	args := mc.Called(ctx, id)
	var rv = args.Get(0).(masherytypes.PackagePlanServiceEndpointMethodFilter)

	return rv, args.Error(1)
}

func (mc *MockClient) ListOrganizationsFiltered(ctx context.Context, qs map[string]string) ([]masherytypes.Organization, error) {
	args := mc.Called(ctx, qs)
	return args.Get(0).([]masherytypes.Organization), args.Error(1)
}

func (mc *MockClient) CreateErrorSet(ctx context.Context, serviceId masherytypes.ServiceIdentifier, set masherytypes.ErrorSet) (masherytypes.ErrorSet, error) {
	args := mc.Called(ctx, serviceId, set)
	var rv = args.Get(0).(masherytypes.ErrorSet)

	return rv, args.Error(1)
}

func (mc *MockClient) GetErrorSet(ctx context.Context, ident masherytypes.ErrorSetIdentifier) (masherytypes.ErrorSet, bool, error) {
	args := mc.Called(ctx, ident)

	var rv = args.Get(0).(masherytypes.ErrorSet)

	return rv, args.Bool(1), args.Error(2)
}
