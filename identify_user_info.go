package analytics

import (
	"errors"

	"github.com/common-fate/analytics-go/acore"
)

func init() {
	registerEvent(&UserInfo{})
}

type Role string

const (
	RoleEndUser Role = "end_user"
	RoleAdmin   Role = "admin"
)

type UserInfo struct {
	ID         string `json:"-"`
	GroupCount int    `json:"group_count"`
	IsAdmin    bool   `json:"-"`
}

func (d *UserInfo) userID() string { return "" }

func (d *UserInfo) Type() string { return "cf:identify:user_info" }

func (d *UserInfo) EmittedWhen() string { return "Access Request created/updated" }

func (d *UserInfo) marshalEvent(ctx marshalContext) ([]acore.Message, error) {
	// role is populated automatically based on IsAdmin.
	// In future we might have additional role types rather than
	// just 'admin' and 'end user'
	role := RoleEndUser
	if d.IsAdmin {
		role = RoleAdmin
	}

	uid, ok := hash(d.ID, "usr")
	if !ok {
		return nil, errors.New("could not hash user ID")
	}

	id := acore.Identify{
		DistinctId: uid,
		Properties: eventToProperties(d).Set("role", role),
	}

	if ctx.DeploymentID != nil {
		id.Groups = acore.NewGroups().Set("deployment", *ctx.DeploymentID)
	}

	m := []acore.Message{id}

	return m, nil
}

func (d *UserInfo) fixture() {
	*d = UserInfo{
		ID:         "usr_123",
		IsAdmin:    true,
		GroupCount: 2,
	}
}
