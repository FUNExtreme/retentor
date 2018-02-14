package main

type RetentionPolicy struct {
	StartsAfterDays int
	ValidForDays    int
	PolicyType      RetentionPolicyType
}

type RetentionPolicyType string

var (
	RetentionPolicyHourly    RetentionPolicyType = "hourly"
	RetentionPolicySixHourly RetentionPolicyType = "sixhourly"
	RetentionPolicyDaily     RetentionPolicyType = "daily"
	RetentionPolicyWeekly    RetentionPolicyType = "weekly"
	RetentionPolicyMonthly   RetentionPolicyType = "monthly"
)

var RetentionPolicies = []RetentionPolicy{
	RetentionPolicy{
		StartsAfterDays: 0,
		ValidForDays:    2,
		PolicyType:      RetentionPolicyHourly,
	},
	RetentionPolicy{
		StartsAfterDays: 2,
		ValidForDays:    5,
		PolicyType:      RetentionPolicySixHourly,
	},
	RetentionPolicy{
		StartsAfterDays: 7,
		ValidForDays:    24,
		PolicyType:      RetentionPolicyDaily,
	},
	RetentionPolicy{
		StartsAfterDays: 31,
		ValidForDays:    334,
		PolicyType:      RetentionPolicyMonthly,
	},
}
