package option

import "net/netip"

type DNSOptions struct {
	Servers         []DNSServerOptions          `json:"servers,omitempty"`
	Rules           []DNSRule                   `json:"rules,omitempty"`
	Final           string                      `json:"final,omitempty"`
	ReverseMapping  bool                        `json:"reverse_mapping,omitempty"`
	MappingOverride bool                        `json:"mapping_override,omitempty"`
	Hosts           map[string]Listable[string] `json:"hosts,omitempty"`
	FakeIP          *DNSFakeIPOptions           `json:"fakeip,omitempty"`
	DNSClientOptions
}

type DNSServerOptions struct {
	Tag                  string           `json:"tag,omitempty"`
	Address              Listable[string] `json:"address"`
	AddressResolver      string           `json:"address_resolver,omitempty"`
	AddressStrategy      DomainStrategy   `json:"address_strategy,omitempty"`
	AddressFallbackDelay Duration         `json:"address_fallback_delay,omitempty"`
	Strategy             DomainStrategy   `json:"strategy,omitempty"`
	Detour               string           `json:"detour,omitempty"`
	ClientSubnet         *AddrPrefix      `json:"client_subnet,omitempty"`
	Insecure             bool             `json:"insecure,omitempty"`
}

type DNSClientOptions struct {
	Strategy         DomainStrategy `json:"strategy,omitempty"`
	DisableCache     bool           `json:"disable_cache,omitempty"`
	DisableExpire    bool           `json:"disable_expire,omitempty"`
	IndependentCache bool           `json:"independent_cache,omitempty"`
	ClientSubnet     *AddrPrefix    `json:"client_subnet,omitempty"`
}

type DNSFakeIPOptions struct {
	Enabled    bool          `json:"enabled,omitempty"`
	Inet4Range *netip.Prefix `json:"inet4_range,omitempty"`
	Inet6Range *netip.Prefix `json:"inet6_range,omitempty"`
}

type DOHInboundOptions struct {
	Network     NetworkList    `json:"network,omitempty"`
	Listen      *ListenAddress `json:"listen,omitempty"`
	ListenPort  uint16         `json:"listen_port,omitempty"`
	QueryPath   string         `json:"query_path,omitempty"`
	UDPFragment *bool          `json:"udp_fragment,omitempty"`
	InboundTLSOptionsContainer
}

type DOQInboundOptions struct {
	Listen           *ListenAddress `json:"listen,omitempty"`
	ListenPort       uint16         `json:"listen_port,omitempty"`
	ZeroRTTHandshake bool           `json:"zero_rtt_handshake,omitempty"`
	UDPFragment      *bool          `json:"udp_fragment,omitempty"`
	InboundTLSOptionsContainer
}
