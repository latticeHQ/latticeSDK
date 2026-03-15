// Package deployments provides deployment management services for the Lattice Runtime API.
package deployments

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Deployment is the JSON representation of a Lattice deployment instance.
type Deployment struct {
	// ID is the unique identifier for the deployment.
	ID uuid.UUID `json:"id"`
	// Slug is the URL-safe identifier for the deployment.
	Slug string `json:"slug"`
	// DisplayName is the human-readable name for the deployment.
	DisplayName string `json:"display_name"`
	// Description is a description of the deployment.
	Description string `json:"description"`
	// AccessURL is the canonical URL for this deployment.
	AccessURL string `json:"access_url"`
	// IsDefault indicates whether this is the default deployment.
	IsDefault bool `json:"is_default"`
	// LogoURL is the URL for the deployment's logo.
	LogoURL string `json:"logo_url"`
	// ApplicationName is the display name of the application for this deployment.
	ApplicationName string `json:"application_name"`
	// OIDCConfig is the per-deployment OIDC configuration as raw JSON.
	OIDCConfig *json.RawMessage `json:"oidc_config,omitempty"`
	// GithubOAuthConfig is the per-deployment GitHub OAuth configuration as raw JSON.
	GithubOAuthConfig *json.RawMessage `json:"github_oauth_config,omitempty"`
	// CreatedAt is the time the deployment was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the deployment was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateDeploymentRequest provides options for creating a new deployment instance.
type CreateDeploymentRequest struct {
	// Slug is the URL-safe identifier for the deployment.
	Slug string `json:"slug"`
	// DisplayName is the human-readable name for the deployment.
	DisplayName string `json:"display_name,omitempty"`
	// Description is a description of the deployment.
	Description string `json:"description,omitempty"`
	// AccessURL is the canonical URL for this deployment (e.g., "https://acme-corp.com").
	AccessURL string `json:"access_url"`
	// LogoURL is the URL for the deployment's logo.
	LogoURL string `json:"logo_url,omitempty"`
	// ApplicationName is the display name of the application for this deployment.
	ApplicationName string `json:"application_name,omitempty"`
}

// UpdateDeploymentRequest provides options for updating an existing deployment.
type UpdateDeploymentRequest struct {
	// Slug is the new URL-safe identifier for the deployment.
	Slug string `json:"slug,omitempty"`
	// DisplayName is the new human-readable name for the deployment.
	DisplayName string `json:"display_name,omitempty"`
	// Description is the new description for the deployment.
	// Use a pointer to distinguish between omitting the field and setting it to empty.
	Description *string `json:"description,omitempty"`
	// AccessURL is the new canonical URL for the deployment.
	AccessURL string `json:"access_url,omitempty"`
	// LogoURL is the new logo URL.
	// Use a pointer to distinguish between omitting the field and setting it to empty.
	LogoURL *string `json:"logo_url,omitempty"`
	// ApplicationName is the new application display name.
	// Use a pointer to distinguish between omitting the field and setting it to empty.
	ApplicationName *string `json:"application_name,omitempty"`
	// OIDCConfig is the per-deployment OIDC configuration as raw JSON.
	// Set to null to use global config from env vars.
	OIDCConfig *json.RawMessage `json:"oidc_config,omitempty"`
	// GithubOAuthConfig is the per-deployment GitHub OAuth configuration as raw JSON.
	// Set to null to use global config from env vars.
	GithubOAuthConfig *json.RawMessage `json:"github_oauth_config,omitempty"`
}

// DeploymentConfig contains both the deployment values and how they are set.
type DeploymentConfig struct {
	// Values contains the raw deployment configuration.
	Values json.RawMessage `json:"config,omitempty"`
	// Options contains the available configuration options.
	Options json.RawMessage `json:"options,omitempty"`
}

// AgentConnectionLatencyMS contains agent connection latency percentiles.
type AgentConnectionLatencyMS struct {
	// P50 is the 50th-percentile (median) latency in milliseconds.
	P50 float64 `json:"p50"`
	// P95 is the 95th-percentile latency in milliseconds.
	P95 float64 `json:"p95"`
}

// AgentDeploymentStats contains statistics about agents in the deployment.
type AgentDeploymentStats struct {
	// Pending is the number of agents in the pending state.
	Pending int64 `json:"pending"`
	// Building is the number of agents currently building.
	Building int64 `json:"building"`
	// Running is the number of agents currently running.
	Running int64 `json:"running"`
	// Failed is the number of agents in the failed state.
	Failed int64 `json:"failed"`
	// Stopped is the number of agents that have been stopped.
	Stopped int64 `json:"stopped"`
	// ConnectionLatencyMS contains connection latency percentiles.
	ConnectionLatencyMS AgentConnectionLatencyMS `json:"connection_latency_ms"`
	// RxBytes is the total number of bytes received by agents.
	RxBytes int64 `json:"rx_bytes"`
	// TxBytes is the total number of bytes transmitted by agents.
	TxBytes int64 `json:"tx_bytes"`
}

// SessionCountDeploymentStats contains session counts by type.
type SessionCountDeploymentStats struct {
	// VSCode is the number of active VS Code sessions.
	VSCode int64 `json:"vscode"`
	// SSH is the number of active SSH sessions.
	SSH int64 `json:"ssh"`
	// JetBrains is the number of active JetBrains sessions.
	JetBrains int64 `json:"jetbrains"`
	// ReconnectingPTY is the number of active reconnecting PTY sessions.
	ReconnectingPTY int64 `json:"reconnecting_pty"`
}

// DeploymentStats contains aggregate statistics about the deployment.
type DeploymentStats struct {
	// AggregatedFrom is the start of the aggregation window.
	AggregatedFrom time.Time `json:"aggregated_from"`
	// CollectedAt is the time the stats were collected.
	CollectedAt time.Time `json:"collected_at"`
	// NextUpdateAt is the time when the next batch of stats will be updated.
	NextUpdateAt time.Time `json:"next_update_at"`
	// Agents contains agent-related statistics.
	Agents AgentDeploymentStats `json:"agents"`
	// SessionCount contains session count statistics by type.
	SessionCount SessionCountDeploymentStats `json:"session_count"`
}

// Entitlement represents whether a feature is licensed.
type Entitlement string

const (
	// EntitlementEntitled indicates the feature is fully licensed.
	EntitlementEntitled Entitlement = "entitled"
	// EntitlementGracePeriod indicates the feature license is in its grace period.
	EntitlementGracePeriod Entitlement = "grace_period"
	// EntitlementNotEntitled indicates the feature is not licensed.
	EntitlementNotEntitled Entitlement = "not_entitled"
)

// FeatureName represents the internal name of a feature.
type FeatureName string

const (
	FeatureUserLimit                  FeatureName = "user_limit"
	FeatureAuditLog                   FeatureName = "audit_log"
	FeatureBrowserOnly                FeatureName = "browser_only"
	FeatureSCIM                       FeatureName = "scim"
	FeatureTemplateRBAC               FeatureName = "template_rbac"
	FeatureUserRoleManagement         FeatureName = "user_role_management"
	FeatureHighAvailability           FeatureName = "high_availability"
	FeatureMultipleExternalAuth       FeatureName = "multiple_external_auth"
	FeatureExternalProvisionerDaemons FeatureName = "external_provisioner_daemons"
	FeatureAppearance                 FeatureName = "appearance"
	FeatureAdvancedTemplateScheduling FeatureName = "advanced_template_scheduling"
	FeatureAgentProxy                 FeatureName = "agent_proxy"
	FeatureExternalTokenEncryption    FeatureName = "external_token_encryption"
	FeatureAgentBatchActions          FeatureName = "agent_batch_actions"
	FeatureAccessControl              FeatureName = "access_control"
	FeatureControlSharedPorts         FeatureName = "control_shared_ports"
	FeatureCustomRoles                FeatureName = "custom_roles"
	FeatureMultipleOrganizations      FeatureName = "multiple_organizations"
	FeatureAIBridge                   FeatureName = "aibridge"
)

// Feature describes the state of a licensed feature.
type Feature struct {
	// Entitlement is the entitlement status of the feature.
	Entitlement Entitlement `json:"entitlement"`
	// Enabled indicates whether the feature is currently enabled.
	Enabled bool `json:"enabled"`
	// Limit is the maximum allowed value for the feature (e.g., user count).
	Limit *int64 `json:"limit,omitempty"`
	// Actual is the current usage value for the feature.
	Actual *int64 `json:"actual,omitempty"`
}

// FeatureSet represents a grouping of features assigned to a license tier.
type FeatureSet string

const (
	FeatureSetNone       FeatureSet = ""
	FeatureSetFree       FeatureSet = "free"
	FeatureSetTeam       FeatureSet = "team"
	FeatureSetBusiness   FeatureSet = "business"
	FeatureSetEnterprise FeatureSet = "enterprise"
)

// Entitlements contains the full set of feature entitlements for a deployment.
type Entitlements struct {
	// Features maps feature names to their entitlement state.
	Features map[FeatureName]Feature `json:"features"`
	// Warnings contains non-fatal entitlement warnings.
	Warnings []string `json:"warnings"`
	// Errors contains entitlement errors.
	Errors []string `json:"errors"`
	// HasLicense indicates whether a valid license is present.
	HasLicense bool `json:"has_license"`
	// Trial indicates whether the deployment is in trial mode.
	Trial bool `json:"trial"`
	// RequireTelemetry indicates whether telemetry is required by the license.
	RequireTelemetry bool `json:"require_telemetry"`
	// RefreshedAt is the time entitlements were last refreshed.
	RefreshedAt time.Time `json:"refreshed_at"`
}

// NavigationConfig stores per-deployment sidebar navigation and theme customization.
type NavigationConfig struct {
	// NavLabelOverrides maps navigation item keys to custom display labels.
	NavLabelOverrides map[string]string `json:"nav_label_overrides,omitempty"`
	// SectionLabelOverrides maps section keys to custom display labels.
	SectionLabelOverrides map[string]string `json:"section_label_overrides,omitempty"`
	// HiddenNavItems maps navigation item keys to their hidden state.
	HiddenNavItems map[string]bool `json:"hidden_nav_items,omitempty"`
	// HideHelpCenter hides the help center from the sidebar.
	HideHelpCenter bool `json:"hide_help_center,omitempty"`
	// ThemePreset is the per-deployment theme.
	// Valid values: "" (default/platform), "healthcare", "finance", "tech", "education".
	ThemePreset string `json:"theme_preset,omitempty"`
}

// BannerConfig configures an announcement banner displayed in the dashboard.
type BannerConfig struct {
	// Enabled indicates whether the banner is displayed.
	Enabled bool `json:"enabled"`
	// Message is the banner text.
	Message string `json:"message,omitempty"`
	// BackgroundColor is the CSS background color for the banner.
	BackgroundColor string `json:"background_color,omitempty"`
}

// LinkConfig describes a support or navigation link.
type LinkConfig struct {
	// Name is the display text for the link.
	Name string `json:"name"`
	// Target is the URL the link points to.
	Target string `json:"target"`
	// Icon is the icon identifier. Valid values: "bug", "chat", "docs".
	Icon string `json:"icon"`
}

// AppearanceConfig contains the visual appearance settings for the dashboard.
type AppearanceConfig struct {
	// ApplicationName is the display name shown in the dashboard.
	ApplicationName string `json:"application_name"`
	// LogoURL is the URL of the logo shown in the dashboard.
	LogoURL string `json:"logo_url"`
	// DocsURL is the URL for the documentation link.
	DocsURL string `json:"docs_url"`
	// ServiceBanner is the legacy banner configuration.
	// Deprecated: Use AnnouncementBanners instead.
	ServiceBanner BannerConfig `json:"service_banner"`
	// AnnouncementBanners is the list of active announcement banners.
	AnnouncementBanners []BannerConfig `json:"announcement_banners"`
	// SupportLinks is the list of support links shown in the dashboard.
	SupportLinks []LinkConfig `json:"support_links,omitempty"`
	// NavigationConfig contains sidebar navigation customization.
	NavigationConfig NavigationConfig `json:"navigation_config"`
}

// UpdateAppearanceConfig provides options for updating dashboard appearance.
type UpdateAppearanceConfig struct {
	// ApplicationName is the display name shown in the dashboard.
	ApplicationName string `json:"application_name"`
	// LogoURL is the URL of the logo shown in the dashboard.
	LogoURL string `json:"logo_url"`
	// ServiceBanner is the legacy banner configuration.
	// Deprecated: Use AnnouncementBanners instead.
	ServiceBanner BannerConfig `json:"service_banner"`
	// AnnouncementBanners is the list of active announcement banners.
	AnnouncementBanners []BannerConfig `json:"announcement_banners"`
	// NavigationConfig contains sidebar navigation customization.
	NavigationConfig NavigationConfig `json:"navigation_config"`
}

// BuildInfoResponse contains build information for the Lattice instance.
type BuildInfoResponse struct {
	// ExternalURL references the current Lattice version.
	// For production builds, this links to a release. For development builds, this links to a commit.
	ExternalURL string `json:"external_url"`
	// Version is the semantic version of the build.
	Version string `json:"version"`
	// DashboardURL is the URL to the deployment's dashboard.
	DashboardURL string `json:"dashboard_url"`
	// Telemetry indicates whether telemetry is enabled.
	Telemetry bool `json:"telemetry"`
	// AgentProxy indicates whether this instance is an agent proxy.
	AgentProxy bool `json:"agent_proxy"`
	// SidecarAPIVersion is the current version of the Sidecar API.
	SidecarAPIVersion string `json:"sidecar_api_version"`
	// ProvisionerAPIVersion is the current version of the Provisioner API.
	ProvisionerAPIVersion string `json:"provisioner_api_version"`
	// UpgradeMessage is the message displayed to users when an outdated client is detected.
	UpgradeMessage string `json:"upgrade_message"`
	// DeploymentID is the unique identifier for this deployment.
	DeploymentID string `json:"deployment_id"`
}

// SSHConfigResponse contains SSH configuration for the Lattice instance.
type SSHConfigResponse struct {
	// HostnamePrefix is the config-ssh hostname prefix.
	HostnamePrefix string `json:"hostname_prefix"`
	// SSHConfigOptions contains additional SSH config options as key-value pairs.
	SSHConfigOptions map[string]string `json:"ssh_config_options"`
}

// AppHostResponse contains the application wildcard hostname.
type AppHostResponse struct {
	// Host is the externally accessible wildcard hostname for the Lattice instance.
	// e.g., "*--apps.latticeruntime.com".
	Host string `json:"host"`
}

// Experiment represents a named experiment that can be enabled on a deployment.
type Experiment string

const (
	// ExperimentExample is a placeholder experiment for testing.
	ExperimentExample Experiment = "example"
	// ExperimentAutoFillParameters enables auto-fill for template parameters.
	ExperimentAutoFillParameters Experiment = "auto-fill-parameters"
	// ExperimentNotifications enables SMTP and webhook notification delivery.
	ExperimentNotifications Experiment = "notifications"
	// ExperimentAgentUsage enables agent usage tracking.
	ExperimentAgentUsage Experiment = "agent-usage"
)

// Experiments is a list of experiments enabled on a deployment.
type Experiments []Experiment

// Enabled returns true if the given experiment is in the list.
func (e Experiments) Enabled(ex Experiment) bool {
	for _, v := range e {
		if v == ex {
			return true
		}
	}
	return false
}

// AvailableExperiments lists experiments that are safe for users to opt in to.
type AvailableExperiments struct {
	// Safe contains experiments that are safe for general use.
	Safe []Experiment `json:"safe"`
}

// DAUsResponse contains daily active user metrics.
type DAUsResponse struct {
	// Entries contains the DAU data points.
	Entries []DAUEntry `json:"entries"`
	// TZHourOffset is the timezone hour offset used for the query.
	TZHourOffset int `json:"tz_hour_offset"`
}

// DAUEntry represents a single day's active user count.
type DAUEntry struct {
	// Date is formatted as "2024-01-31". Timezone and time information is not included.
	Date string `json:"date"`
	// Amount is the number of active users on that date.
	Amount int `json:"amount"`
}

// UpdateCheckResponse contains information about the latest Lattice release.
type UpdateCheckResponse struct {
	// Current indicates whether the server version matches the latest release.
	Current bool `json:"current"`
	// Version is the semantic version of the latest release.
	Version string `json:"version"`
	// URL is the download link for the latest release.
	URL string `json:"url"`
}

// CryptoKeyFeature identifies the purpose of a cryptographic key.
type CryptoKeyFeature string

const (
	// CryptoKeyFeatureAgentAppsAPIKey is used for agent app API key signing.
	CryptoKeyFeatureAgentAppsAPIKey CryptoKeyFeature = "agent_apps_api_key"
	// CryptoKeyFeatureAgentAppsToken is used for agent app token signing.
	CryptoKeyFeatureAgentAppsToken CryptoKeyFeature = "agent_apps_token"
	// CryptoKeyFeatureOIDCConvert is used for OIDC token conversion.
	CryptoKeyFeatureOIDCConvert CryptoKeyFeature = "oidc_convert"
	// CryptoKeyFeatureTailnetResume is used for tailnet session resumption.
	CryptoKeyFeatureTailnetResume CryptoKeyFeature = "tailnet_resume"
)

// CryptoKey represents a cryptographic key used by the deployment.
type CryptoKey struct {
	// Feature identifies what the key is used for.
	Feature CryptoKeyFeature `json:"feature"`
	// Secret is the key material.
	Secret string `json:"secret"`
	// DeletesAt is the time after which the key will be deleted.
	DeletesAt time.Time `json:"deletes_at"`
	// Sequence is the key's rotation sequence number.
	Sequence int32 `json:"sequence"`
	// StartsAt is the time the key becomes active.
	StartsAt time.Time `json:"starts_at"`
}

// NotificationsConfig contains notification delivery configuration.
type NotificationsConfig struct {
	// MaxSendAttempts is the upper limit of attempts to send a notification.
	MaxSendAttempts int64 `json:"max_send_attempts"`
	// RetryInterval is the minimum time between retries.
	RetryInterval string `json:"retry_interval"`
	// Method is the delivery method (e.g., "smtp", "webhook").
	Method string `json:"method"`
}
