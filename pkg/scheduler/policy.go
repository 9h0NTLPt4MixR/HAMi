package scheduler

import (
	"fmt"
	"sort"

	v1 "k8s.io/api/core/v1"
)

// PolicyType defines the scheduling policy for GPU allocation.
type PolicyType string

const (
	// PolicyBinpack prefers nodes with the most allocated resources to minimize fragmentation.
	PolicyBinpack PolicyType = "binpack"
	// PolicySpread prefers nodes with the least allocated resources to spread workloads.
	PolicySpread PolicyType = "spread"
)

// NodeScore represents a node and its computed scheduling score.
type NodeScore struct {
	Node  *v1.Node
	Score int64
}

// NodeScoreList is a sortable list of NodeScore.
type NodeScoreList []NodeScore

func (n NodeScoreList) Len() int           { return len(n) }
func (n NodeScoreList) Less(i, j int) bool { return n[i].Score > n[j].Score }
func (n NodeScoreList) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

// Policy defines the interface for GPU scheduling policies.
type Policy interface {
	// Score computes a score for the given node based on current resource usage.
	Score(node *v1.Node, available, total int64) int64
	// Name returns the name of the policy.
	Name() PolicyType
}

// BinpackPolicy scores nodes higher when they have less available capacity,
// encouraging tight packing of workloads.
type BinpackPolicy struct{}

func (b *BinpackPolicy) Name() PolicyType { return PolicyBinpack }

func (b *BinpackPolicy) Score(node *v1.Node, available, total int64) int64 {
	if total == 0 {
		return 0
	}
	// Higher score for nodes with less available (more utilized)
	used := total - available
	return (used * 100) / total
}

// SpreadPolicy scores nodes higher when they have more available capacity,
// encouraging even distribution of workloads.
type SpreadPolicy struct{}

func (s *SpreadPolicy) Name() PolicyType { return PolicySpread }

func (s *SpreadPolicy) Score(node *v1.Node, available, total int64) int64 {
	if total == 0 {
		return 0
	}
	// Higher score for nodes with more available (less utilized)
	return (available * 100) / total
}

// NewPolicy creates a Policy instance for the given PolicyType.
func NewPolicy(p PolicyType) (Policy, error) {
	switch p {
	case PolicyBinpack:
		return &BinpackPolicy{}, nil
	case PolicySpread:
		return &SpreadPolicy{}, nil
	default:
		return nil, fmt.Errorf("unknown scheduling policy: %s", p)
	}
}

// RankNodes sorts the provided NodeScoreList according to the active policy
// and returns an ordered slice of nodes, highest score first.
func RankNodes(scores NodeScoreList) []*v1.Node {
	sort.Sort(scores)
	nodes := make([]*v1.Node, 0, len(scores))
	for _, ns := range scores {
		nodes = append(nodes, ns.Node)
	}
	return nodes
}
