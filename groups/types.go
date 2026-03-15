package groups

import (
	"regexp"

	"github.com/google/uuid"
)

// GroupSource indicates how a group was created.
type GroupSource string

const (
	// GroupSourceUser indicates the group was created by a user.
	GroupSourceUser GroupSource = "user"
	// GroupSourceOIDC indicates the group was provisioned by an OIDC provider.
	GroupSourceOIDC GroupSource = "oidc"
)

// ReducedUser is a minimal user representation included in group membership.
type ReducedUser struct {
	ID        uuid.UUID `json:"id" format:"uuid"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
}

// Group represents a group of users within an organization.
type Group struct {
	ID                      uuid.UUID   `json:"id" format:"uuid"`
	Name                    string      `json:"name"`
	DisplayName             string      `json:"display_name"`
	OrganizationID          uuid.UUID   `json:"organization_id" format:"uuid"`
	Members                 []ReducedUser `json:"members"`
	// TotalMemberCount shows the total count of members, even if the caller
	// is not authorized to read group member details. May be greater than
	// len(Members).
	TotalMemberCount        int         `json:"total_member_count"`
	AvatarURL               string      `json:"avatar_url"`
	QuotaAllowance          int         `json:"quota_allowance"`
	Source                  GroupSource `json:"source"`
	OrganizationName        string      `json:"organization_name"`
	OrganizationDisplayName string      `json:"organization_display_name"`
}

// IsEveryone returns true if the group is the implicit "everyone" group for
// its organization.
func (g Group) IsEveryone() bool {
	return g.ID == g.OrganizationID
}

// CreateGroupRequest contains the parameters for creating a new group.
type CreateGroupRequest struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name,omitempty"`
	AvatarURL      string `json:"avatar_url,omitempty"`
	QuotaAllowance int    `json:"quota_allowance,omitempty"`
}

// PatchGroupRequest contains the parameters for updating an existing group.
type PatchGroupRequest struct {
	AddUsers       []string `json:"add_users,omitempty"`
	RemoveUsers    []string `json:"remove_users,omitempty"`
	Name           string   `json:"name,omitempty"`
	DisplayName    *string  `json:"display_name,omitempty"`
	AvatarURL      *string  `json:"avatar_url,omitempty"`
	QuotaAllowance *int     `json:"quota_allowance,omitempty"`
}

// GroupArguments contains filter parameters for listing groups.
type GroupArguments struct {
	// Organization can be an org UUID or name.
	Organization string
	// HasMember can be a user UUID or username.
	HasMember string
	// GroupIDs is a list of group UUIDs to filter by.
	// If empty, all groups will be returned.
	GroupIDs []uuid.UUID
}

// GroupSyncSettings configures how OIDC groups are synced into Lattice groups.
type GroupSyncSettings struct {
	// Field selects the claim field to be used as the created user's groups.
	// If empty, no group updates will come from the OIDC provider.
	Field string `json:"field"`
	// Mapping maps from an OIDC group to Lattice group IDs.
	Mapping map[string][]uuid.UUID `json:"mapping"`
	// RegexFilter is a regular expression that filters the groups returned by
	// the OIDC provider. Any group not matched will be ignored.
	RegexFilter *regexp.Regexp `json:"regex_filter"`
	// AutoCreateMissing controls whether groups returned by the OIDC provider
	// are automatically created in Lattice if they are missing.
	AutoCreateMissing bool `json:"auto_create_missing_groups"`
	// LegacyNameMapping is deprecated. It remaps an IDP group name to
	// a Lattice group name. Use Mapping instead.
	//
	// Deprecated: Use Mapping instead.
	LegacyNameMapping map[string]string `json:"legacy_group_name_mapping,omitempty"`
}

// RoleSyncSettings configures how OIDC groups are synced into organization roles.
type RoleSyncSettings struct {
	// Field selects the claim field to be used as the created user's roles.
	// If empty, no role updates will come from the OIDC provider.
	Field string `json:"field"`
	// Mapping maps from an OIDC group to Lattice organization roles.
	Mapping map[string][]string `json:"mapping"`
}

// OrganizationSyncSettings configures how OIDC claims are synced to
// organization membership.
type OrganizationSyncSettings struct {
	// Field selects the claim field to be used as the created user's
	// organizations. If empty, no organization updates will come from
	// the OIDC provider.
	Field string `json:"field"`
	// Mapping maps from an OIDC claim to Lattice organization UUIDs.
	Mapping map[string][]uuid.UUID `json:"mapping"`
	// AssignDefault ensures the default org is always included for every user,
	// regardless of their claims. This preserves legacy behavior.
	AssignDefault bool `json:"organization_assign_default"`
}
