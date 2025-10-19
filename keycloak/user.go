package keycloak

type keycloak struct {
	host         string
	clientSecret string
}

func New(host, clientSecret string) (*keycloak, error) {
	k := &keycloak{
		host:         host,
		clientSecret: clientSecret,
	}
	return k, nil
}

func (k keycloak) CreateClubRootUser() error {
	return nil
}
