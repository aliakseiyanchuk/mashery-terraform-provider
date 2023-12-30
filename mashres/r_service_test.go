package mashres

import (
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
	"testing"
	"time"
)

func TestCreatingServiceWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	h.givenStateFieldSetTo(t, mashschema.MashObjNamePrefix, "name-prefix")
	givenCreateServiceSucceeds(h)
	givenReadingServiceRolesOfACreatedServiceYieldsNothing(h)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceIdentifier) {
		assert.Equal(t, id.ServiceId, "serviceId")
	})
	h.willHaveFieldSetTo(t, mashschema.MashObjName, "service-name")
}

func TestCreatingServiceWIthRolesWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	h.givenStateFieldSetTo(t, mashschema.MashObjNamePrefix, "name-prefix")
	h.givenStateFieldSetTo(t, mashschema.MashSvcInteractiveDocsRoles, []string{"role-id"})
	givenCreateServiceSucceeds(h)
	givenReadingServiceRolesOfACreatedServiceYieldsData(h)

	h.thenExecutingCreate(t)
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceIdentifier) {
		assert.Equal(t, id.ServiceId, "serviceId")
	})
	h.willHaveFieldSetTo(t, mashschema.MashObjName, "service-name")
	h.willHaveStringArrayFieldSetTo(t, mashschema.MashSvcInteractiveDocsRoles, []string{"role-read-back-id"})
}

func TestCreatingServiceAndFailingRolesWillTaintResource(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	h.givenStateFieldSetTo(t, mashschema.MashObjNamePrefix, "name-prefix")
	h.givenStateFieldSetTo(t, mashschema.MashSvcInteractiveDocsRoles, []string{"role-id"})
	givenCreateServiceSucceeds(h)
	givenReadingServiceRolesOfACreatedServiceFails(h)

	h.thenExecutingCreateWillYieldDiagnostic(t, "unexpected error returned from Mashery V3 api during creating object")
	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.ServiceIdentifier) {
		assert.Equal(t, id.ServiceId, "serviceId")
	})
	h.willHaveFieldSetTo(t, mashschema.MashObjName, "service-name")
}

func TestCreatingServiceWillReturnDiagnosticIfNotSuccessful(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	h.givenStateFieldSetTo(t, mashschema.MashObjNamePrefix, "name-prefix")
	givenCreateServiceFails(h)

	h.thenExecutingCreateWillYieldDiagnostic(t, "unexpected error returned from Mashery V3 api during creating object")
	h.willHaveStateId(t, "")
}

func TestReadServiceWithNoRolesWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceResource)
	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, serviceId)
	givenReadServiceSucceeds(h, serviceId)
	givenReadingServiceRolesYieldsNothing(h, serviceId)

	h.thenExecutingRead(t)
	// It just works.
}

func TestReadServiceWithRolesWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceResource)
	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, serviceId)
	givenReadServiceSucceeds(h, serviceId)
	givenReadingServiceRolesYieldsData(h, serviceId)

	h.thenExecutingRead(t)
	h.willHaveStringArrayFieldSetTo(t, mashschema.MashSvcInteractiveDocsRoles, []string{"role-Id"})
}

func TestReadServiceWithReturnDiagOnAPIFailure(t *testing.T) {
	h := CreateTestResource(ServiceResource)
	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, serviceId)
	givenReadServiceFails(h, serviceId)

	h.thenExecutingReadWillYieldDiagnostic(t, "unexpected error returned from Mashery V3 api")
}

func TestReadServiceWithReturnDiagOnRolesReadAPIFailure(t *testing.T) {
	h := CreateTestResource(ServiceResource)
	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, serviceId)
	givenReadServiceSucceeds(h, serviceId)
	givenReadingServiceRolesFails(h, serviceId)

	h.thenExecutingReadWillYieldDiagnostic(t, "unexpected error returned from Mashery V3 api")
}

// ----------------------------------------------------------------------------------------
// Update operations

func TestUpdateServiceWithNoRolesWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, serviceId)
	givenUpdatingServiceSucceeds(h, serviceId)
	givenReadingServiceRolesYieldsNothing(h, serviceId)

	h.thenExecutingUpdate(t)
	// just works
}

func TestUpdateServiceWithNoRolesWillReturnDiagnosticIfServiceUpdateErrs(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, serviceId)
	givenUpdatingServiceFails(h, serviceId)

	h.thenExecutingUpdateWillYieldDiagnostic(t, "unexpected error returned from Mashery V3 api during update object")
}

