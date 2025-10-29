package api

import (
	"fmt"
	"sort"
	"strings"
)

type ConstraintGroupType string

const (
	GroupExactlyOneOf  ConstraintGroupType = "exactly_one_of"
	GroupAtLeastOneOf  ConstraintGroupType = "at_least_one_of"
	GroupRequiredWith  ConstraintGroupType = "required_with"
	GroupConflictsWith ConstraintGroupType = "conflicts_with"
)

type ConstraintGroup struct {
	Type    ConstraintGroupType
	Key     string
	Members *[]string
}

type ConstraintGroupRegistry struct {
	groups map[string]*ConstraintGroup
}

func NewConstraintGroupRegistry() *ConstraintGroupRegistry {
	return &ConstraintGroupRegistry{
		groups: make(map[string]*ConstraintGroup),
	}
}

func (r *ConstraintGroupRegistry) GetOrCreateGroup(groupType ConstraintGroupType, members []string) *ConstraintGroup {
	key := r.generateGroupKey(groupType, members)
	if group, exists := r.groups[key]; exists {
		return group
	}

	membersCopy := make([]string, len(members))
	copy(membersCopy, members)
	group := &ConstraintGroup{
		Type:    groupType,
		Key:     key,
		Members: &membersCopy,
	}

	r.groups[key] = group
	return group
}

func (r *ConstraintGroupRegistry) AddMembersToGroup(groupType ConstraintGroupType, originalMembers []string, newMembers []string) {
	group := r.GetOrCreateGroup(groupType, originalMembers)
	updatedMembers := append(*group.Members, newMembers...)
	*group.Members = deduplicateSliceOfStrings(updatedMembers)
}

func (r *ConstraintGroupRegistry) FindGroupByKey(key string) *ConstraintGroup {
	return r.groups[key]
}

func (r *ConstraintGroupRegistry) FindGroupByMembers(groupType ConstraintGroupType, members []string) *ConstraintGroup {
	key := r.generateGroupKey(groupType, members)
	return r.FindGroupByKey(key)
}

func (r *ConstraintGroupRegistry) generateGroupKey(groupType ConstraintGroupType, members []string) string {
	sortedMembers := make([]string, len(members))
	copy(sortedMembers, members)
	sort.Strings(sortedMembers)
	return fmt.Sprintf("%s:%s", groupType, strings.Join(sortedMembers, ","))
}

func (r *ConstraintGroupRegistry) GetAllGroups() map[string]*ConstraintGroup {
	return r.groups
}
