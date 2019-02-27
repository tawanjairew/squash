// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sync"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

var (
	mApiSnapshotIn  = stats.Int64("api.squash.solo.io/snap_emitter/snap_in", "The number of snapshots in", "1")
	mApiSnapshotOut = stats.Int64("api.squash.solo.io/snap_emitter/snap_out", "The number of snapshots out", "1")

	apisnapshotInView = &view.View{
		Name:        "api.squash.solo.io_snap_emitter/snap_in",
		Measure:     mApiSnapshotIn,
		Description: "The number of snapshots updates coming in",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
	apisnapshotOutView = &view.View{
		Name:        "api.squash.solo.io/snap_emitter/snap_out",
		Measure:     mApiSnapshotOut,
		Description: "The number of snapshots updates going out",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(apisnapshotInView, apisnapshotOutView)
}

type ApiEmitter interface {
	Register() error
	DebugAttachment() DebugAttachmentClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ApiSnapshot, <-chan error, error)
}

func NewApiEmitter(debugAttachmentClient DebugAttachmentClient) ApiEmitter {
	return NewApiEmitterWithEmit(debugAttachmentClient, make(chan struct{}))
}

func NewApiEmitterWithEmit(debugAttachmentClient DebugAttachmentClient, emit <-chan struct{}) ApiEmitter {
	return &apiEmitter{
		debugAttachment: debugAttachmentClient,
		forceEmit:       emit,
	}
}

type apiEmitter struct {
	forceEmit       <-chan struct{}
	debugAttachment DebugAttachmentClient
}

func (c *apiEmitter) Register() error {
	if err := c.debugAttachment.Register(); err != nil {
		return err
	}
	return nil
}

func (c *apiEmitter) DebugAttachment() DebugAttachmentClient {
	return c.debugAttachment
}

func (c *apiEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ApiSnapshot, <-chan error, error) {

	if len(watchNamespaces) == 0 {
		watchNamespaces = []string{""}
	}

	for _, ns := range watchNamespaces {
		if ns == "" && len(watchNamespaces) > 1 {
			return nil, nil, errors.Errorf("the \"\" namespace is used to watch all namespaces. Snapshots can either be tracked for " +
				"specific namespaces or \"\" AllNamespaces, but not both.")
		}
	}

	errs := make(chan error)
	var done sync.WaitGroup
	ctx := opts.Ctx
	/* Create channel for DebugAttachment */
	type debugAttachmentListWithNamespace struct {
		list      DebugAttachmentList
		namespace string
	}
	debugAttachmentChan := make(chan debugAttachmentListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup namespaced watch for DebugAttachment */
		debugAttachmentNamespacesChan, debugAttachmentErrs, err := c.debugAttachment.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting DebugAttachment watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, debugAttachmentErrs, namespace+"-debugattachments")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-ctx.Done():
					return
				case debugAttachmentList := <-debugAttachmentNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case debugAttachmentChan <- debugAttachmentListWithNamespace{list: debugAttachmentList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *ApiSnapshot)
	go func() {
		originalSnapshot := ApiSnapshot{}
		currentSnapshot := originalSnapshot.Clone()
		timer := time.NewTicker(time.Second * 1)
		sync := func() {
			if originalSnapshot.Hash() == currentSnapshot.Hash() {
				return
			}

			stats.Record(ctx, mApiSnapshotOut.M(1))
			originalSnapshot = currentSnapshot.Clone()
			sentSnapshot := currentSnapshot.Clone()
			snapshots <- &sentSnapshot
		}

		for {
			record := func() { stats.Record(ctx, mApiSnapshotIn.M(1)) }

			select {
			case <-timer.C:
				sync()
			case <-ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case debugAttachmentNamespacedList := <-debugAttachmentChan:
				record()

				namespace := debugAttachmentNamespacedList.namespace
				debugAttachmentList := debugAttachmentNamespacedList.list

				currentSnapshot.Debugattachments[namespace] = debugAttachmentList
			}
		}
	}()
	return snapshots, errs, nil
}
