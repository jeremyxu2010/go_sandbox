package impl

import (
	profilesvcapi "personal/jeremyxu/sandbox/test/profilesvc/api"
)

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a profilesvc
// server.
func MakeServerEndpoints(s profilesvcapi.Service) profilesvcapi.Endpoints {
	return profilesvcapi.Endpoints{
		PostProfileEndpoint:   profilesvcapi.MakePostProfileEndpoint(s),
		GetProfileEndpoint:    profilesvcapi.MakeGetProfileEndpoint(s),
		PutProfileEndpoint:    profilesvcapi.MakePutProfileEndpoint(s),
		PatchProfileEndpoint:  profilesvcapi.MakePatchProfileEndpoint(s),
		DeleteProfileEndpoint: profilesvcapi.MakeDeleteProfileEndpoint(s),
		GetAddressesEndpoint:  profilesvcapi.MakeGetAddressesEndpoint(s),
		GetAddressEndpoint:    profilesvcapi.MakeGetAddressEndpoint(s),
		PostAddressEndpoint:   profilesvcapi.MakePostAddressEndpoint(s),
		DeleteAddressEndpoint: profilesvcapi.MakeDeleteAddressEndpoint(s),
	}
}
