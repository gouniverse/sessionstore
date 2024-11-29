package sessionstore

import "errors"

type SessionQueryInterface interface {
	Validate() error

	IsCountOnly() bool

	Columns() []string
	SetColumns(columns []string) SessionQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) SessionQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) SessionQueryInterface

	HasExpiresAtGte() bool
	ExpiresAtGte() string
	SetExpiresAtGte(expiresAtGte string) SessionQueryInterface

	HasExpiresAtLte() bool
	ExpiresAtLte() string
	SetExpiresAtLte(expiresAtLte string) SessionQueryInterface

	HasID() bool
	ID() string
	SetID(id string) SessionQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) SessionQueryInterface

	HasKey() bool
	Key() string
	SetKey(key string) SessionQueryInterface

	HasUserID() bool
	UserID() string
	SetUserID(userID string) SessionQueryInterface

	HasUserIpAddress() bool
	UserIpAddress() string
	SetUserIpAddress(userIpAddress string) SessionQueryInterface

	HasUserAgent() bool
	UserAgent() string
	SetUserAgent(userAgent string) SessionQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) SessionQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) SessionQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) SessionQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) SessionQueryInterface

	HasCountOnly() bool
	SetCountOnly(countOnly bool) SessionQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withSoftDeleted bool) SessionQueryInterface
}

func SessionQuery() SessionQueryInterface {
	return &sessionQuery{
		properties: make(map[string]interface{}),
	}
}

var _ SessionQueryInterface = (*sessionQuery)(nil)

type sessionQuery struct {
	properties map[string]interface{}
}

func (q *sessionQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("Session query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("Session query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("Session query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("Session query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("Session query. limit cannot be negative")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("Session query. offset cannot be negative")
	}

	return nil
}

func (q *sessionQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *sessionQuery) SetColumns(columns []string) SessionQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *sessionQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *sessionQuery) IsCountOnly() bool {
	return q.hasProperty("count_only") && q.properties["count_only"].(bool)
}

func (q *sessionQuery) SetCountOnly(countOnly bool) SessionQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *sessionQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *sessionQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *sessionQuery) SetCreatedAtGte(createdAtGte string) SessionQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *sessionQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *sessionQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *sessionQuery) SetCreatedAtLte(createdAtLte string) SessionQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *sessionQuery) HasExpiresAtGte() bool {
	return q.hasProperty("expires_at_gte")
}

func (q *sessionQuery) ExpiresAtGte() string {
	return q.properties["expires_at_gte"].(string)
}

func (q *sessionQuery) SetExpiresAtGte(expiresAtGte string) SessionQueryInterface {
	q.properties["expires_at_gte"] = expiresAtGte
	return q
}

func (q *sessionQuery) HasExpiresAtLte() bool {
	return q.hasProperty("expires_at_lte")
}

func (q *sessionQuery) ExpiresAtLte() string {
	return q.properties["expires_at_lte"].(string)
}

func (q *sessionQuery) SetExpiresAtLte(expiresAtLte string) SessionQueryInterface {
	q.properties["expires_at_lte"] = expiresAtLte
	return q
}

func (q *sessionQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *sessionQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *sessionQuery) SetID(id string) SessionQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *sessionQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *sessionQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *sessionQuery) SetIDIn(idIn []string) SessionQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *sessionQuery) HasKey() bool {
	return q.hasProperty("key")
}

func (q *sessionQuery) Key() string {
	return q.properties["key"].(string)
}

func (q *sessionQuery) SetKey(key string) SessionQueryInterface {
	q.properties["key"] = key
	return q
}

func (q *sessionQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *sessionQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *sessionQuery) SetLimit(limit int) SessionQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *sessionQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *sessionQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *sessionQuery) SetOffset(offset int) SessionQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *sessionQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *sessionQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *sessionQuery) SetOrderBy(orderBy string) SessionQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *sessionQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

func (q *sessionQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}

	return q.properties["soft_deleted_included"].(bool)
}

func (q *sessionQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) SessionQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

func (q *sessionQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *sessionQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *sessionQuery) SetSortOrder(sortOrder string) SessionQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *sessionQuery) HasUserAgent() bool {
	return q.hasProperty("user_agent")
}

func (q *sessionQuery) UserAgent() string {
	return q.properties["user_agent"].(string)
}

func (q *sessionQuery) SetUserAgent(userAgent string) SessionQueryInterface {
	q.properties["user_agent"] = userAgent
	return q
}

func (q *sessionQuery) HasUserID() bool {
	return q.hasProperty("user_id")
}

func (q *sessionQuery) UserID() string {
	return q.properties["user_id"].(string)
}

func (q *sessionQuery) SetUserID(userID string) SessionQueryInterface {
	q.properties["user_id"] = userID
	return q
}

func (q *sessionQuery) HasUserIpAddress() bool {
	return q.hasProperty("user_ip_address")
}

func (q *sessionQuery) UserIpAddress() string {
	return q.properties["user_ip_address"].(string)
}

func (q *sessionQuery) SetUserIpAddress(userIpAddress string) SessionQueryInterface {
	q.properties["user_ip_address"] = userIpAddress
	return q
}

func (q *sessionQuery) hasProperty(key string) bool {
	_, ok := q.properties[key]
	return ok
}
