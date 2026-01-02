package resource

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

type resource struct {
	Api struct {
		Url string `json:"url"`
	}
	Auth struct {
		Url string `json:"url"`
	}
	StripePublic struct {
		Value string `json:"value"`
	}
	IpinfoToken struct {
		Value string `json:"value"`
	}
	AuthFingerprintKey struct {
		Value string `json:"value"`
	}
	SSHKey struct {
		Public  string `json:"public"`
		Private string `json:"private"`
	}
}

var Resource resource

func init() {
	val := reflect.ValueOf(&Resource).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := val.Type().Field(i)
		envVarName := fmt.Sprintf("SST_RESOURCE_%s", typeField.Name)
		envValue, exists := os.LookupEnv(envVarName)
		if !exists {
			// Use default values for local development
			switch typeField.Name {
			case "SSHKey":
				// Try to load SSH keys from common locations
				if privKey, err := loadSSHKey(); err == nil {
					envValue = privKey
				} else {
					envValue = `{"public":"","private":""}`
				}
			case "Api":
				envValue = `{"url":"http://localhost:3000"}`
			case "Auth":
				envValue = `{"url":"http://localhost:3001"}`
			case "StripePublic":
				envValue = `{"value":"pk_test_local"}`
			case "IpinfoToken":
				envValue = `{"value":"test_local"}`
			case "AuthFingerprintKey":
				envValue = `{"value":"test_local"}`
			default:
				panic(fmt.Sprintf("Environment variable %s is required", envVarName))
			}
		}
		if err := json.Unmarshal([]byte(envValue), field.Addr().Interface()); err != nil {
			panic(err)
		}
	}
}

func loadSSHKey() (string, error) {
	// Try to load SSH keys from temp directory
	tempDir := os.TempDir()
	keyPaths := []string{
		filepath.Join(tempDir, "terminal-ssh", "terminal_key"),
		filepath.Join(tempDir, "ssh-terminal", "terminal_key"),
	}

	for _, privKeyPath := range keyPaths {
		privKeyBytes, err := os.ReadFile(privKeyPath)
		if err == nil {
			pubKeyPath := privKeyPath + ".pub"
			pubKeyBytes, err := os.ReadFile(pubKeyPath)
			if err == nil {
				jsonStr := fmt.Sprintf(`{"private":%s,"public":%s}`,
					marshal(string(privKeyBytes)),
					marshal(string(pubKeyBytes)))
				return jsonStr, nil
			}
		}
	}
	return "", fmt.Errorf("no SSH key found")
}

func marshal(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
