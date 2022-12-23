package rest

type ReaderLogic interface {
	SignUp()
	SignOut()
	SignIn()
}

type BookLogic interface {
	Add()
	Get()
	Fetch()
	AddToFavorites()
	AddToWishlist()
	ExportLibrary()
}
