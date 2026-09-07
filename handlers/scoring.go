package handlers

// ScoreQuestionResponse applies FP proficiency scoring rules for a single answer.
// Yes earns the full weight and counts toward possible; No earns 0 but counts;
// NA is excluded from both achieved and possible totals.
func ScoreQuestionResponse(response string, weight int) (pointsEarned int, countsTowardPossible bool) {
	switch response {
	case "Yes":
		return weight, true
	case "No":
		return 0, true
	default:
		return 0, false
	}
}

// PercentageScore returns achieved/possible as a percentage (0 when possible is 0).
func PercentageScore(achieved, possible int) float64 {
	if possible <= 0 {
		return 0
	}
	return (float64(achieved) / float64(possible)) * 100
}

// PerformanceLevelFromPercentage maps a percentage to Proficient / Competent / Not Acceptable.
func PerformanceLevelFromPercentage(percentage float64) string {
	if percentage > 90 {
		return "Proficient"
	}
	if percentage >= 70 {
		return "Competent"
	}
	return "Not Acceptable"
}
