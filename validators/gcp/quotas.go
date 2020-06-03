package gcp

var (
	// OpenShift4QuotaNeeds holds the quota needs for an OpenShift 4 project.
	// See https://docs.openshift.com/container-platform/4.4/installing/installing_gcp/installing-gcp-account.html#installation-gcp-limits_installing-gcp-account.
	OpenShift4QuotaNeeds = map[string]float64{
		"CPUS_ALL_REGIONS": 28,
		"FIREWALLS":        11,
		"FORWARDING_RULES": 2,
		"HEALTH_CHECKS":    3,
		"IMAGES":           1,
		"IN_USE_ADDRESSES": 4,
		"NETWORKS":         2,
		"ROUTERS":          1,
		"ROUTES":           2,
		"STATIC_ADDRESSES": 4,
		"SUBNETWORKS":      2,
		"TARGET_POOLS":     3,
	}
)

// ValidateQuotas checks if quota limits are not exceeded.
// The given list of needed reserves are taken into consideration.
func (v *Validator) ValidateQuotas(needs map[string]float64) (bool, error) {
	quotas, err := v.client.GetQuotas()
	if err != nil {
		return false, err
	}

	valid := true
	for _, quota := range quotas {
		need := needOrDefault(quota.Metric, needs)
		if quota.Usage+need > quota.Limit {
			v.log.Info(
				"Quota exceeded",
				"metric", quota.Metric,
				"need", need,
				"usage", quota.Usage,
				"limit", quota.Limit,
			)
			valid = false
		}
	}

	return valid, nil
}

func needOrDefault(name string, needs map[string]float64) float64 {
	if need, ok := needs[name]; ok {
		return need
	}

	return 0
}
