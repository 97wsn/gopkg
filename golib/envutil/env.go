package envutil

import "os"

const (
	Prod = "prod"
	Test = "test"
	Dev  = "dev"
)

func IsProd() bool {
	return os.Getenv("APP_ENV") == Prod
}

func IsTest() bool {
	return os.Getenv("APP_ENV") == Test
}

func IsDev() bool {
	return os.Getenv("APP_ENV") == Dev
}

func Env() string {
	return os.Getenv("APP_ENV")
}
