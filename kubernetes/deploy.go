/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package kubernetes

import (
	"fmt"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

func DeployManifest(pathToKubeconfig, manifest, namespace string) {

	clusterConfig, _ := sthingsK8s.GetKubeConfig(pathToKubeconfig)

	fmt.Println(manifest)

	resourceCreationStatus, err := sthingsK8s.CreateDynamicResourcesFromTemplate(clusterConfig, []byte(manifest), namespace)

	if err != nil {
		fmt.Errorf("FAILED TO CREATE RESOURCE ON CLUSTER: %w", err)
	}

	if resourceCreationStatus {
		fmt.Println("RESOURCE CREATED (OR PATCHED) SUCCESSFULLY")
	}

}
