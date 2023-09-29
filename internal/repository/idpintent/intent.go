package idpintent

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	StartedEventType       = instanceEventTypePrefix + "started"
	SucceededEventType     = instanceEventTypePrefix + "succeeded"
	SAMLSucceededEventType = instanceEventTypePrefix + "saml.succeeded"
	SAMLRequestEventType   = instanceEventTypePrefix + "saml.requested"
	LDAPSucceededEventType = instanceEventTypePrefix + "ldap.succeeded"
	FailedEventType        = instanceEventTypePrefix + "failed"
)

type StartedEvent struct {
	eventstore.BaseEvent `json:"-"`

	SuccessURL *url.URL `json:"successURL"`
	FailureURL *url.URL `json:"failureURL"`
	IDPID      string   `json:"idpId"`
}

func NewStartedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	successURL,
	failureURL *url.URL,
	idpID string,
) *StartedEvent {
	return &StartedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			StartedEventType,
		),
		SuccessURL: successURL,
		FailureURL: failureURL,
		IDPID:      idpID,
	}
}

func (e *StartedEvent) Data() interface{} {
	return e
}

func (e *StartedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func StartedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &StartedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sf3f1", "unable to unmarshal event")
	}

	return e, nil
}

type SucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPUser     []byte `json:"idpUser"`
	IDPUserID   string `json:"idpUserId,omitempty"`
	IDPUserName string `json:"idpUserName,omitempty"`
	UserID      string `json:"userId,omitempty"`

	IDPAccessToken *crypto.CryptoValue `json:"idpAccessToken,omitempty"`
	IDPIDToken     string              `json:"idpIdToken,omitempty"`
}

func NewSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpUser []byte,
	idpUserID,
	idpUserName,
	userID string,
	idpAccessToken *crypto.CryptoValue,
	idpIDToken string,
) *SucceededEvent {
	return &SucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SucceededEventType,
		),
		IDPUser:        idpUser,
		IDPUserID:      idpUserID,
		IDPUserName:    idpUserName,
		UserID:         userID,
		IDPAccessToken: idpAccessToken,
		IDPIDToken:     idpIDToken,
	}
}

func (e *SucceededEvent) Data() interface{} {
	return e
}

func (e *SucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SucceededEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &SucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-HBreq", "unable to unmarshal event")
	}

	return e, nil
}

type SAMLSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPUser     []byte `json:"idpUser"`
	IDPUserID   string `json:"idpUserId,omitempty"`
	IDPUserName string `json:"idpUserName,omitempty"`
	UserID      string `json:"userId,omitempty"`

	Assertion *crypto.CryptoValue `json:"assertion,omitempty"`
}

func NewSAMLSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpUser []byte,
	idpUserID,
	idpUserName,
	userID string,
	assertion *crypto.CryptoValue,
) *SAMLSucceededEvent {
	return &SAMLSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLSucceededEventType,
		),
		IDPUser:     idpUser,
		IDPUserID:   idpUserID,
		IDPUserName: idpUserName,
		UserID:      userID,
		Assertion:   assertion,
	}
}

func (e *SAMLSucceededEvent) Data() interface{} {
	return e
}

func (e *SAMLSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SAMLSucceededEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &SAMLSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-l4tw23y6lq", "unable to unmarshal event")
	}

	return e, nil
}

type SAMLRequestEvent struct {
	eventstore.BaseEvent `json:"-"`

	RequestID string `json:"requestId"`
}

func NewSAMLRequestEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	requestID string,
) *SAMLRequestEvent {
	return &SAMLRequestEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLRequestEventType,
		),
		RequestID: requestID,
	}
}

func (e *SAMLRequestEvent) Data() interface{} {
	return e
}

func (e *SAMLRequestEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SAMLRequestEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &SAMLRequestEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-l85678vwlf", "unable to unmarshal event")
	}

	return e, nil
}

type LDAPSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPUser     []byte `json:"idpUser"`
	IDPUserID   string `json:"idpUserId,omitempty"`
	IDPUserName string `json:"idpUserName,omitempty"`
	UserID      string `json:"userId,omitempty"`

	EntryAttributes map[string][]string `json:"user,omitempty"`
}

func NewLDAPSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpUser []byte,
	idpUserID,
	idpUserName,
	userID string,
	attributes map[string][]string,
) *LDAPSucceededEvent {
	return &LDAPSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LDAPSucceededEventType,
		),
		IDPUser:         idpUser,
		IDPUserID:       idpUserID,
		IDPUserName:     idpUserName,
		UserID:          userID,
		EntryAttributes: attributes,
	}
}

func (e *LDAPSucceededEvent) Data() interface{} {
	return e
}

func (e *LDAPSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func LDAPSucceededEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &LDAPSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-HBreq", "unable to unmarshal event")
	}

	return e, nil
}

type FailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Reason string `json:"reason,omitempty"`
}

func NewFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	reason string,
) *FailedEvent {
	return &FailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FailedEventType,
		),
		Reason: reason,
	}
}

func (e *FailedEvent) Data() interface{} {
	return e
}

func (e *FailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func FailedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &FailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sfer3", "unable to unmarshal event")
	}

	return e, nil
}
