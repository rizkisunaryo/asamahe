package util

// import "fmt"
// import "strconv"
import (
    "strconv"
)


func FloatToString(input_num float64) string {
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}