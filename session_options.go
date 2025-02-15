package sessionstore

// SessionOptionsInterface is an interface for session options
type SessionOptionsInterface interface {
	HasUserID() bool
	GetUserID() string
	SetUserID(userID string)

	HasIPAddress() bool
	GetIPAddress() string
	SetIPAddress(ipAddress string)

	HasUserAgent() bool
	GetUserAgent() string
	SetUserAgent(userAgent string)
}

// SessionOptions shortcut for NewSessionOptions
func SessionOptions() SessionOptionsInterface {
	return NewSessionOptions()
}

// NewSessionOptions creates a new session options
func NewSessionOptions() SessionOptionsInterface {
	return &sessionOptions{
		properties: make(map[string]any),
	}
}

// sessionOptions is a struct for session options
type sessionOptions struct {
	properties map[string]any
}

// HasIPAddress returns true if the session has an IP address
func (s *sessionOptions) HasIPAddress() bool {
	return s.hasProperty("ip_address")
}

// GetIPAddress returns the IP address of the session
func (s *sessionOptions) GetIPAddress() string {
	if !s.HasIPAddress() {
		return ""
	}

	return s.properties["ip_address"].(string)
}

// SetIPAddress sets the IP address of the session
func (s *sessionOptions) SetIPAddress(ipAddress string) {
	s.properties["ip_address"] = ipAddress
}

// HasUserAgent returns true if the session has a user agent
func (s *sessionOptions) HasUserAgent() bool {
	return s.hasProperty("user_agent")
}

// GetUserAgent returns the user agent of the session
func (s *sessionOptions) GetUserAgent() string {
	if !s.HasUserAgent() {
		return ""
	}

	return s.properties["user_agent"].(string)
}

// SetUserAgent sets the user agent of the session
func (s *sessionOptions) SetUserAgent(userAgent string) {
	s.properties["user_agent"] = userAgent
}

// HasUserID returns true if the session has a user ID
func (s *sessionOptions) HasUserID() bool {
	return s.hasProperty("user_id")
}

// GetUserID returns the user ID of the session
func (s *sessionOptions) GetUserID() string {
	if !s.HasUserID() {
		return ""
	}

	return s.properties["user_id"].(string)
}

// SetUserID sets the user ID of the session
func (s *sessionOptions) SetUserID(userID string) {
	s.properties["user_id"] = userID
}

// hasProperty returns true if the session has a property
func (s *sessionOptions) hasProperty(key string) bool {
	_, ok := s.properties[key]
	return ok
}
