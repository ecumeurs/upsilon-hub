package endpoint

import (
	"fmt"

	"github.com/ecumeurs/upsiloncli/internal/api"
	"github.com/ecumeurs/upsiloncli/internal/session"
)

// AuthLogin implements Endpoint for POST /api/v1/auth/login.
type AuthLogin struct{}

func (e *AuthLogin) Name() string        { return "auth_login" }
func (e *AuthLogin) Description() string { return "Authenticate and receive JWT" }
func (e *AuthLogin) Method() string      { return "POST" }
func (e *AuthLogin) Path() string        { return "/api/v1/auth/login" }
func (e *AuthLogin) Auth() bool          { return false }
func (e *AuthLogin) Params() []Param {
	return []Param{
		{Name: "account_name", Hint: "tactical identifier", Required: true},
		{Name: "password", Hint: "min 15 chars", Required: true},
	}
}

func (e *AuthLogin) Next() []string {
	return []string{"matchmaking_status", "matchmaking_join", "profile_get"}
}

func (e *AuthLogin) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Post(e.Path(), map[string]string{
		"account_name": inputs["account_name"],
		"password":     inputs["password"],
	})
	if err != nil {
		return err
	}

	// Extract token and user_id from response
	if data, ok := resp.Data.(map[string]interface{}); ok {
		if token, ok := data["token"].(string); ok && token != "" {
			sess.SetToken(token)
			client.Printer.System("JWT cached from login response.")
		}
		if user, ok := data["user"].(map[string]interface{}); ok {
			if uid, ok := user["id"].(string); ok {
				sess.Set("user_id", uid)
			}
			if name, ok := user["account_name"].(string); ok {
				sess.Set("account_name", name)
			}
		}
	}
	client.Printer.System(fmt.Sprintf("Login: %s", resp.Message))
	return nil
}

// AuthRegister implements Endpoint for POST /api/v1/auth/register.
type AuthRegister struct{}

func (e *AuthRegister) Name() string        { return "auth_register" }
func (e *AuthRegister) Description() string { return "Create account + initial roster, receive JWT" }
func (e *AuthRegister) Method() string      { return "POST" }
func (e *AuthRegister) Path() string        { return "/api/v1/auth/register" }
func (e *AuthRegister) Auth() bool          { return false }
func (e *AuthRegister) Params() []Param {
	return []Param{
		{Name: "account_name", Hint: "display name", Required: true},
		{Name: "email", Hint: "unique email", Required: true},
		{Name: "password", Hint: "min 15 chars", Required: true},
		{Name: "password_confirmation", Hint: "must match password", Required: true},
		{Name: "full_address", Hint: "residential address", Required: true},
		{Name: "birth_date", Hint: "ISO8601 e.g. 1990-01-15", Required: true},
	}
}

func (e *AuthRegister) Next() []string {
	return []string{"profile_get", "matchmaking_join"}
}

func (e *AuthRegister) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Post(e.Path(), inputs)
	if err != nil {
		return err
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		if token, ok := data["token"].(string); ok && token != "" {
			sess.SetToken(token)
			client.Printer.System("JWT cached from register response.")
		}
		if user, ok := data["user"].(map[string]interface{}); ok {
			if uid, ok := user["id"].(string); ok {
				sess.Set("user_id", uid)
			}
			if name, ok := user["account_name"].(string); ok {
				sess.Set("account_name", name)
			}
		}
	}
	client.Printer.System(fmt.Sprintf("Register: %s", resp.Message))
	return nil
}

// AuthLogout implements Endpoint for POST /api/v1/auth/logout.
type AuthLogout struct{}

func (e *AuthLogout) Name() string        { return "auth_logout" }
func (e *AuthLogout) Description() string { return "Terminate session and revoke token" }
func (e *AuthLogout) Method() string      { return "POST" }
func (e *AuthLogout) Path() string        { return "/api/v1/auth/logout" }
func (e *AuthLogout) Auth() bool          { return true }
func (e *AuthLogout) Params() []Param     { return nil }

func (e *AuthLogout) Next() []string {
	return []string{"auth_login"}
}

