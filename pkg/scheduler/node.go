package scheduler

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

const (
	// GPUResourceName is the resource name for NVIDIA GPUs
	GPUResourceName = "nvidia.com/gpu"
	// GPUMemoryResourceName is the resource name for GPU memory (in MiB)
	GPUMemoryResourceName = "hami.io/gpu-memory"
	// GPUCoreResourceName is the resource name for GPU core utilization (percentage)
	GPUCoreResourceName = "hami.io/gpu-core"
)

// GPUDevice represents a single GPU device on a node.
type GPUDevice struct {
	// Index is the GPU device index (e.g., 0, 1, 2)
	Index int
	// UUID is the unique identifier for the GPU
	UUID string
	// MemoryTotal is the total GPU memory in MiB
	MemoryTotal int64
	// MemoryUsed is the currently allocated GPU memory in MiB
	MemoryUsed int64
	// CoreUsed is the currently allocated GPU core utilization percentage
	CoreUsed int64
}

// MemoryFree returns the available GPU memory in MiB.
func (d *GPUDevice) MemoryFree() int64 {
	return d.MemoryTotal - d.MemoryUsed
}

// NodeInfo holds GPU resource information for a Kubernetes node.
type NodeInfo struct {
	// Name is the Kubernetes node name
	Name string
	// GPUDevices is the list of GPU devices on this node
	GPUDevices []*GPUDevice
}

// NewNodeInfo creates a NodeInfo from a Kubernetes Node object and its
// associated GPU annotations set by the HAMi device plugin.
func NewNodeInfo(node *v1.Node) (*NodeInfo, error) {
	if node == nil {
		return nil, fmt.Errorf("node must not be nil")
	}
	ni := &NodeInfo{
		Name:       node.Name,
		GPUDevices: make([]*GPUDevice, 0),
	}
	return ni, nil
}

// TotalMemory returns the sum of all GPU memory across devices on the node.
func (n *NodeInfo) TotalMemory() int64 {
	var total int64
	for _, d := range n.GPUDevices {
		total += d.MemoryTotal
	}
	return total
}

// FreeMemory returns the sum of free GPU memory across all devices on the node.
func (n *NodeInfo) FreeMemory() int64 {
	var free int64
	for _, d := range n.GPUDevices {
		free += d.MemoryFree()
	}
	return free
}

// GPUCount returns the number of GPU devices on the node.
func (n *NodeInfo) GPUCount() int {
	return len(n.GPUDevices)
}

// CanFit reports whether the node has enough free GPU memory and devices to
// satisfy the requested memory (memMiB) across the requested number of GPUs.
func (n *NodeInfo) CanFit(gpuCount int, memMiB int64) bool {
	if gpuCount <= 0 || gpuCount > n.GPUCount() {
		return false
	}
	// Count devices that individually satisfy the per-GPU memory request.
	eligible := 0
	for _, d := range n.GPUDevices {
		if d.MemoryFree() >= memMiB {
			eligible++
		}
	}
	return eligible >= gpuCount
}
