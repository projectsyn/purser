package gcp

import (
	"github.com/go-logr/logr"
	"google.golang.org/api/compute/v1"
)

// Validator represents validator for Google Compute Platform projects.
type Validator struct {
	client              Client
	log                 logr.Logger
	expectedServices    []string
	expectedManagedZone string
	quotaNeeds          map[string]float64
}

// Client defines the interface used to fetch data from the Google APIs.
type Client interface {
	GetManagedZones() (map[string]int, error)
	GetQuotas() ([]*compute.Quota, error)
	GetServiceStates() (map[string]bool, error)
}

// New creates a new validator instance.
func New(c Client, log logr.Logger) *Validator {
	return &Validator{
		client:           c,
		log:              log,
		expectedServices: OpenShift4Required,
		quotaNeeds:       OpenShift4QuotaNeeds,
	}
}

// ValidateAll runs all validations on a GCP project.
func (v *Validator) ValidateAll() (bool, error) {
	validators := []validatorFunction{
		func() (bool, error) {
			return v.ValidateServices(v.expectedServices)
		},
		func() (bool, error) {
			return v.ValidatePublicManagedZones(v.expectedManagedZone)
		},
		func() (bool, error) {
			return v.ValidateQuotas(v.quotaNeeds)
		},
	}

	return validate(validators...)
}

// SetExpectedServices sets the list of services to be expected enabled.
func (v *Validator) SetExpectedServices(expected []string) {
	v.expectedServices = expected
}

// SetExpectedManagedZone sets the domain of a public managed zone to be expected.
func (v *Validator) SetExpectedManagedZone(expected string) {
	v.expectedManagedZone = expected
}

type validatorFunction func() (bool, error)

func validate(validators ...validatorFunction) (bool, error) {
	valid := true
	for _, validator := range validators {
		v, err := validator()
		if err != nil {
			return false, err
		}
		valid = valid && v
	}

	return valid, nil
}
