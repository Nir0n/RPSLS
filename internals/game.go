package internals

func AnnounceRandomChoice(rand int) int {
	switch {
	case rand <= 20:
		return 1
	case 20 < rand && rand <= 40:
		return 2
	case 40 < rand && rand <= 60:
		return 3
	case 60 < rand && rand <= 80:
		return 4
	default:
		return 5
	}
}

func ResultCalculator(choice1 int, choice2 int) string {
	res := Combinations[choice1*10+choice2]
	return Outcomes[res]

}
