package proxmoxve

import (
	"github.com/yuuki0xff/clustertest/config"
	"github.com/yuuki0xff/clustertest/models"
)

func init() {
	config.SpecInitializers[models.SpecType("proxmox-ve")] = func() models.Spec { return &ProxmoxVESpec{} }
}

type ProxmoxVESpec struct {
	// Proxmox VE account settings..
	Proxmox *struct {
		// IP address or FQDN of the Proxmox VE API server.
		Address string
		Account struct {
			User     string
			Password string
		}
		// Fingerprint of the Proxmox VE API server.
		// If you need the server public key pinning to make it more secure.
		// TODO: https://medium.com/@zmanian/server-public-key-pinning-in-go-7a57bbe39438
		Fingerprint string
	}
	// Addresses to assign to VMs.
	AddressPools []*struct {
		StartAddress string `yaml:"start_address"`
		EndAddress   string `yaml:"end_address"`
		CIDR         int
		Gateway      string
	}
	// User information.
	// This user will create by cloud-init at VM start-up.
	User *struct {
		User         string
		Password     string
		SSHPublicKey string `yaml:"ssh_public_key"`
	}
	VMs []*struct {
		// Template name.
		Template string
		// Number of VMs.
		Nodes int
		// Number of processors.
		Processors int
		// RAM size (GiB).
		MemorySize int `yaml:"memory_size"`
		// Minimal storage size (GiB).
		// The storage may be large than specified size.
		StorageSize int `yaml:"storage_size"`
		// Define tasks to execute on VMs.
		Scripts *config.ScriptConfig
	}
}

func (s *ProxmoxVESpec) String() string {
	return "<ProxmoxVESpec>"
}
func (s *ProxmoxVESpec) Type() models.SpecType {
	return models.SpecType("proxmox-ve")
}

// TODO: write the provisioning logic.