func TestUpdateServiceWithRolesWillSucceed(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, serviceId)
	h.givenStateFieldSetTo(t, mashschema.MashSvcInteractiveDocsRoles, []string{"role-id"})
	givenUpdatingServiceSucceeds(h, serviceId)
	givenReadingServiceRolesYieldsData(h, serviceId)

	h.thenExecutingUpdate(t)
	// just works
}

func TestUpdateServiceWithRolesWillReturnDiagnosticIfRoleReadbackErrs(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	serviceId := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, serviceId)
	h.givenStateFieldSetTo(t, mashschema.MashSvcInteractiveDocsRoles, []string{"role-id"})
	givenUpdatingServiceSucceeds(h, serviceId)
	givenReadingServiceRolesFails(h, serviceId)

	h.thenExecutingUpdateWillYieldDiagnostic(t, "unexpected error returned from Mashery V3 api during update object")
}

// ----------------------------------------------------------------------------------------
// Delete operations

func TestDeleteServiceWillComplete(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	testIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, testIdent)
	givenThereAreEndpointOfService(h, 0)
	givenServiceDeleteCompletes(h)

	h.thenExecutingDelete(t)
	h.willHaveStateId(t, "")
}

func TestDeleteServiceWillBlockOnOffendingEndpoints(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	testIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, testIdent)
	givenThereAreEndpointOfService(h, 1)

	h.thenExecutingDeleteWillYieldDiagnostic(t, "offending objets prevent deletion")
	h.willHaveStateIdDefined(t)
}

func TestDeleteServiceWillBlockOnFailureToCountOffendingEndpoints(t *testing.T) {
	h := CreateTestResource(ServiceResource)

	testIdent := masherytypes.ServiceIdentifier{
		ServiceId: "abcdefg",
	}

	h.givenIdentity(t, testIdent)
	givenCountingServiceEndpointFails(h)

	h.thenExecutingDeleteWillYieldDiagnostic(t, "query for offending objects of has returned an error")
	h.willHaveStateIdDefined(t)
}

func givenDeletingServiceRolesSucceeds(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	h.mockClientWill().
		On("DeleteServiceRoles", mock.Anything, ident).
		Return(nil).
		Once()
}

func givenDeletingServiceRolesFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	h.mockClientWill().
		On("DeleteServiceRoles", mock.Anything, ident).
		Return(errors.New("unit text error on deleting the service roles")).
		Once()
}

func givenCreateServiceSucceeds(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	mashTime := masherytypes.MasheryJSONTime(time.Now())

	returnedService := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:        "serviceId",
			Name:      "service-name",
			Created:   &mashTime,
			Updated:   &mashTime,
			Retrieved: time.Time{},
		},
		Cache:             nil,
		Endpoints:         nil,
		EditorHandle:      "author",
		RevisionNumber:    0,
		RobotsPolicy:      "",
		CrossdomainPolicy: "",
		Description:       "",
		ErrorSets:         nil,
		QpsLimitOverall:   nil,
		RFC3986Encode:     true,
		SecurityProfile:   nil,
		Version:           "1.0",
		Roles:             nil,
	}
	h.mockClientWill().
		On("CreateService", mock.Anything, mock.Anything).
		Return(returnedService, nil).
		Once()
}

func givenReadServiceSucceeds(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	mashTime := masherytypes.MasheryJSONTime(time.Now())

	returnedService := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:        ident.ServiceId,
			Name:      "service-name",
			Created:   &mashTime,
			Updated:   &mashTime,
			Retrieved: time.Time{},
		},
		Cache:             nil,
		Endpoints:         nil,
		EditorHandle:      "author",
		RevisionNumber:    0,
		RobotsPolicy:      "",
		CrossdomainPolicy: "",
		Description:       "",
		ErrorSets:         nil,
		QpsLimitOverall:   nil,
		RFC3986Encode:     true,
		SecurityProfile:   nil,
		Version:           "1.0",
		Roles:             nil,
	}
	h.mockClientWill().
		On("GetService", mock.Anything, ident).
		Return(returnedService, true, nil).
		Once()
}

