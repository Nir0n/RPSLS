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
	var outcomes = [3]string{"win", "lose", "tie"}
	var combinations = map[int]int{11: 2, 12: 1, 13: 0, 14: 0, 15: 1, 21: 0, 22: 2, 23: 1, 24: 1, 25: 0, 31: 1, 32: 0, 33: 2, 34: 1, 35: 1, 41: 1, 42: 0, 43: 1, 44: 2, 45: 0, 51: 0, 52: 1, 53: 0, 54: 1, 55: 2}
	res := combinations[choice1*10+choice2]
	return outcomes[res]

}
