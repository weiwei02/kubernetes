package cacher

import (
	"context"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/cacher/metrics"
	"k8s.io/klog/v2"
	"k8s.io/utils/clock"
)

const (
	currentResourceVersionInterval = "APISERVER_GET_CURRENT_RESOURCE_VERSION_INTERVAL"
)

type cceStorageVersionManager struct {
	// for testing timeouts.
	clock          clock.Clock
	resourcePrefix string
	objectType     string

	// getCurrentStorageRV is a function that returns the current storage version from storage
	getCurrentStorageRV  func(ctx context.Context) (uint64, error)
	requestWatchProgress func() error

	// lock effectively protects access to the underlying source
	lock sync.Locker
	// refreshStorageInterval is the interval at which we will try to update storage version
	// default is 0s
	refreshStorageInterval time.Duration
	cacherResourceVersion  uint64
	storageResourceVersion uint64
	// lastStorageResourceVersionUpdateTime is the last time we updated storage version
	// default is zero time.Time
	lastStorageResourceVersionUpdateTime time.Time
}

func newCCEStorageVersionManager(s storage.Interface, newListFunc func() runtime.Object, resourcePrefix, objectType string, contextMetadata metadata.MD) *cceStorageVersionManager {
	var refreshStorageInterval, _ = time.ParseDuration(os.Getenv(currentResourceVersionInterval))
	vm := &cceStorageVersionManager{
		lock:                                 &sync.Mutex{},
		resourcePrefix:                       resourcePrefix,
		objectType:                           objectType,
		refreshStorageInterval:               refreshStorageInterval,
		storageResourceVersion:               0,
		lastStorageResourceVersionUpdateTime: time.Time{},
		clock:                                clock.RealClock{},
	}
	vm.getCurrentStorageRV = func(ctx context.Context) (uint64, error) {
		return storage.GetCurrentResourceVersionFromStorage(ctx, s, newListFunc, resourcePrefix, objectType)
	}

	// Add grpc context metadata to watch and progress notify requests done by cacher to:
	// * Prevent starvation of watch opened by cacher, by moving it to separate Watch RPC than watch request that bypass cacher.
	// * Ensure that progress notification requests are executed on the same Watch RPC as their watch, which is required for it to work.
	vm.requestWatchProgress = func() error {
		ctx := context.Background()
		if contextMetadata != nil {
			ctx = metadata.NewOutgoingContext(context.Background(), contextMetadata)
		}
		return s.RequestWatchProgress(ctx)
	}
	return vm
}

func (vm *cceStorageVersionManager) run() {
	if vm.refreshStorageInterval > 0 {
		klog.V(2).Infof("starting cce periodic storage version refresher for %s", vm.resourcePrefix)
		timer := time.NewTicker(vm.refreshStorageInterval)
		go func() {
			for {
				select {
				case <-timer.C:
					rv, err := vm.getCurrentStorageRV(context.TODO())
					if err != nil {
						klog.V(2).Infof("failed to get current resource version for %s from storage %v", vm.resourcePrefix, err)
						continue
					}
					vm.updateStorageResourceVersion(rv)

					err = vm.requestWatchProgress()
					if err != nil {
						klog.V(4).InfoS("Error requesting bookmark", "err", err)
					}
				}
			}
		}()
	}
}

func (vm *cceStorageVersionManager) updateCacherResourceVersion(resourceVersion uint64) {
	metrics.CCEWatchCacheResourceVersion.WithLabelValues(vm.resourcePrefix, metrics.SourceCacher).Set(float64(resourceVersion))

	vm.lock.Lock()
	defer vm.lock.Unlock()
	vm.cacherResourceVersion = resourceVersion
	if resourceVersion >= vm.storageResourceVersion && !vm.lastStorageResourceVersionUpdateTime.IsZero() {
		metrics.CCEWatchCacheReadWait.WithLabelValues(vm.resourcePrefix).Observe(float64(vm.clock.Since(vm.lastStorageResourceVersionUpdateTime).Seconds()))
		vm.lastStorageResourceVersionUpdateTime = time.Time{}
	}
}

func (vm *cceStorageVersionManager) updateStorageResourceVersion(resourceVersion uint64) {
	metrics.CCEWatchCacheResourceVersion.WithLabelValues(vm.resourcePrefix, metrics.SourceStorage).Set(float64(resourceVersion))

	vm.lock.Lock()
	defer vm.lock.Unlock()
	// only update storage version if last update storage version is zero
	// This is to prevent the value of lastStorageResourceVersionUpdateTime time
	// from being refreshed when the delay between apiserver and etcd is greater
	// than refreshStorageInterval, so that the issue of excessively long timeouts
	// can be calculated
	if vm.lastStorageResourceVersionUpdateTime.IsZero() {
		if vm.cacherResourceVersion >= resourceVersion {
			metrics.CCEWatchCacheReadWait.WithLabelValues(vm.resourcePrefix).Observe(0)
			return
		}
		vm.storageResourceVersion = resourceVersion
		vm.lastStorageResourceVersionUpdateTime = vm.clock.Now()
	}
}