func givenUpdatingServiceSucceeds(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	mashTime := masherytypes.MasheryJSONTime(time.Now())

	returnedService := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:        ident.ServiceId,
			Name:      "service-name",
			Created:   &mashTime,
			Updated:   &mashTime,
			Retrieved: time.Time{},
		},
		Cache:             nil,
		Endpoints:         nil,
		EditorHandle:      "author",
		RevisionNumber:    0,
		RobotsPolicy:      "",
		CrossdomainPolicy: "",
		Description:       "",
		ErrorSets:         nil,
		QpsLimitOverall:   nil,
		RFC3986Encode:     true,
		SecurityProfile:   nil,
		Version:           "1.0",
		Roles:             nil,
	}
	h.mockClientWill().
		On("UpdateService", mock.Anything, mock.Anything).
		Return(returnedService, nil).
		Once()
}

func givenUpdatingServiceFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	h.mockClientWill().
		On("UpdateService", mock.Anything, mock.Anything).
		Return(masherytypes.Service{}, errors.New("unit test rejection of UpdateService method")).
		Once()
}

func givenReadServiceFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	h.mockClientWill().
		On("GetService", mock.Anything, ident).
		Return(masherytypes.Service{}, false, errors.New("unit test error on reading the service")).
		Once()
}

func givenReadingServiceRolesYieldsData(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	returnedPermissions := []masherytypes.RolePermission{
		{
			Role: masherytypes.Role{
				AddressableV3Object: masherytypes.AddressableV3Object{
					Id:   "role-Id",
					Name: "role-name",
				},
			},
			Action: "read",
		},
	}
	h.mockClientWill().
		On("GetServiceRoles", mock.Anything, ident).
		Return(returnedPermissions, true, nil).
		Once()
}

func givenReadingServiceRolesOfACreatedServiceYieldsData(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	returnedPermissions := []masherytypes.RolePermission{
		{
			Role: masherytypes.Role{
				AddressableV3Object: masherytypes.AddressableV3Object{
					Id:   "role-read-back-id",
					Name: "role-name",
				},
			},
			Action: "read",
		},
	}
	h.mockClientWill().
		On("GetServiceRoles", mock.Anything, mock.Anything).
		Return(returnedPermissions, true, nil).
		Once()
}

func givenReadingServiceRolesYieldsNothing(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	h.mockClientWill().
		On("GetServiceRoles", mock.Anything, ident).
		Return([]masherytypes.RolePermission{}, false, nil).
		Once()
}

func givenReadingServiceRolesOfACreatedServiceYieldsNothing(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	h.mockClientWill().
		On("GetServiceRoles", mock.Anything, mock.Anything).
		Return([]masherytypes.RolePermission{}, true, nil).
		Once()
}

func givenReadingServiceRolesFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	h.mockClientWill().
		On("GetServiceRoles", mock.Anything, ident).
		Return([]masherytypes.RolePermission{}, false, errors.New("unit-test rejection while reading roles")).
		Once()
}

func givenReadingServiceRolesOfACreatedServiceFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	h.mockClientWill().
		On("GetServiceRoles", mock.Anything, mock.Anything).
		Return([]masherytypes.RolePermission{}, false, errors.New("unit-test rejection while reading roles")).
		Once()
}

func givenCreateServiceRoleSucceeds(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	h.mockClientWill().
		On("SetServiceRoles", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Once()
}

func givenUpdateServiceRoleSucceeds(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	h.mockClientWill().
		On("SetServiceRoles", mock.Anything, ident, mock.Anything).
		Return(nil).
		Once()
}

func givenUpdatingServiceRoleFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], ident masherytypes.ServiceIdentifier) {
	h.mockClientWill().
		On("SetServiceRoles", mock.Anything, ident, mock.Anything).
		Return(errors.New("unit test error for SetServiceRoles method")).
		Once()
}

func givenCreateServiceRolesFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	h.mockClientWill().
		On("SetServiceRoles", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("unit test failure on creating service roles")).
		Once()
}

func givenCreateServiceFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	h.mockClientWill().
		On("CreateService", mock.Anything, mock.Anything).
		Return(masherytypes.Service{}, errors.New("service creation failed for unit test reason")).
		Once()

}

func givenThereAreEndpointOfService(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service], n int64) {
	h.mockClientWill().
		On("CountEndpointsOf", mock.Anything, mock.Anything).
		Return(n, nil).
		Once()
}

func givenCountingServiceEndpointFails(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	h.mockClientWill().
		On("CountEndpointsOf", mock.Anything, mock.Anything).
		Return(int64(0), errors.New("unit test error on counting service endpoints")).
		Once()
}

func givenServiceDeleteCompletes(h *ResourceTemplateMockHelper[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]) {
	h.mockClientWill().
		On("DeleteService", mock.Anything, mock.Anything).
		Return(nil).
		Once()
}
