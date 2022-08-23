package ports

type HTTPServer interface {
	SetupRouter()
	Run() error
}
