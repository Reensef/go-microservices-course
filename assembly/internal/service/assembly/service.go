package assembly

type service struct{}

func NewService() *service {
	return &service{}
}
