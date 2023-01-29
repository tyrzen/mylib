package usecases

type Author struct {
	repo AuthorRepository
}

func NewAuthor(repo AuthorRepository) Author {
	return Author{
		repo: repo,
	}
}
