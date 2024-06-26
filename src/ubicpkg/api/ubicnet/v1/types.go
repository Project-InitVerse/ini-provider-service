package v1

import (
	"math"
	"strconv"

	ctypes "providerService/src/cluster/types/v1"
	clusterutil "providerService/src/cluster/ubicutil"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	maniv2beta1 "github.com/ovrclk/akash/manifest/v2beta1"
	types "github.com/ovrclk/akash/types/v1beta2"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Manifest store metadata, specifications and status of the Lease
type Manifest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   ManifestSpec   `json:"spec,omitempty"`
	Status ManifestStatus `json:"status,omitempty"`
}

// ManifestStatus stores state and message of manifest
type ManifestStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// ManifestSpec stores LeaseID, Group and metadata details
type ManifestSpec struct {
	LeaseID LeaseID       `json:"lease_id"`
	Group   ManifestGroup `json:"group"`
}

// Deployment returns the cluster.Deployment that the saved manifest represents.
func (m Manifest) Deployment() (ctypes.Deployment, error) {
	lid, err := m.Spec.LeaseID.ToIniType()
	if err != nil {
		return nil, err
	}

	group, err := m.Spec.Group.toIni()
	if err != nil {
		return nil, err
	}
	return deployment{lid: lid, group: group}, nil
}

type deployment struct {
	lid   ctypes.LeaseID
	group maniv2beta1.Group
}

func (d deployment) LeaseID() ctypes.LeaseID {
	return d.lid
}

func (d deployment) ManifestGroup() maniv2beta1.Group {
	return d.group
}

// NewManifest creates new manifest with provided details. Returns error in case of failure.
func NewManifest(ns string, lid ctypes.LeaseID, mgroup *maniv2beta1.Group) (*Manifest, error) {
	group, err := manifestGroupFromIni(mgroup)
	if err != nil {
		return nil, err
	}

	return &Manifest{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Manifest",
			APIVersion: "ini.net/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterutil.LeaseIDToNamespace(lid),
			Namespace: ns,
		},
		Spec: ManifestSpec{
			Group:   group,
			LeaseID: LeaseIDFromInitype(lid),
		},
	}, nil
}

// LeaseID stores deployment, group sequence, order, provider and metadata
type LeaseID struct {
	Owner    string `json:"owner"`
	OSeq     uint64 `json:"oseq"`
	Provider string `json:"provider"`
}

// ToIniType returns LeaseID from LeaseID details
func (id LeaseID) ToIniType() (ctypes.LeaseID, error) {
	return ctypes.LeaseID{
		Owner:    id.Owner,
		OSeq:     id.OSeq,
		Provider: id.Provider,
	}, nil
}

// LeaseIDFromInitype returns LeaseID instance from ini
func LeaseIDFromInitype(id ctypes.LeaseID) LeaseID {
	return LeaseID{
		Owner:    id.Owner,
		OSeq:     id.OSeq,
		Provider: id.Provider,
	}
}

// ManifestGroup stores metadata, name and list of SDL manifest services
type ManifestGroup struct {
	// Placement profile name
	Name string `json:"name,omitempty"`
	// Service definitions
	Services []ManifestService `json:"services,omitempty"`
}

// toIni returns ini group details formatted from manifest group
func (m ManifestGroup) toIni() (maniv2beta1.Group, error) {
	am := maniv2beta1.Group{
		Name:     m.Name,
		Services: make([]maniv2beta1.Service, 0, len(m.Services)),
	}

	for _, svc := range m.Services {
		asvc, err := svc.toIni()
		if err != nil {
			return am, err
		}
		am.Services = append(am.Services, asvc)
	}

	return am, nil
}

// manifestGroupFromIni returns manifest group instance from ini group
func manifestGroupFromIni(m *maniv2beta1.Group) (ManifestGroup, error) {
	ma := ManifestGroup{
		Name:     m.Name,
		Services: make([]ManifestService, 0, len(m.Services)),
	}

	for _, svc := range m.Services {
		service, err := manifestServiceFromIni(svc)
		if err != nil {
			return ManifestGroup{}, err
		}

		ma.Services = append(ma.Services, service)
	}

	return ma, nil
}

// ManifestStorageParams is struct
type ManifestStorageParams struct {
	Name     string `json:"name" yaml:"name"`
	Mount    string `json:"mount" yaml:"mount"`
	ReadOnly bool   `json:"readOnly" yaml:"readOnly"`
}

// ManifestServiceParams is struct
type ManifestServiceParams struct {
	Storage []ManifestStorageParams `json:"storage,omitempty"`
}

