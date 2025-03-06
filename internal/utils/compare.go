package utils

import (
    "strconv"
    "strings"
)

func CompareVersions(v1, v2 string) int {
    parts1 := strings.Split(v1, ".")
    parts2 := strings.Split(v2, ".")

    for i := 0; i < len(parts1) && i < len(parts2); i++ {
        num1, _ := strconv.Atoi(parts1[i])
        num2, _ := strconv.Atoi(parts2[i])

        if num1 < num2 {
            return -1
        }
        if num1 > num2 {
            return 1
        }
    }

    if len(parts1) < len(parts2) {
        return -1
    }
    if len(parts1) > len(parts2) {
        return 1
    }

    return 0
}
