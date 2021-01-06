package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Project struct {
	Name        string
	Inventories []*Inventory
	Playbooks   []*Playbook
}

var filterRegex = regexp.MustCompile("(\\w*)(\\W*)(\\w*)")

func ParseFilter(filter string) (key, op, value string) {
	result := filterRegex.FindStringSubmatch(filter)
	return result[1], result[2], result[3]
}

func ConditionEval(left, right, op string) bool {
	switch op {
	case "==":
		return strings.EqualFold(left, right)
	case "!=":
		return !strings.EqualFold(left, right)
	case "$=":
		return strings.HasSuffix(left, right)
	case "~=":
		return strings.Contains(left, right)
	case "^=":
		return strings.HasPrefix(left, right)
	default:
		log.Fatalf("Unsuported filter operation %s", op)
		return false
	}
}

func AllTrue(a map[string]bool) bool {
	for _, value := range a {
		if !value {
			return false
		}
	}
	return true
}

func (project *Project) FilterFromVars(filters []string) (filtered []*Inventory) {
	for _, inventory := range project.Inventories {
		if inventory.Data != nil {

			type condition = string
			matchFilter := make(map[condition]bool, len(filters))

			for _, filter := range filters {
				key, op, value := ParseFilter(filter)
				if ConditionEval(inventory.Data.Groups["all"].Vars[key], value, op) {
					matchFilter[filter] = true
				}
			}

			if AllTrue(matchFilter) {
				filtered = append(filtered, inventory)
			}

		}
	}
	return
}

// TODO: Add assert on file system ( readable, permissions ...)
func LoadFromPath(inventoryPath string, playbookPath string) (project Project) {
	project = Project{
		Name:      inventoryPath,
		Playbooks: nil,
	}
	fmt.Println(playbookPath)
	playbooks, errPlaybooks := readPlaybook(playbookPath)
	inventories, errInventories := readInventories(inventoryPath)
	project.Playbooks = playbooks
	project.Inventories = inventories

	if errPlaybooks != nil {
		log.Fatalln("Cannot parse directory for playbooks: ", errPlaybooks.Error())
	}
	if errInventories != nil {
		log.Fatalln("Cannot parse directory for inventories: ", errInventories.Error())
	}
	for _, inventory := range inventories {
		inventory.make()
	}
	return
}
