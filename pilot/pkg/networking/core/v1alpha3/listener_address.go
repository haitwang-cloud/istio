// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha3

import (
	"istio.io/istio/pilot/pkg/model"
)

const (
	// WildcardAddress binds to all IP addresses
	WildcardAddress = "0.0.0.0"

	// WildcardIPv6Address binds to all IPv6 addresses
	WildcardIPv6Address = "::"

	// LocalhostAddress for local binding
	LocalhostAddress = "127.0.0.1"

	// LocalhostIPv6Address for local binding
	LocalhostIPv6Address = "::1"

	// 6 is the magical number for inbound: 15006, 127.0.0.6, ::6
	InboundPassthroughBindIpv4 = "127.0.0.6"
	InboundPassthroughBindIpv6 = "::6"
)

var (
	wildCards  map[model.IPMode][]string
	localHosts map[model.IPMode][]string
)

// TODO: getActualWildcardAndLocalHost would be removed once the dual stack support in Istio
// getActualWildcardAndLocalHost will return corresponding Wildcard and LocalHost
// depending on value of proxy's IPAddresses.
func getActualWildcardAndLocalHost(node *model.Proxy) (string, string) {
	if node.SupportsIPv4() {
		return WildcardAddress, LocalhostAddress
	}
	return WildcardIPv6Address, LocalhostIPv6Address
}

func getPassthroughBindIP(node *model.Proxy) string {
	if node.SupportsIPv4() {
		return InboundPassthroughBindIpv4
	}
	return InboundPassthroughBindIpv6
}

// getSidecarInboundBindIPs returns the IP that the proxy can bind to along with the sidecar specified port.
// It looks for an unicast address, if none found, then the default wildcard address is used.
// This will make the inbound listener bind to instance_ip:port instead of 0.0.0.0:port where applicable.
func getSidecarInboundBindIPs(node *model.Proxy) []string {
	// Return the IP if its a global unicast address.
	if len(node.GlobalUnicastIP) > 0 {
		return []string{node.GlobalUnicastIP}
	}
	defaultInboundIPs, _ := getWildcardsAndLocalHostForDualStack(node.GetIPMode())
	return defaultInboundIPs
}

func getWildcardsAndLocalHostForDualStack(ipMode model.IPMode) ([]string, []string) {
	return wildCards[ipMode], localHosts[ipMode]
}

func init() {
	// maintain 2 maps to return wildCards and localHosts according to IP mode of proxy
	wildCards = make(map[model.IPMode][]string)
	localHosts = make(map[model.IPMode][]string)

	wildCards[model.IPv4] = []string{WildcardAddress}
	wildCards[model.IPv6] = []string{WildcardIPv6Address}
	wildCards[model.Dual] = []string{WildcardAddress, WildcardIPv6Address}

	localHosts[model.IPv4] = []string{LocalhostAddress}
	localHosts[model.IPv6] = []string{LocalhostIPv6Address}
	localHosts[model.Dual] = []string{LocalhostAddress, LocalhostIPv6Address}
}
