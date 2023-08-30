package testsuite

func SetPanicFunc(f func(message string)) {
	panicFunc = f
}

func PanicFunc() func(message string) {
	return panicFunc
}
