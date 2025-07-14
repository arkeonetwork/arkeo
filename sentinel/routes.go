package sentinel

const (
	RoutesMetaData       = "/metadata.json"
	RoutesActiveContract = "/active-contract/{service}/{spender}"
	RoutesClaim          = "/claim/{id}"
	RoutesOpenClaims     = "/open-claims"
	RoutesClaims         = "/claims"
	RouteManage          = "/manage/contract/{id}"
	RouteProviderData    = "/provider/{service}"
)
