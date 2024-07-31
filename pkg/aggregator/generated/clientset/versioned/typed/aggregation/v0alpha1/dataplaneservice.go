// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by client-gen. DO NOT EDIT.

package v0alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v0alpha1 "github.com/grafana/grafana/pkg/aggregator/apis/aggregation/v0alpha1"
	aggregationv0alpha1 "github.com/grafana/grafana/pkg/aggregator/generated/applyconfiguration/aggregation/v0alpha1"
	scheme "github.com/grafana/grafana/pkg/aggregator/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// DataPlaneServicesGetter has a method to return a DataPlaneServiceInterface.
// A group's client should implement this interface.
type DataPlaneServicesGetter interface {
	DataPlaneServices() DataPlaneServiceInterface
}

// DataPlaneServiceInterface has methods to work with DataPlaneService resources.
type DataPlaneServiceInterface interface {
	Create(ctx context.Context, dataPlaneService *v0alpha1.DataPlaneService, opts v1.CreateOptions) (*v0alpha1.DataPlaneService, error)
	Update(ctx context.Context, dataPlaneService *v0alpha1.DataPlaneService, opts v1.UpdateOptions) (*v0alpha1.DataPlaneService, error)
	UpdateStatus(ctx context.Context, dataPlaneService *v0alpha1.DataPlaneService, opts v1.UpdateOptions) (*v0alpha1.DataPlaneService, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v0alpha1.DataPlaneService, error)
	List(ctx context.Context, opts v1.ListOptions) (*v0alpha1.DataPlaneServiceList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v0alpha1.DataPlaneService, err error)
	Apply(ctx context.Context, dataPlaneService *aggregationv0alpha1.DataPlaneServiceApplyConfiguration, opts v1.ApplyOptions) (result *v0alpha1.DataPlaneService, err error)
	ApplyStatus(ctx context.Context, dataPlaneService *aggregationv0alpha1.DataPlaneServiceApplyConfiguration, opts v1.ApplyOptions) (result *v0alpha1.DataPlaneService, err error)
	DataPlaneServiceExpansion
}

// dataPlaneServices implements DataPlaneServiceInterface
type dataPlaneServices struct {
	client rest.Interface
}

// newDataPlaneServices returns a DataPlaneServices
func newDataPlaneServices(c *AggregationV0alpha1Client) *dataPlaneServices {
	return &dataPlaneServices{
		client: c.RESTClient(),
	}
}

// Get takes name of the dataPlaneService, and returns the corresponding dataPlaneService object, and an error if there is any.
func (c *dataPlaneServices) Get(ctx context.Context, name string, options v1.GetOptions) (result *v0alpha1.DataPlaneService, err error) {
	result = &v0alpha1.DataPlaneService{}
	err = c.client.Get().
		Resource("dataplaneservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DataPlaneServices that match those selectors.
func (c *dataPlaneServices) List(ctx context.Context, opts v1.ListOptions) (result *v0alpha1.DataPlaneServiceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v0alpha1.DataPlaneServiceList{}
	err = c.client.Get().
		Resource("dataplaneservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested dataPlaneServices.
func (c *dataPlaneServices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("dataplaneservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a dataPlaneService and creates it.  Returns the server's representation of the dataPlaneService, and an error, if there is any.
func (c *dataPlaneServices) Create(ctx context.Context, dataPlaneService *v0alpha1.DataPlaneService, opts v1.CreateOptions) (result *v0alpha1.DataPlaneService, err error) {
	result = &v0alpha1.DataPlaneService{}
	err = c.client.Post().
		Resource("dataplaneservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(dataPlaneService).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a dataPlaneService and updates it. Returns the server's representation of the dataPlaneService, and an error, if there is any.
func (c *dataPlaneServices) Update(ctx context.Context, dataPlaneService *v0alpha1.DataPlaneService, opts v1.UpdateOptions) (result *v0alpha1.DataPlaneService, err error) {
	result = &v0alpha1.DataPlaneService{}
	err = c.client.Put().
		Resource("dataplaneservices").
		Name(dataPlaneService.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(dataPlaneService).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *dataPlaneServices) UpdateStatus(ctx context.Context, dataPlaneService *v0alpha1.DataPlaneService, opts v1.UpdateOptions) (result *v0alpha1.DataPlaneService, err error) {
	result = &v0alpha1.DataPlaneService{}
	err = c.client.Put().
		Resource("dataplaneservices").
		Name(dataPlaneService.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(dataPlaneService).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the dataPlaneService and deletes it. Returns an error if one occurs.
func (c *dataPlaneServices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("dataplaneservices").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *dataPlaneServices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("dataplaneservices").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched dataPlaneService.
func (c *dataPlaneServices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v0alpha1.DataPlaneService, err error) {
	result = &v0alpha1.DataPlaneService{}
	err = c.client.Patch(pt).
		Resource("dataplaneservices").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied dataPlaneService.
func (c *dataPlaneServices) Apply(ctx context.Context, dataPlaneService *aggregationv0alpha1.DataPlaneServiceApplyConfiguration, opts v1.ApplyOptions) (result *v0alpha1.DataPlaneService, err error) {
	if dataPlaneService == nil {
		return nil, fmt.Errorf("dataPlaneService provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(dataPlaneService)
	if err != nil {
		return nil, err
	}
	name := dataPlaneService.Name
	if name == nil {
		return nil, fmt.Errorf("dataPlaneService.Name must be provided to Apply")
	}
	result = &v0alpha1.DataPlaneService{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("dataplaneservices").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *dataPlaneServices) ApplyStatus(ctx context.Context, dataPlaneService *aggregationv0alpha1.DataPlaneServiceApplyConfiguration, opts v1.ApplyOptions) (result *v0alpha1.DataPlaneService, err error) {
	if dataPlaneService == nil {
		return nil, fmt.Errorf("dataPlaneService provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(dataPlaneService)
	if err != nil {
		return nil, err
	}

	name := dataPlaneService.Name
	if name == nil {
		return nil, fmt.Errorf("dataPlaneService.Name must be provided to Apply")
	}

	result = &v0alpha1.DataPlaneService{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("dataplaneservices").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
