package views

func StringArrToStr(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}