func (e *AuthLogout) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Post(e.Path(), nil)
	if err != nil {
		return err
	}
	sess.ClearAll()
	client.Printer.System(fmt.Sprintf("Logout: %s — Session cleared.", resp.Message))
	return nil
}

// AuthUpdate implements Endpoint for POST /api/v1/auth/update.
type AuthUpdate struct{}

func (e *AuthUpdate) Name() string        { return "auth_update" }
func (e *AuthUpdate) Description() string { return "Update account identity (nickname, email, address)" }
func (e *AuthUpdate) Method() string      { return "POST" }
func (e *AuthUpdate) Path() string        { return "/api/v1/auth/update" }
func (e *AuthUpdate) Auth() bool          { return true }
func (e *AuthUpdate) Params() []Param {
	return []Param{
		{Name: "account_name", Hint: "new display name (optional)", ContextKey: "account_name"},
		{Name: "email", Hint: "new email (optional)"},
		{Name: "full_address", Hint: "new address (optional)"},
		{Name: "birth_date", Hint: "new birth date ISO8601 (optional)"},
	}
}

func (e *AuthUpdate) Next() []string {
	return []string{"profile_get"}
}

func (e *AuthUpdate) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	// Filter out empty optional fields
	body := make(map[string]string)
	for k, v := range inputs {
		if v != "" {
			body[k] = v
		}
	}
	resp, err := client.Post(e.Path(), body)
	if err != nil {
		return err
	}
	client.Printer.System(fmt.Sprintf("Update: %s", resp.Message))
	return nil
}

// AuthPassword implements Endpoint for POST /api/v1/auth/password.
type AuthPassword struct{}

func (e *AuthPassword) Name() string        { return "auth_password" }
func (e *AuthPassword) Description() string { return "Change account password" }
func (e *AuthPassword) Method() string      { return "POST" }
func (e *AuthPassword) Path() string        { return "/api/v1/auth/password" }
func (e *AuthPassword) Auth() bool          { return true }
func (e *AuthPassword) Params() []Param {
	return []Param{
		{Name: "current_password", Hint: "current password", Required: true},
		{Name: "password", Hint: "new password (min 15 chars)", Required: true},
		{Name: "password_confirmation", Hint: "confirm new password", Required: true},
	}
}

func (e *AuthPassword) Next() []string {
	return []string{"profile_get", "auth_login"}
}

func (e *AuthPassword) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Post(e.Path(), inputs)
	if err != nil {
		return err
	}
	client.Printer.System(fmt.Sprintf("Password: %s", resp.Message))
	return nil
}

// AuthExport implements Endpoint for GET /api/v1/auth/export.
type AuthExport struct{}

func (e *AuthExport) Name() string        { return "auth_export" }
func (e *AuthExport) Description() string { return "GDPR data portability export" }
func (e *AuthExport) Method() string      { return "GET" }
func (e *AuthExport) Path() string        { return "/api/v1/auth/export" }
func (e *AuthExport) Auth() bool          { return true }
func (e *AuthExport) Params() []Param     { return nil }

func (e *AuthExport) Next() []string {
	return nil
}

func (e *AuthExport) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	_, err := client.Get(e.Path())
	return err
}

// AuthDelete implements Endpoint for DELETE /api/v1/auth/delete.
type AuthDelete struct{}

func (e *AuthDelete) Name() string        { return "auth_delete" }
func (e *AuthDelete) Description() string { return "GDPR right to be forgotten — anonymize & delete" }
func (e *AuthDelete) Method() string      { return "DELETE" }
func (e *AuthDelete) Path() string        { return "/api/v1/auth/delete" }
func (e *AuthDelete) Auth() bool          { return true }
func (e *AuthDelete) Params() []Param     { return nil }

func (e *AuthDelete) Next() []string {
	return []string{"auth_register"}
}

func (e *AuthDelete) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Delete(e.Path())
	if err != nil {
		return err
	}
	sess.ClearAll()
	client.Printer.System(fmt.Sprintf("Delete: %s — Session cleared.", resp.Message))
	return nil
}
