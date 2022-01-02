package domain

type Port struct {
	Name            string
	ServicePort     int32
	NetworkProtocol string
	ExposedPort     int32
}
