package sendlimiter

import (
	"encoding/json"
	"io"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

// ChannelCriterion is a criterion for a send limiter that is based on a channel.
type ChannelCriterion struct {
	// ExcludeRoleIDs is a list of role IDs that are ignored for the criterion.
	ExcludeRoleIDs []discord.RoleID `json:"exclude_role_ids"`

	// IncludeRoleIDs is a list of role IDs required for the criterion.
	IncludeRoleIDs []discord.RoleID `json:"include_role_ids"`

	// ChannelCooldown is the cooldown for a channel.
	ChannelCooldown time.Duration `json:"channel_cooldown"`

	// UserCooldown is the cooldown for a user.
	UserCooldown time.Duration `json:"user_cooldown"`

	// MinimumContentLength is the minimum content length for a message.
	MinimumContentLength int `json:"minimum_content_length"`

	// MaximumContentLength is the maximum content length for a message.
	MaximumContentLength int `json:"maximum_content_length"`

	// MinimumAttachmentLength is the minimum attachment length for a message.
	MinimumAttachmentLength int `json:"minimum_attachment_length"`

	// MaximumAttachmentLength is the maximum attachment length for a message.
	MaximumAttachmentLength int `json:"maximum_attachment_length"`
}

// ChannelCriteria is a map of channel IDs to channel criteria.
type ChannelCriteria map[discord.ChannelID]ChannelCriterion

// UnmarshalChannelCriteria unmarshals a channel criteria from a reader.
func UnmarshalChannelCriteria(r io.Reader) (ChannelCriteria, error) {
	var criteria ChannelCriteria
	if err := json.NewDecoder(r).Decode(&criteria); err != nil {
		return nil, err
	}

	return criteria, nil
}
