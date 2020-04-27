package repository

type Repository interface {
	Health() error
	ProjectRepository
	UserRepository
}
