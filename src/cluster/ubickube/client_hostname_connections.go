package kube

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strings"

	crd "providerService/src/ubicpkg/api/ubicnet/v1"

	ctypes "providerService/src/cluster/types/v1"

	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/pager"

	"providerService/src/cluster/ubickube/builder"
)

type hostnameResourceEvent struct {
	eventType ctypes.ProviderResourceEvent
	hostname  string

	owner        common.Address
	oseq         uint64
	provider     common.Address
	serviceName  string
	externalPort uint32
}

func (c *client) DeclareHostname(ctx context.Context, lID ctypes.LeaseID, host string, serviceName string, externalPort uint32) error {
	// Label each entry with the standard labels
	labels := map[string]string{
		builder.UbicManagedLabelName: "true",
	}
	builder.AppendLeaseLabels(lID, labels)

	foundEntry, err := c.ac.UbicnetV1().ProviderHosts(c.ns).Get(ctx, host, metav1.GetOptions{})
	exists := true
	var resourceVersion string

	if err != nil {
		if kubeErrors.IsNotFound(err) {
			exists = false
		} else {
			return err
		}
	} else {
		resourceVersion = foundEntry.ObjectMeta.ResourceVersion
	}

	obj := crd.ProviderHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:            host, // Name is always the hostname, to prevent duplicates
			Labels:          labels,
			ResourceVersion: resourceVersion,
		},
		Spec: crd.ProviderHostSpec{
			Hostname:     host,
			Owner:        lID.Owner,
			Oseq:         lID.OSeq,
			Provider:     lID.Provider,
			ServiceName:  serviceName,
			ExternalPort: externalPort,
		},
		Status: crd.ProviderHostStatus{},
	}

	c.log.Info("declaring hostname", "lease", lID, "service-name", serviceName, "external-port", externalPort, "host", host)
	// Create or update the entry
	if exists {
		_, err = c.ac.UbicnetV1().ProviderHosts(c.ns).Update(ctx, &obj, metav1.UpdateOptions{})
	} else {
		obj.ResourceVersion = ""
		_, err = c.ac.UbicnetV1().ProviderHosts(c.ns).Create(ctx, &obj, metav1.CreateOptions{})
	}
	return err
}

func (c *client) PurgeDeclaredHostname(ctx context.Context, lID ctypes.LeaseID, hostname string) error {
	labelSelector := &strings.Builder{}
	kubeSelectorForLease(labelSelector, lID)

	return c.ac.UbicnetV1().ProviderHosts(c.ns).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: labelSelector.String(),
		FieldSelector: fmt.Sprintf("metadata.name=%s", hostname),
	})
}

func (c *client) PurgeDeclaredHostnames(ctx context.Context, lID ctypes.LeaseID) error {
	labelSelector := &strings.Builder{}
	kubeSelectorForLease(labelSelector, lID)
	result := c.ac.UbicnetV1().ProviderHosts(c.ns).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: labelSelector.String(),
	})

	return result
}

func (ev hostnameResourceEvent) GetLeaseID() ctypes.LeaseID {
	return ctypes.LeaseID{
		Owner:    ev.owner.String(),
		OSeq:     ev.oseq,
		Provider: ev.provider.String(),
	}
}

func (ev hostnameResourceEvent) GetHostname() string {
	return ev.hostname
}

func (ev hostnameResourceEvent) GetEventType() ctypes.ProviderResourceEvent {
	return ev.eventType
}

func (ev hostnameResourceEvent) GetServiceName() string {
	return ev.serviceName
}

func (ev hostnameResourceEvent) GetExternalPort() uint32 {
	return ev.externalPort
}

