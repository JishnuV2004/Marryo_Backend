package repositories

type Repository interface {
	Create(req interface{}) error
	FindOne(model interface{}, query string, args ...any) error
	Save(model interface{}) error
	FindByID(model interface{}, Id uint, preloads ... string) error
	// FindByIDWithPreload(model interface{}, Id uint) error
	DeleteByID(model interface{}, id uint) error
}