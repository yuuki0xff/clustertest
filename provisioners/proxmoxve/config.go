package proxmoxve

import (
	"fmt"
	"github.com/yuuki0xff/clustertest/config"
	"github.com/yuuki0xff/clustertest/models"
	"net"
	"net/url"
)

const PveSpecType models.SpecType = "proxmox-ve"

func init() {
	config.SpecInitializers[PveSpecType] = func() models.Spec { return &PveSpec{} }
}

type PveSpec struct {
	// Identifier of the spec.
	Name string
	// Proxmox VE account settings..
	Proxmox struct {
		// URL of the Proxmox VE API server.
		// Example: https://pve.local:8006
		Address string
		// Login information.
		Account struct {
			User     string
			Password string
		}
		// Fingerprint of the Proxmox VE API server.
		// If you need the server certificate pinning to make it more secure.
		Fingerprint string
	}
	// Addresses to assign to VMs.
	AddressPools []*struct {
		StartAddress string `yaml:"start_address"`
		EndAddress   string `yaml:"end_address"`
		CIDR         int
		Gateway      string
	} `yaml:"address_pools"`
	// User information.
	// This user will create by cloud-init at VM start-up.
	User struct {
		User string
		// todo: password is not used.
		Password     string
		SSHPublicKey string `yaml:"ssh_public_key"`
	}
	VMs map[string]*PveVM
}
type PveVM struct {
	// Template name.
	Template string
	// Pool name to which the VM belongs.
	Pool string
	// Number of VMs.
	Nodes int
	// Number of processors.
	Processors int
	// RAM size (MiB).
	MemorySize int `yaml:"memory_size"`
	// Minimal storage size (GiB).
	// The storage may be large than specified size.
	StorageSize int `yaml:"storage_size"`
	// Define tasks to execute on VMs.
	Scripts *config.ScriptConfigSet
}

func (s *PveSpec) String() string {
	return "<PveSpec>"
}
func (s *PveSpec) Type() models.SpecType {
	return PveSpecType
}
func (s *PveSpec) LoadDefault(defaultSpec models.Spec) error {
	default_ := defaultSpec.(*PveSpec)
	empty := PveSpec{}

	if s.Proxmox == empty.Proxmox {
		s.Proxmox = default_.Proxmox
	}
	if len(s.AddressPools) == 0 {
		s.AddressPools = default_.AddressPools
	}
	if s.User == empty.User {
		s.User = default_.User
	}
	return nil
}
func (s *PveSpec) validateCommon() *config.Validator {
	v := &config.Validator{}
	empty := PveSpec{}

	v.Validate("name", s.Name != "").SetReason("name must not be empty")
	if s.Proxmox != empty.Proxmox {
		addr, err := url.Parse(s.Proxmox.Address)
		v.ValidateError("proxmox.url", err)
		if addr != nil {
			v.Validate("proxmox.url", addr.Scheme == "https").SetReason(`schema must be "https""`)
			v.Validate("proxmox.url", addr.Hostname() != "").SetReason("host name must not be empty")
		}

		if s.Proxmox.Account != empty.Proxmox.Account {
			v.Validate("proxmox.account.user", s.Proxmox.Account.User != "").SetReason("user must not be empty")
			v.Validate("proxmox.account.password", s.Proxmox.Account.Password != "").SetReason("password must not be empty")
		}
	}

	for idx, addrRange := range s.AddressPools {
		base := fmt.Sprintf("address_pools[%d].", idx)
		addr := net.ParseIP(addrRange.StartAddress)
		v.Validate(base+"start_address", addr != nil).SetReason("invalid IP address")

		addr = net.ParseIP(addrRange.EndAddress)
		v.Validate(base+"end_address", addr != nil).SetReason("invalid IP address")

		addr = net.ParseIP(addrRange.Gateway)
		v.Validate(base+"gateway", addr != nil).SetReason("invalid IP address")
	}

	if s.User != empty.User {
		v.Validate("user.user", s.User.User != "").SetReason("user must not be empty")
		v.Validate("user.ssh_public_key", s.User.SSHPublicKey != "").SetReason("ssh_public_key must not be empty")
	}
	return v
}
func (s *PveSpec) Validate() error {
	v := s.validateCommon()
	for _, vm := range s.VMs {
		vm.validate(v)
	}
	return v.Error()
}
func (s *PveSpec) ValidateDefault() error {
	v := s.validateCommon()
	v.Validate("vms", len(s.VMs) == 0).SetReason("vms must not define in default settings")
	return v.Error()
}

func (vm PveVM) validate(v *config.Validator) {
	v.Validate("template", vm.Template != "").SetReason("template is empty")
	v.Validate("pool", vm.Template != "").SetReason("pool is empty")
	v.Validate("nodes", vm.Nodes > 0).SetReason("nodes must be greater than 0")
	v.Validate("processors", vm.Processors > 0).SetReason("processors must be greater than 0")
	v.Validate("memory_size", vm.MemorySize > 0).SetReason("memory_size must be greater than 0")
	v.Validate("storage_size", vm.StorageSize > 0).SetReason("storage_size must be greater than 0")

	if s := vm.Scripts.Before.Get(); s != nil {
		v.Merge("scripts.before", s.Validate())
	}
	if s := vm.Scripts.Main.Get(); s != nil {
		v.Merge("scripts.main", s.Validate())
	}
	if s := vm.Scripts.After.Get(); s != nil {
		v.Merge("scripts.after", s.Validate())
	}
}
