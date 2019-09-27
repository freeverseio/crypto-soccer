package utils

import "errors"

func PreferredPosition(forwardness uint8, leftishness uint8) (string, error) {
	f, err := ForwardnessToString(forwardness)
	if err != nil {
		return "", err
	}
	l, err := LeftishnessToString(leftishness)
	if err != nil {
		return "", err
	}
	if len(l) == 0 {
		return f, nil
	}
	return f + " " + l, nil
}

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

func LeftishnessToString(value uint8) (string, error) {
	if value >= 8 {
		return "", errors.New("unexistent leftishness")
	}
	var result string
	if (value & (0x1 << 2)) != 0 {
		result += "L"
	}
	if (value & (0x1 << 1)) != 0 {
		result += "C"
	}
	if (value & 0x1) != 0 {
		result += "R"
	}
	return result, nil
}
