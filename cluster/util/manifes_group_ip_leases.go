package util

import (
	mani "github.com/akash-network/akash-api/go/manifest/v2beta2"
	cluster "github.com/akash-network/provider/cluster"
)

func GetLeasedIpsFromDeploy(mgroup *mani.Group) []cluster.ServiceExposeWithServiceName {
	leasedIPs := make([]cluster.ServiceExposeWithServiceName, 0)
	for _, service := range mgroup.Services {
		for _, expose := range service.Expose {
			if expose.Global && len(expose.IP) != 0 {
				v := cluster.ServiceExposeWithServiceName{Expose: expose, Name: service.Name}
				leasedIPs = append(leasedIPs, v)

			}
		}
	}
	return leasedIPs
}
