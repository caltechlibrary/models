package models

var (
	// Debug turns on verbose log output for model debugging
	Debug bool
)

// SetDebug sets the model's Debug flag to debug
func SetDebug(debug bool) {
	Debug = debug
}
