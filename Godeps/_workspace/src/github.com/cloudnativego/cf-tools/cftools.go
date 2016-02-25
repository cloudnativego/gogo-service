package cftools

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfenv"
)

//GetVCAPServiceProperty retrieves a property from bound service credentials.
func GetVCAPServiceProperty(serviceName string, propertyName string, appEnv *cfenv.App) (propertyValue string, err error) {
	if propertyName == "" {
		return "", fmt.Errorf("Must supply a property name value.")
	}

	service, err := getVCAPService(serviceName, appEnv)
	if err != nil {
		return "", err
	}

	propertyValue = service.Credentials[propertyName].(string)
	if propertyValue == "" {
		return "", fmt.Errorf("Error retrieving property %v", propertyName)
	}
	return
}

func getVCAPService(serviceName string, appEnv *cfenv.App) (service *cfenv.Service, err error) {
	service, err = appEnv.Services.WithName(serviceName)
	return
}
