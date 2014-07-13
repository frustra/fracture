package network

type Server interface {
	Serve() error

	NodePort() int
}
