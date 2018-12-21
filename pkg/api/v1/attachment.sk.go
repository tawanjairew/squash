// Code generated by protoc-gen-solo-kit. DO NOT EDIT.

package v1

import (
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func NewAttachment(namespace, name string) *Attachment {
	return &Attachment{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Attachment) SetStatus(status core.Status) {
	r.Status = status
}

func (r *Attachment) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

type AttachmentList []*Attachment
type AttachmentsByNamespace map[string]AttachmentList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list AttachmentList) Find(namespace, name string) (*Attachment, error) {
	for _, attachment := range list {
		if attachment.Metadata.Name == name {
			if namespace == "" || attachment.Metadata.Namespace == namespace {
				return attachment, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find attachment %v.%v", namespace, name)
}

func (list AttachmentList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, attachment := range list {
		ress = append(ress, attachment)
	}
	return ress
}

func (list AttachmentList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, attachment := range list {
		ress = append(ress, attachment)
	}
	return ress
}

func (list AttachmentList) Names() []string {
	var names []string
	for _, attachment := range list {
		names = append(names, attachment.Metadata.Name)
	}
	return names
}

func (list AttachmentList) NamespacesDotNames() []string {
	var names []string
	for _, attachment := range list {
		names = append(names, attachment.Metadata.Namespace+"."+attachment.Metadata.Name)
	}
	return names
}

func (list AttachmentList) Sort() AttachmentList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
	return list
}

func (list AttachmentList) Clone() AttachmentList {
	var attachmentList AttachmentList
	for _, attachment := range list {
		attachmentList = append(attachmentList, proto.Clone(attachment).(*Attachment))
	}
	return attachmentList
}

func (list AttachmentList) ByNamespace() AttachmentsByNamespace {
	byNamespace := make(AttachmentsByNamespace)
	for _, attachment := range list {
		byNamespace.Add(attachment)
	}
	return byNamespace
}

func (byNamespace AttachmentsByNamespace) Add(attachment ...*Attachment) {
	for _, item := range attachment {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace AttachmentsByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace AttachmentsByNamespace) List() AttachmentList {
	var list AttachmentList
	for _, attachmentList := range byNamespace {
		list = append(list, attachmentList...)
	}
	return list.Sort()
}

func (byNamespace AttachmentsByNamespace) Clone() AttachmentsByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &Attachment{}

// Kubernetes Adapter for Attachment

func (o *Attachment) GetObjectKind() schema.ObjectKind {
	t := AttachmentCrd.TypeMeta()
	return &t
}

func (o *Attachment) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Attachment)
}

var AttachmentCrd = crd.NewCrd("squash.solo.io",
	"attachments",
	"squash.solo.io",
	"v1",
	"Attachment",
	"att",
	&Attachment{})
