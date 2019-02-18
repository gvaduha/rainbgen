package sha2svc

import "os"

// GetEnvOrDefault return named variable environment value or default
func GetEnvOrDefault(name string, defvalue string) (ret string) {
	if ret = os.Getenv(name); ret == "" {
		ret = defvalue
	}

	return
}
