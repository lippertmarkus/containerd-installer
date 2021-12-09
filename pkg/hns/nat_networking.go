package hns

import (
	"encoding/json"
	"github.com/Microsoft/hcsshim"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type NatNetworkDetails struct {
	AddressPrefix  string
	GatewayAddress string
}

func GetOrCreateNatNetwork() (NatNetworkDetails, error) {
	logrus.Debug("Check for existing HNS NAT network")
	var natNetworkDetails NatNetworkDetails

	existingNetwork, err := hcsshim.GetHNSNetworkByName("nat")
	if err == nil {
		natNetworkDetails = NatNetworkDetails{
			AddressPrefix:  existingNetwork.Subnets[0].AddressPrefix,
			GatewayAddress: existingNetwork.Subnets[0].GatewayAddress,
		}
		return natNetworkDetails, nil
	}

	logrus.Debug("No existing NAT network can be retrieved, creating a new one")
	natNetwork := hcsshim.HNSNetwork{
		Name: "nat",
		Type: "nat",
	}
	natNetworkConfig, err := json.Marshal(natNetwork)
	if err != nil {
		return natNetworkDetails, errors.Wrap(err, "Couldn't marshal HNS network object to JSON")
	}
	createdNetwork, err := hcsshim.HNSNetworkRequest("POST", "", string(natNetworkConfig))
	if err != nil {
		return natNetworkDetails, errors.Wrap(err, "NAT network couldn't be created")
	}

	natNetworkDetails = NatNetworkDetails{
		AddressPrefix:  createdNetwork.Subnets[0].AddressPrefix,
		GatewayAddress: createdNetwork.Subnets[0].GatewayAddress,
	}

	logrus.Debugf("Created NAT network: gateway: %s, subnet: %s", natNetworkDetails.GatewayAddress, natNetworkDetails.AddressPrefix)

	return natNetworkDetails, nil
}
