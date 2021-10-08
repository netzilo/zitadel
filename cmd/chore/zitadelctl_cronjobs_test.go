package chore_test

import "gopkg.in/yaml.v3"

type cronjob struct {
	Metadata struct {
		Name string
	}
	Spec struct {
		Schedule string
	}
}

func getCronJobScheduleWithName(kubectl kubectlCmd, namespace, name string) string {
	cron, err := getCronJobWithName(kubectl, namespace, name)
	if err != nil {
		return ""
	}
	return cron.Spec.Schedule
}

func getCronJobWithName(kubectl kubectlCmd, namespace, name string) (cronjob, error) {
	cronjob := cronjob{}
	args := []string{
		"get", "cronjobs",
		"--namespace", namespace,
		name,
		"--output", "yaml",
	}

	cmd := kubectl(args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return cronjob, err
	}

	if err := yaml.Unmarshal(out, &cronjob); err != nil {
		return cronjob, err
	}
	return cronjob, nil
}
