package matchingvalidator

type MatchingRepo interface {
	// needed funcitons
}

type MatchingValidator struct {
	repo MatchingRepo
}

func New(repo MatchingRepo) MatchingValidator {
	return MatchingValidator{
		repo: repo,
	}
}
