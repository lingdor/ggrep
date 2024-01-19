package util

func CheckPanic(err error) {
	if err != nil {
		panic(err)
	}
}
