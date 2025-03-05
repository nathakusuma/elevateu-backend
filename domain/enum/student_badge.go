package enum

type StudentBadge string

const (
	BadgeNothing StudentBadge = "nothing"
	BadgeBronze  StudentBadge = "bronze"
	BadgeSilver  StudentBadge = "silver"
	BadgeGold    StudentBadge = "gold"
)

func GetBadge(points int) StudentBadge {
	switch {
	case points >= 2000:
		return BadgeGold
	case points >= 1000:
		return BadgeSilver
	case points >= 500:
		return BadgeBronze
	default:
		return BadgeNothing
	}
}
