package ring

import (
	"context"
	"sort"
	"time"
)

// ReplicationSet describes the instances to talk to for a given key, and how
// many errors to tolerate.
type ReplicationSet struct {
	Instances []InstanceDesc

	// Maximum number of tolerated failing instances. Max errors and max unavailable zones are
	// mutually exclusive.
	MaxErrors int

	// Maximum number of different zones in which instances can fail. Max unavailable zones and
	// max errors are mutually exclusive.
	MaxUnavailableZones int
}

// Do function f in parallel for all replicas in the set, erroring if we exceed
// MaxErrors and returning early otherwise.
// Return a slice of all results from f, or nil if an error occurred.
func (r ReplicationSet) Do(ctx context.Context, delay time.Duration, f func(context.Context, *InstanceDesc) (interface{}, error)) ([]interface{}, error) {
	type instanceResult struct {
		res      interface{}
		err      error
		instance *InstanceDesc
	}

	// Initialise the result tracker, which is use to keep track of successes and failures.
	var tracker replicationSetResultTracker
	if r.MaxUnavailableZones > 0 {
		tracker = newZoneAwareResultTracker(ctx, r.Instances, r.MaxUnavailableZones)
	} else {
		tracker = newDefaultResultTracker(ctx, r.Instances, r.MaxErrors)
	}

	var (
		ch         = make(chan instanceResult, len(r.Instances))
		forceStart = make(chan struct{}, r.MaxErrors)
	)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Spawn a goroutine for each instance.
	for i := range r.Instances {
		go func(i int, ing *InstanceDesc) {
			// Wait to send extra requests. Works only when zone-awareness is disabled.
			if delay > 0 && r.MaxUnavailableZones == 0 && i >= len(r.Instances)-r.MaxErrors {
				after := time.NewTimer(delay)
				defer after.Stop()
				select {
				case <-ctx.Done():
					return
				case <-forceStart:
				case <-after.C:
				}
			}
			result, err := f(ctx, ing)
			ch <- instanceResult{
				res:      result,
				err:      err,
				instance: ing,
			}
		}(i, &r.Instances[i])
	}

	results := make([]interface{}, 0, len(r.Instances))

	for !tracker.succeeded() {
		select {
		case res := <-ch:
			tracker.done(res.instance, res.err)
			if res.err != nil {
				if tracker.failed() {
					return nil, res.err
				}

				// force one of the delayed requests to start
				if delay > 0 && r.MaxUnavailableZones == 0 {
					forceStart <- struct{}{}
				}
			} else {
				results = append(results, res.res)
			}

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return results, nil
}

// DoUntilQuorum runs function f in parallel for all replicas in r.
//
// If r.MaxUnavailableZones is greater than zero:
//   - DoUntilQuorum returns an error if calls to f for instances in r.MaxUnavailableZones zones return errors
//   - Otherwise, DoUntilQuorum returns all results from all replicas in the first zones for which f succeeds
//     for every instance in that zone (eg. if there are 3 zones and r.MaxUnavailableZones is 1, DoUntilQuorum will
//     return the results from all instances in 2 zones, even if all calls to f succeed).
//
// Otherwise:
//   - DoUntilQuorum returns an error if r.MaxErrors calls to f return errors
//   - Otherwise, DoUntilQuorum returns all results from the first len(r.Instances) - r.MaxErrors instances
//     (eg. if there are 6 replicas and r.MaxErrors is 2, DoUntilQuorum will return the results from the first 4
//     successful calls to f, even if all 6 calls to f succeed).
//
// Any results from successful calls to f that are not returned by DoUntilQuorum will be passed to cleanupFunc,
// including when DoUntilQuorum returns an error or only returns a subset of successful results. cleanupFunc may
// be called both before and after DoUntilQuorum returns.
//
// DoUntilQuorum cancels the context.Context passed to each invocation of f if the result of that invocation of
// f will not be returned. If the result of that invocation of f will be returned, the context.Context passed
// to that invocation of f will not be cancelled by DoUntilQuorum, but the context.Context is a child of ctx
// and so will be cancelled if ctx is cancelled.
func DoUntilQuorum[T any](ctx context.Context, r ReplicationSet, f func(context.Context, *InstanceDesc) (T, error), cleanupFunc func(T)) ([]T, error) {
	resultsChan := make(chan instanceResult[T], len(r.Instances))
	resultsRemaining := len(r.Instances)

	defer func() {
		go func() {
			for resultsRemaining > 0 {
				result := <-resultsChan
				resultsRemaining--

				if result.err == nil {
					cleanupFunc(result.result)
				}
			}
		}()
	}()

	var tracker replicationSetResultTracker
	if r.MaxUnavailableZones > 0 {
		tracker = newZoneAwareResultTracker(ctx, r.Instances, r.MaxUnavailableZones)
	} else {
		tracker = newDefaultResultTracker(ctx, r.Instances, r.MaxErrors)
	}

	for i := range r.Instances {
		instance := &r.Instances[i]
		instanceCtx := tracker.contextFor(instance)

		go func(desc *InstanceDesc) {
			result, err := f(instanceCtx, desc)
			resultsChan <- instanceResult[T]{
				result:   result,
				err:      err,
				instance: desc,
			}
		}(instance)
	}

	resultsMap := make(map[*InstanceDesc]T, len(r.Instances))
	cleanupResultsAlreadyReceived := func() {
		for _, result := range resultsMap {
			cleanupFunc(result)
		}
	}

	for !tracker.succeeded() {
		select {
		case <-ctx.Done():
			// No need to cancel individual instance contexts, as they inherit the cancellation from ctx.
			cleanupResultsAlreadyReceived()

			return nil, ctx.Err()
		case result := <-resultsChan:
			resultsRemaining--
			tracker.done(result.instance, result.err)

			if result.err == nil {
				resultsMap[result.instance] = result.result
			} else if tracker.failed() {
				tracker.cancelAllContexts()
				cleanupResultsAlreadyReceived()
				return nil, result.err
			}
		}
	}

	results := make([]T, 0, len(r.Instances))

	for i := range r.Instances {
		instance := &r.Instances[i]
		result, haveResult := resultsMap[instance]

		if haveResult {
			if tracker.shouldIncludeResultFrom(instance) {
				results = append(results, result)
			} else {
				tracker.cancelContextFor(instance)
				cleanupFunc(result)
			}
		} else {
			// Nothing to clean up (yet) - this will be handled by deferred call above.
			tracker.cancelContextFor(instance)
		}
	}

	return results, nil
}

type instanceResult[T any] struct {
	result   T
	err      error
	instance *InstanceDesc
}

// Includes returns whether the replication set includes the replica with the provided addr.
func (r ReplicationSet) Includes(addr string) bool {
	for _, instance := range r.Instances {
		if instance.GetAddr() == addr {
			return true
		}
	}

	return false
}

// GetAddresses returns the addresses of all instances within the replication set. Returned slice
// order is not guaranteed.
func (r ReplicationSet) GetAddresses() []string {
	addrs := make([]string, 0, len(r.Instances))
	for _, desc := range r.Instances {
		addrs = append(addrs, desc.Addr)
	}
	return addrs
}

// GetAddressesWithout returns the addresses of all instances within the replication set while
// excluding the specified address. Returned slice order is not guaranteed.
func (r ReplicationSet) GetAddressesWithout(exclude string) []string {
	addrs := make([]string, 0, len(r.Instances))
	for _, desc := range r.Instances {
		if desc.Addr != exclude {
			addrs = append(addrs, desc.Addr)
		}
	}
	return addrs
}

// ZoneCount returns the number of unique zones represented by instances within the replication set.
func (r ReplicationSet) ZoneCount() int {
	zones := map[string]struct{}{}

	for _, i := range r.Instances {
		zones[i.Zone] = struct{}{}
	}

	return len(zones)
}

// HasReplicationSetChanged returns true if two replications sets are the same (with possibly different timestamps),
// false if they differ in any way (number of instances, instance states, tokens, zones, ...).
func HasReplicationSetChanged(before, after ReplicationSet) bool {
	return hasReplicationSetChangedExcluding(before, after, func(i *InstanceDesc) {
		i.Timestamp = 0
	})
}

// HasReplicationSetChangedWithoutState returns true if two replications sets
// are the same (with possibly different timestamps and instance states),
// false if they differ in any other way (number of instances, tokens, zones, ...).
func HasReplicationSetChangedWithoutState(before, after ReplicationSet) bool {
	return hasReplicationSetChangedExcluding(before, after, func(i *InstanceDesc) {
		i.Timestamp = 0
		i.State = PENDING
	})
}

// Do comparison of replicasets, but apply a function first
// to be able to exclude (reset) some values
func hasReplicationSetChangedExcluding(before, after ReplicationSet, exclude func(*InstanceDesc)) bool {
	beforeInstances := before.Instances
	afterInstances := after.Instances

	if len(beforeInstances) != len(afterInstances) {
		return true
	}

	sort.Sort(ByAddr(beforeInstances))
	sort.Sort(ByAddr(afterInstances))

	for i := 0; i < len(beforeInstances); i++ {
		b := beforeInstances[i]
		a := afterInstances[i]

		exclude(&a)
		exclude(&b)

		if !b.Equal(a) {
			return true
		}
	}

	return false
}
