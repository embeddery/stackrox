package multipliers

import (
	"testing"

	"bitbucket.org/stack-rox/apollo/central/risk/getters"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestViolationsScore(t *testing.T) {
	cases := []struct {
		name     string
		alerts   []*v1.ListAlert
		expected *v1.Risk_Result
	}{
		{
			name:     "No alerts",
			alerts:   nil,
			expected: nil,
		},
		{
			name: "One critical",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_CRITICAL_SEVERITY,
						Name:     "Policy 1",
					},
				},
			},
			expected: &v1.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []string{
					"Policy 1 (severity: Critical)",
				},
				Score: 1.2,
			},
		},
		{
			name: "Two critical",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_CRITICAL_SEVERITY,
						Name:     "Policy 1",
					},
				},
			},
			expected: &v1.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []string{
					"Policy 1 (severity: Critical)",
				},
				Score: 1.2,
			},
		},
		{
			name: "Mix of severities (1)",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_HIGH_SEVERITY,
						Name:     "Policy 1",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_MEDIUM_SEVERITY,
						Name:     "Policy 2",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_LOW_SEVERITY,
						Name:     "Policy 3",
					},
				},
			},
			expected: &v1.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []string{
					"Policy 1 (severity: High)",
					"Policy 2 (severity: Medium)",
					"Policy 3 (severity: Low)",
				},
				Score: 1.3,
			},
		},
		{
			name: "Mix of severities (2)",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_CRITICAL_SEVERITY,
						Name:     "Policy 1",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_HIGH_SEVERITY,
						Name:     "Policy 2",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_LOW_SEVERITY,
						Name:     "Policy 3",
					},
				},
			},
			expected: &v1.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []string{
					"Policy 1 (severity: Critical)",
					"Policy 2 (severity: High)",
					"Policy 3 (severity: Low)",
				},
				Score: 1.4,
			},
		},
		{
			name: "Don't include stale alerts",
			alerts: []*v1.ListAlert{
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_CRITICAL_SEVERITY,
						Name:     "Policy 3",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_HIGH_SEVERITY,
						Name:     "Policy 2",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_LOW_SEVERITY,
						Name:     "Policy 1",
					},
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_CRITICAL_SEVERITY,
						Name:     "Policy Don't Show Me!",
					},
					Stale: true,
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_HIGH_SEVERITY,
						Name:     "Policy Don't Show Me!",
					},
					Stale: true,
				},
				{
					Policy: &v1.ListAlertPolicy{
						Severity: v1.Severity_LOW_SEVERITY,
						Name:     "Policy Don't Show Me!",
					},
					Stale: true,
				},
			},
			expected: &v1.Risk_Result{
				Name: PolicyViolationsHeading,
				Factors: []string{
					"Policy 3 (severity: Critical)",
					"Policy 2 (severity: High)",
					"Policy 1 (severity: Low)",
				},
				Score: 1.4,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mult := NewViolations(&getters.MockAlertsGetter{
				Alerts: c.alerts,
			})
			deployment := getMockDeployment()
			result := mult.Score(deployment)
			assert.Equal(t, c.expected, result)
		})
	}
}
