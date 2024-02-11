package mashres

//func TestCreatingApplicationPackageKeySendsCeilingData(t *testing.T) {
//	h := CreateTestResource(ApplicationPackageKeyResource)
//
//	expIdent := mashschemag.ApplicationPackageKeyIdentifier{
//		PackageKeyId: "pack-key-abc",
//	}
//
//	h.givenStateFieldSetTo(t, mashschema.ApplicationPackageKeyQpsLimitCeiling, 0)
//	h.givenStateFieldSetTo(t, mashschema.ApplicationPackageKeyRateLimitCeiling, 0)
//
//	givenCreatingPackageKeySucceeds(h, expIdent)
//	h.thenExecutingCreate(t)
//	h.thenAssignedIdIs(t, func(t *testing.T, id masherytypes.PackageIdentifier) {
//		assert.True(t, reflect.DeepEqual(expIdent, id))
//	})
//}
//
//func givenCreatingPackageKeySucceeds(h *ResourceTemplateMockHelper[
//	masherytypes.ApplicationIdentifier,
//	mashschemag.ApplicationPackageKeyIdentifier, masherytypes.Package], ident masherytypes.PackageIdentifier) {
//	apiKey := "api-key"
//	apiSecret := "api-secret"
//
//	zero := int64(0)
//
//	rv := masherytypes.PackageKey{
//		AddressableV3Object: masherytypes.AddressableV3Object{
//			Id:   ident.PackageId,
//			Name: "Created Package key",
//		},
//		Apikey:           &apiKey,
//		Secret:           &apiSecret,
//		RateLimitCeiling: &zero,
//		RateLimitExempt:  false,
//		QpsLimitCeiling:  &zero,
//		QpsLimitExempt:   false,
//		Status:           "active",
//	}
//
//	h.mockClientWill().
//		On("CreatePackage", mock.Anything, mock.Anything).
//		Return(rv, nil).
//		Once()
//}
