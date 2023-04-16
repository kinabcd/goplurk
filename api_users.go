package goplurk

type APIUsers struct {
	client *Client
}

func (u *APIUsers) Me() (*User, error) {
	res := &User{}
	if err := u.client.Engine.CallAPIUnmarshal("/APP/Users/me", map[string]string{}, res); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
