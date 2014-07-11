package network

type Server interface {
	Serve() error

	NodeType() string
	NodePort() int
}