func (c *client) ObserveHostnameState(ctx context.Context) (<-chan ctypes.HostnameResourceEvent, error) {
	var lastResourceVersion string
	phpager := pager.New(func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		resources, err := c.ac.UbicnetV1().ProviderHosts(c.ns).List(ctx, opts)

		if err == nil && len(resources.GetResourceVersion()) != 0 {
			lastResourceVersion = resources.GetResourceVersion()
		}
		return resources, err
	})

	data := make([]crd.ProviderHost, 0, 128)
	err := phpager.EachListItem(ctx, metav1.ListOptions{}, func(obj runtime.Object) error {
		ph := obj.(*crd.ProviderHost)
		data = append(data, *ph)
		return nil
	})

	if err != nil {
		return nil, err
	}

	c.log.Info("starting hostname watch", "resourceVersion", lastResourceVersion)
	watcher, err := c.ac.UbicnetV1().ProviderHosts(c.ns).Watch(ctx, metav1.ListOptions{
		TypeMeta:             metav1.TypeMeta{},
		LabelSelector:        "",
		FieldSelector:        "",
		Watch:                false,
		AllowWatchBookmarks:  false,
		ResourceVersion:      lastResourceVersion,
		ResourceVersionMatch: "",
		TimeoutSeconds:       nil,
		Limit:                0,
		Continue:             "",
	})
	if err != nil {
		return nil, err
	}

	evData := make([]hostnameResourceEvent, len(data))
	for i, v := range data {
		//ownerAddr, err := sdktypes.AccAddressFromBech32(v.Spec.Owner)
		ownerAddr := common.HexToAddress(v.Spec.Owner)
		//providerAddr, err := sdktypes.AccAddressFromBech32(v.Spec.Provider)
		providerAddr := common.HexToAddress(v.Spec.Provider)
		ev := hostnameResourceEvent{
			eventType:    ctypes.ProviderResourceAdd,
			hostname:     v.Spec.Hostname,
			oseq:         v.Spec.Oseq,
			owner:        ownerAddr,
			provider:     providerAddr,
			serviceName:  v.Spec.ServiceName,
			externalPort: v.Spec.ExternalPort,
		}
		evData[i] = ev
	}

	data = nil

	output := make(chan ctypes.HostnameResourceEvent)

	go func() {
		defer close(output)
		for _, v := range evData {
			output <- v
		}
		evData = nil // do not hold the reference

		results := watcher.ResultChan()
		for {
			select {
			case result, ok := <-results:
				if !ok { // Channel closed when an error happens
					return
				}
				ph,ok := result.Object.(*crd.ProviderHost)
				if !ok{
				    continue
				}

				ownerAddr := common.HexToAddress(ph.Spec.Owner)
				providerAddr := common.HexToAddress(ph.Spec.Provider)
				ev := hostnameResourceEvent{
					hostname:     ph.Spec.Hostname,
					oseq:         ph.Spec.Oseq,
					owner:        ownerAddr,
					provider:     providerAddr,
					serviceName:  ph.Spec.ServiceName,
					externalPort: ph.Spec.ExternalPort,
				}
				switch result.Type {

				case watch.Added:
					ev.eventType = ctypes.ProviderResourceAdd
				case watch.Modified:
					ev.eventType = ctypes.ProviderResourceUpdate
				case watch.Deleted:
					ev.eventType = ctypes.ProviderResourceDelete

				case watch.Error:
					// Based on examination of the implementation code, this is basically never called anyways
					c.log.Error("watch error", "err", result.Object)

				default:

					continue
				}

				output <- ev

			case <-ctx.Done():
				return
			}
		}
	}()

	return output, nil
}

func (c *client) AllHostnames(ctx context.Context) ([]ctypes.ActiveHostname, error) {
	ingressPager := pager.New(func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		return c.ac.UbicnetV1().ProviderHosts(c.ns).List(ctx, opts)
	})

	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=true", builder.UbicManagedLabelName),
	}

	result := make([]ctypes.ActiveHostname, 0)

	err := ingressPager.EachListItem(ctx, listOptions, func(obj runtime.Object) error {
		ph := obj.(*crd.ProviderHost)
		hostname := ph.Spec.Hostname
		oseq := ph.Spec.Oseq

		owner, ok := ph.Labels[builder.UbicLeaseOwnerLabelName]
		if !ok || len(owner) == 0 {
			c.log.Error("providerhost missing owner label", "host", hostname)
			return nil
		}
		provider, ok := ph.Labels[builder.UbicLeaseProviderLabelName]
		if !ok || len(provider) == 0 {
			c.log.Error("providerhost missing provider label", "host", hostname)
			return nil
		}

		leaseID := ctypes.LeaseID{
			Owner:    owner,
			OSeq:     oseq,
			Provider: provider,
		}

		result = append(result, ctypes.ActiveHostname{
			ID:       leaseID,
			Hostname: hostname,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
