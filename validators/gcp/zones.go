package gcp

import "github.com/projectsyn/probatorem/clients/googleapis"

// ValidatePublicManagedZones checks there is at least one public zone and if
// the expected domain is present.
func (v *Validator) ValidatePublicManagedZones(expected string) (bool, error) {
	zones, err := v.client.GetManagedZones()
	if err != nil {
		return false, err
	}

	count := 0
	found := len(expected) == 0

	for zone, visibilty := range zones {
		if visibilty == googleapis.ManagedZoneVisibilityPublic {
			count += 1
		}
		found = found || zone == expected
	}

	if count < 1 {
		v.log.Info("The project must have at least one public DNS zone.")
	}

	if !found {
		v.log.Info("DNS zone not managed by project.", "domain", expected)
	}

	return count > 0 && found, nil
}
