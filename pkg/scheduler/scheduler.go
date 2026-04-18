// Package scheduler provides GPU resource scheduling capabilities for HAMi.
// It implements a Kubernetes scheduler extender that handles GPU device allocation.
package scheduler

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// Config holds the configuration for the HAMi scheduler.
type Config struct {
	// SchedulerPort is the port the scheduler extender listens on.
	SchedulerPort int
	// KubeConfigPath is the path to the kubeconfig file.
	KubeConfigPath string
	// NodeSelectors are labels used to identify GPU nodes.
	NodeSelectors map[string]string
}

// Scheduler is the main HAMi GPU scheduler extender.
type Scheduler struct {
	config    *Config
	client    kubernetes.Interface
	mu        sync.RWMutex
	stopCh    chan struct{}
}

// New creates a new Scheduler instance with the provided configuration.
func New(cfg *Config, client kubernetes.Interface) (*Scheduler, error) {
	if cfg == nil {
		return nil, fmt.Errorf("scheduler config must not be nil")
	}
	if client == nil {
		return nil, fmt.Errorf("kubernetes client must not be nil")
	}
	if cfg.SchedulerPort <= 0 {
		cfg.SchedulerPort = 443
	}

	return &Scheduler{
		config: cfg,
		client: client,
		stopCh: make(chan struct{}),
	}, nil
}

// Start begins the scheduler extender, registering handlers and starting
// any background reconciliation loops.
func (s *Scheduler) Start(ctx context.Context) error {
	klog.InfoS("Starting HAMi scheduler extender", "port", s.config.SchedulerPort)

	errCh := make(chan error, 1)

	go func() {
		if err := s.run(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("scheduler run error: %w", err)
	case <-ctx.Done():
		klog.InfoS("Scheduler context cancelled, shutting down")
		close(s.stopCh)
		return nil
	}
}

// run is the internal loop for the scheduler extender.
func (s *Scheduler) run(ctx context.Context) error {
	// TODO: Register HTTP handlers for filter, prioritize, and bind endpoints.
	// TODO: Start node GPU resource informer.
	klog.V(4).InfoS("Scheduler run loop started")
	<-ctx.Done()
	return nil
}

// Stop gracefully shuts down the scheduler.
func (s *Scheduler) Stop() {
	klog.InfoS("Stopping HAMi scheduler")
	select {
	case <-s.stopCh:
		// already closed
	default:
		close(s.stopCh)
	}
}
