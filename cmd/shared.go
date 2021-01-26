package cmd

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/ansible"
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func filterCompletion(toComplete string, path string) ([]string, cobra.ShellCompDirective) {
	// logrus.SetLevel(logrus.PanicLevel)

	key, op, value := ansible.ParseFilter(toComplete)

	cobra.CompDebug(fmt.Sprintf("key:%s op:%s value:%s", key, op, value), true)

	k8s := ansible.Projects.LoadFromPath(path)

	availableKeys := k8s.InventoryKeys()

	blank := key == utils.EmptyString && op == utils.EmptyString && value == utils.EmptyString
	if blank {
		return availableKeys, cobra.ShellCompDirectiveDefault
	}

	writingKey := key != utils.EmptyString && op == utils.EmptyString && value == utils.EmptyString
	if writingKey {
		var keysCompletion []string
		for _, availableKey := range availableKeys {
			if strings.HasPrefix(availableKey, key) {
				keysCompletion = append(keysCompletion, availableKey)
			}
		}

		if len(keysCompletion) == 1 {
			var prefixedOperator []string

			for _, allowedOperator := range ansible.AllowedOperators {
				prefixedOperator = append(prefixedOperator, fmt.Sprintf("%s%s", keysCompletion[0], allowedOperator))
			}
			return prefixedOperator, cobra.ShellCompDirectiveDefault
		}

		return keysCompletion, cobra.ShellCompDirectiveDefault
	}

	writingOp := key != utils.EmptyString && op != utils.EmptyString && value == utils.EmptyString
	if writingOp {
		var prefixedOperator []string

		for _, allowedOperator := range ansible.AllowedOperators {

			if op == allowedOperator {
				availableValues := k8s.InventoryValues(key)

				var prefixedValues []string

				for _, availableValue := range availableValues {

					if availableValue != utils.EmptyString {
						prefixedValues = append(prefixedValues, fmt.Sprintf("%s%s%s", key, op, availableValue))
					}

				}

				return prefixedValues, cobra.ShellCompDirectiveDefault

			}

			if allowedOperator[0] == op[0] {
				prefixedOperator = append(prefixedOperator, fmt.Sprintf("%s%s", key, allowedOperator))
			}

		}

		if len(prefixedOperator) == 1 {
			availableValues := k8s.InventoryValues(key)

			_, foundOp, _ := ansible.ParseFilter(prefixedOperator[0])

			var prefixedValues []string

			for _, availableValue := range availableValues {

				if availableValue != utils.EmptyString {
					prefixedValues = append(prefixedValues, fmt.Sprintf("%s%s%s", key, foundOp, availableValue))
				}

			}

			return prefixedValues, cobra.ShellCompDirectiveDefault
		}

		return prefixedOperator, cobra.ShellCompDirectiveDefault
	}

	writingValue := key != utils.EmptyString && op != utils.EmptyString && value != utils.EmptyString
	if writingValue {
		for _, allowedOperator := range ansible.AllowedOperators {

			if op == allowedOperator {
				availableValues := k8s.InventoryValues(key)

				var prefixedValues []string

				for _, availableValue := range availableValues {
					if availableValue != utils.EmptyString && strings.HasPrefix(availableValue, value) {
						prefixedValues = append(prefixedValues, fmt.Sprintf("%s%s%s", key, op, availableValue))
					}

				}

				return prefixedValues, cobra.ShellCompDirectiveDefault

			}

		}
		return []string{}, cobra.ShellCompDirectiveDefault

	}

	return k8s.InventoryKeys(), cobra.ShellCompDirectiveDefault
}

func playbookCompletion(toComplete string, path string) ([]string, cobra.ShellCompDirective) {
	logrus.SetLevel(logrus.PanicLevel)
	k8s := ansible.Projects.LoadFromPath(path)
	return k8s.PlaybookPaths(), cobra.ShellCompDirectiveDefault
}

func tagsCompletion(toComplete string, path string, playbookPath string) ([]string, cobra.ShellCompDirective) {
	logrus.SetLevel(logrus.PanicLevel)
	var _ = regexp.MustCompile("([\\w-.\\/]+)([,]|)")
	if len(playbookPath) == 0 {
		return nil, cobra.ShellCompDirectiveDefault
	}
	project := ansible.Projects.LoadFromPath(path)

	//TODO unmanaged error
	playbook, _ := project.PlaybookPath(playbookPath)

	return playbook.AllTags().List(), cobra.ShellCompDirectiveDefault
}
