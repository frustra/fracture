package network

type Server interface {
	Serve()

	NodeType() string
	NodePort() int
}
