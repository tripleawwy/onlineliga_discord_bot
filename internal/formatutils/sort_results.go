package formatutils

import (
	"github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
)

// SortResults sorts the results by league and position
func SortResults(results [][]string, logger *logrus.Logger) [][]string {
	// Sort the results by league and position
	sort.Slice(results, func(i, j int) bool {
		// Compare the league
		if results[i][0] != results[j][0] {
			return results[i][0] < results[j][0]
		}
		// Compare the position
		// Consider P1 < P2 < P10
		if results[i][1] != results[j][1] {
			// Convert the position to an int
			iPosition, err := convertPositionToInt(results[i][1])
			if err != nil {
				logger.WithError(err).Error("Failed to convert position to int")
				return false
			}
			jPosition, err := convertPositionToInt(results[j][1])
			if err != nil {
				logger.WithError(err).Error("Failed to convert position to int")
				return false
			}
			return iPosition < jPosition
		}

		// Compare the home team
		if results[i][2] != results[j][2] {
			return results[i][2] < results[j][2]
		}
		// Compare the result
		if results[i][3] != results[j][3] {
			return results[i][3] < results[j][3]
		}
		// Compare the away team
		return results[i][4] < results[j][4]
	})

	return results
}

func convertPositionToInt(s string) (int, error) {
	// Strip the "P" from the position
	s = strings.Replace(s, "P", "", 1)
	// Convert the numerical string to an int
	position, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return position, nil
}
