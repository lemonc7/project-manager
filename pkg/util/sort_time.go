package util

import "github.com/lemonc7/project-manager/internal/model"

func SortByRemaining(a, b model.HasRemainingDays) int {
	switch {
	case a.GetRemainingDays() != nil && b.GetRemainingDays() == nil:
		return -1
	case a.GetRemainingDays() == nil && b.GetRemainingDays() != nil:
		return 1
	case a.GetRemainingDays() != nil && b.GetRemainingDays() != nil:
		if *a.GetRemainingDays() < *b.GetRemainingDays() {
			return -1
		} else if *a.GetRemainingDays() > *b.GetRemainingDays() {
			return 1
		}
	}
	return 0
}
