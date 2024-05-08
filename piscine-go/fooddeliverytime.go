package piscine

type food struct {
	preptime int
}

func FoodDeliveryTime(order string) int {
	var burger, chips, nuggets food
	burger.preptime = 15
	chips.preptime = 10
	nuggets.preptime = 12
	switch order {
	case "burger":
		return burger.preptime
	case "chips":
		return chips.preptime
	case "nuggets":
		return nuggets.preptime
	}
	return 404
}
