package finder

func NewUserFinder(impl string) UserFinder {
	switch impl {
	case "ranged":
		return NewRangedUserFinder()

	case "tested":
		return NewTestedUserFinder()

	default:
		panic("invalid user finder impl")
	}
}
