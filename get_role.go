package analytics

// getRole returns a string name for the role.
// Roles are either "admin", or "end_user".
func getRole(isAdmin bool) string {
	if isAdmin {
		return "admin"
	} else {
		return "end_user"
	}
}
