package main

import (
	"github.com/ca-gip/dploy/cmd"
)

func main() {

	//home, _ := os.UserHomeDir()
	//path := fmt.Sprintf("%s/%s", home, "Projects/ansible-kube/inventories")
	//k8s := services.LoadFromPath(path, path+"/..")
	//
	//filter := []string{"platform==os", "customer!=cacf_corp_hors_prod"}
	//filteredInventories := k8s.FilterFromVars(filter)
	//fmt.Println("Filtering ", len(filteredInventories), "/", len(k8s.Inventories))
	//for _, i := range filteredInventories {
	//	fmt.Println(i.FilePath)
	//}

	//fmt.Println("Playbooks")
	//
	//for _, i := range k8s.Playbooks {
	//	fmt.Println(i.Name, i.Plays)
	//	for _, t := range i.Plays {
	//		fmt.Printf("\ntags:%v", t.Tags)
	//	}
	//}

	cmd.Execute()
}
