package gcp

var (
	OpenShift4Required = []string{
		"compute.googleapis.com",
		"cloudapis.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"dns.googleapis.com",
		"iamcredentials.googleapis.com",
		"iam.googleapis.com",
		"servicemanagement.googleapis.com",
		"serviceusage.googleapis.com",
		"storage-api.googleapis.com",
		"storage-component.googleapis.com",
	}
)

// ValidateService validates the services expected to be enabled.
func (v *Validator) ValidateServices(expected []string) (bool, error) {
	states, err := v.client.GetServiceStates()
	if err != nil {
		return false, err
	}

	isValid := true
	for _, name := range expected {
		if state, ok := states[name]; !ok || !state {
			v.log.Info("Required service is not enabled.", "service", name)
			isValid = false
		}
	}

	return isValid, nil
}
