package utils

import "errors"

func ForwardnessToString(value uint8) (string, error) {
	switch value {
	case 0:
		return "GK", nil
	case 1:
		return "D", nil
	case 2:
		return "M", nil
	case 3:
		return "F", nil
	case 4:
		return "MF", nil
	case 5:
		return "MD", nil
	default:
		return "", errors.New("unexistent forwardness")
	}
}