// ManifestService stores name, image, args, env, unit, count and expose list of service
type ManifestService struct {
	// Service name
	Name string `json:"name,omitempty"`
	// Docker image
	Image   string   `json:"image,omitempty"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	Env     []string `json:"env,omitempty"`
	// Resource requirements
	// in current version of CRD it is named as unit
	Resources ResourceUnits `json:"unit"`
	// Number of instances
	Count uint32 `json:"count,omitempty"`
	// Overlay Network Links
	Expose []ManifestServiceExpose `json:"expose,omitempty"`
	// Miscellaneous service parameters
	Params *ManifestServiceParams `json:"params,omitempty"`
}

func (ms ManifestService) toIni() (maniv2beta1.Service, error) {
	res, err := ms.Resources.toIni()
	if err != nil {
		return maniv2beta1.Service{}, err
	}

	ams := &maniv2beta1.Service{
		Name:      ms.Name,
		Image:     ms.Image,
		Command:   ms.Command,
		Args:      ms.Args,
		Env:       ms.Env,
		Resources: res,
		Count:     ms.Count,
		Expose:    make([]maniv2beta1.ServiceExpose, 0, len(ms.Expose)),
	}

	for _, expose := range ms.Expose {
		value, err := expose.toIni()
		if err != nil {
			return maniv2beta1.Service{}, err
		}
		ams.Expose = append(ams.Expose, value)

		if len(value.IP) != 0 {
			res.Endpoints = append(res.Endpoints, types.Endpoint{
				Kind:           types.Endpoint_LEASED_IP,
				SequenceNumber: value.EndpointSequenceNumber,
			})
		}
	}

	if ms.Params != nil {
		ams.Params = &maniv2beta1.ServiceParams{
			Storage: make([]maniv2beta1.StorageParams, 0, len(ms.Params.Storage)),
		}

		for _, storage := range ms.Params.Storage {
			ams.Params.Storage = append(ams.Params.Storage, maniv2beta1.StorageParams{
				Name:     storage.Name,
				Mount:    storage.Mount,
				ReadOnly: storage.ReadOnly,
			})
		}
	}

	return *ams, nil
}

func manifestServiceFromIni(ams maniv2beta1.Service) (ManifestService, error) {
	resources, err := resourceUnitsFromIni(ams.Resources)
	if err != nil {
		return ManifestService{}, err
	}

	ms := ManifestService{
		Name:      ams.Name,
		Image:     ams.Image,
		Command:   ams.Command,
		Args:      ams.Args,
		Env:       ams.Env,
		Resources: resources,
		Count:     ams.Count,
		Expose:    make([]ManifestServiceExpose, 0, len(ams.Expose)),
	}

	for _, expose := range ams.Expose {
		ms.Expose = append(ms.Expose, manifestServiceExposeFromIni(expose))
	}

	if ams.Params != nil {
		ms.Params = &ManifestServiceParams{
			Storage: make([]ManifestStorageParams, 0, len(ams.Params.Storage)),
		}

		for _, storage := range ams.Params.Storage {
			ms.Params.Storage = append(ms.Params.Storage, ManifestStorageParams{
				Name:     storage.Name,
				Mount:    storage.Mount,
				ReadOnly: storage.ReadOnly,
			})
		}
	}

	return ms, nil
}

// ManifestServiceExpose stores exposed ports and accepted hosts details
type ManifestServiceExpose struct {
	Port         uint16 `json:"port,omitempty"`
	ExternalPort uint16 `json:"external_port,omitempty"`
	Proto        string `json:"proto,omitempty"`
	Service      string `json:"service,omitempty"`
	Global       bool   `json:"global,omitempty"`
	// accepted hostnames
	Hosts                  []string                         `json:"hosts,omitempty"`
	HTTPOptions            ManifestServiceExposeHTTPOptions `json:"http_options,omitempty"`
	IP                     string                           `json:"ip,omitempty"`
	EndpointSequenceNumber uint32                           `json:"endpoint_sequence_number"`
}

// ManifestServiceExposeHTTPOptions is struct
type ManifestServiceExposeHTTPOptions struct {
	MaxBodySize uint32   `json:"max_body_size,omitempty"`
	ReadTimeout uint32   `json:"read_timeout,omitempty"`
	SendTimeout uint32   `json:"send_timeout,omitempty"`
	NextTries   uint32   `json:"next_tries,omitempty"`
	NextTimeout uint32   `json:"next_timeout,omitempty"`
	NextCases   []string `json:"next_cases,omitempty"`
}

func (mse ManifestServiceExpose) toIni() (maniv2beta1.ServiceExpose, error) {
	proto, err := maniv2beta1.ParseServiceProtocol(mse.Proto)
	if err != nil {
		return maniv2beta1.ServiceExpose{}, err
	}
	return maniv2beta1.ServiceExpose{
		Port:                   mse.Port,
		ExternalPort:           mse.ExternalPort,
		Proto:                  proto,
		Service:                mse.Service,
		Global:                 mse.Global,
		Hosts:                  mse.Hosts,
		EndpointSequenceNumber: mse.EndpointSequenceNumber,
		IP:                     mse.IP,
		HTTPOptions: maniv2beta1.ServiceExposeHTTPOptions{
			MaxBodySize: mse.HTTPOptions.MaxBodySize,
			ReadTimeout: mse.HTTPOptions.ReadTimeout,
			SendTimeout: mse.HTTPOptions.SendTimeout,
			NextTries:   mse.HTTPOptions.NextTries,
			NextTimeout: mse.HTTPOptions.NextTimeout,
			NextCases:   mse.HTTPOptions.NextCases,
		},
	}, nil
}

// DetermineExposedExternalPort is function
func (mse ManifestServiceExpose) DetermineExposedExternalPort() uint16 {
	if mse.ExternalPort == 0 {
		return mse.Port
	}
	return mse.ExternalPort
}

func manifestServiceExposeFromIni(amse maniv2beta1.ServiceExpose) ManifestServiceExpose {
	return ManifestServiceExpose{
		Port:                   amse.Port,
		ExternalPort:           amse.ExternalPort,
		Proto:                  amse.Proto.ToString(),
		Service:                amse.Service,
		Global:                 amse.Global,
		Hosts:                  amse.Hosts,
		IP:                     amse.IP,
		EndpointSequenceNumber: amse.EndpointSequenceNumber,
		HTTPOptions: ManifestServiceExposeHTTPOptions{
			MaxBodySize: amse.HTTPOptions.MaxBodySize,
			ReadTimeout: amse.HTTPOptions.ReadTimeout,
			SendTimeout: amse.HTTPOptions.SendTimeout,
			NextTries:   amse.HTTPOptions.NextTries,
			NextTimeout: amse.HTTPOptions.NextTimeout,
			NextCases:   amse.HTTPOptions.NextCases,
		},
	}
}

// ManifestServiceStorage stores name and size
type ManifestServiceStorage struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

// ResourceUnits stores cpu, memory and storage details
type ResourceUnits struct {
	CPU     uint32                   `json:"cpu,omitempty"`
	Memory  string                   `json:"memory,omitempty"`
	Storage []ManifestServiceStorage `json:"storage,omitempty"`
}

func (ru ResourceUnits) toIni() (types.ResourceUnits, error) {
	memory, err := strconv.ParseUint(ru.Memory, 10, 64)
	if err != nil {
		return types.ResourceUnits{}, err
	}

	storage := make([]types.Storage, 0, len(ru.Storage))
	for _, st := range ru.Storage {
		size, err := strconv.ParseUint(st.Size, 10, 64)
		if err != nil {
			return types.ResourceUnits{}, err
		}

		storage = append(storage, types.Storage{
			Name:     st.Name,
			Quantity: types.NewResourceValue(size),
		})
	}

	return types.ResourceUnits{
		CPU: &types.CPU{
			Units: types.NewResourceValue(uint64(ru.CPU)),
		},
		Memory: &types.Memory{
			Quantity: types.NewResourceValue(memory),
		},
		Storage: storage,
	}, nil
}

func resourceUnitsFromIni(aru types.ResourceUnits) (ResourceUnits, error) {
	res := ResourceUnits{}
	if aru.CPU != nil {
		if aru.CPU.Units.Value() > math.MaxUint32 {
			return ResourceUnits{}, errors.New("k8s api: cpu units value overflows uint32")
		}
		res.CPU = uint32(aru.CPU.Units.Value())
	}
	if aru.Memory != nil {
		res.Memory = strconv.FormatUint(aru.Memory.Quantity.Value(), 10)
	}

	res.Storage = make([]ManifestServiceStorage, 0, len(aru.Storage))
	for _, storage := range aru.Storage {
		res.Storage = append(res.Storage, ManifestServiceStorage{
			Name: storage.Name,
			Size: strconv.FormatUint(storage.Quantity.Value(), 10),
		})
	}

	return res, nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManifestList stores metadata and items list of manifest
type ManifestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Manifest `json:"items"`
}

// ProviderHost stores spec and status
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProviderHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   ProviderHostSpec   `json:"spec,omitempty"`
	Status ProviderHostStatus `json:"status,omitempty"`
}

// ProviderHostStatus stores state msg
type ProviderHostStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// ProviderHostSpec is struct
type ProviderHostSpec struct {
	Owner        string `json:"owner"`
	Provider     string `json:"provider"`
	Hostname     string `json:"hostname"`
	Oseq         uint64 `json:"oseq"`
	ServiceName  string `json:"service_name"`
	ExternalPort uint32 `json:"external_port"`
}

// ProviderHostList is struct
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProviderHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ProviderHost `json:"items"`
}

// ProviderLeasedIP is struct
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProviderLeasedIP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   ProviderLeasedIPSpec   `json:"spec,omitempty"`
	Status ProviderLeasedIPStatus `json:"status,omitempty"`
}

// ProviderLeasedIPList is struct
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProviderLeasedIPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ProviderLeasedIP `json:"items"`
}

// ProviderLeasedIPStatus is struct
type ProviderLeasedIPStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// ProviderLeasedIPSpec is struct
type ProviderLeasedIPSpec struct {
	LeaseID      LeaseID `json:"lease_id"`
	ServiceName  string  `json:"service_name"`
	Port         uint32  `json:"port"`
	ExternalPort uint32  `json:"external_port"`
	SharingKey   string  `json:"sharing_key"`
	Protocol     string  `json:"protocol"`
}
