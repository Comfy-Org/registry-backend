package mapper

import (
	"fmt"
	"regexp"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
)

func ApiCreatePublisherToDb(publisher *drip.Publisher, client *ent.Client) (*ent.PublisherCreate, error) {
	newPublisher := client.Publisher.Create()
	if publisher.Description != nil {
		newPublisher.SetDescription(*publisher.Description)
	}
	if publisher.Id != nil {
		newPublisher.SetID(*publisher.Id)
	}
	if publisher.Logo != nil {
		newPublisher.SetLogoURL(*publisher.Logo)
	}
	if publisher.Name != nil {
		newPublisher.SetName(*publisher.Name)
	}
	if publisher.SourceCodeRepo != nil {
		newPublisher.SetSourceCodeRepo(*publisher.SourceCodeRepo)
	}
	if publisher.Support != nil {
		newPublisher.SetSupportEmail(*publisher.Support)
	}
	if publisher.Website != nil {
		newPublisher.SetWebsite(*publisher.Website)
	}

	return newPublisher, nil
}

func ApiUpdatePublisherToUpdateFields(publisherId string, publisher *drip.Publisher, client *ent.Client) *ent.PublisherUpdateOne {
	update := client.Publisher.UpdateOneID(publisherId)
	if publisher.Description != nil {
		update.SetDescription(*publisher.Description)
	}
	if publisher.Logo != nil {
		update.SetLogoURL(*publisher.Logo)
	}
	if publisher.Name != nil {
		update.SetName(*publisher.Name)
	}
	if publisher.SourceCodeRepo != nil {
		update.SetSourceCodeRepo(*publisher.SourceCodeRepo)
	}
	if publisher.Support != nil {
		update.SetSupportEmail(*publisher.Support)
	}
	if publisher.Website != nil {
		update.SetWebsite(*publisher.Website)
	}

	return update
}

func ValidatePublisher(publisher *drip.Publisher) error {
	if publisher.Id != nil {
		if !IsValidPublisherID(*publisher.Id) {
			return fmt.Errorf("invalid publisher id")
		}
	}
	return nil
}

func IsValidPublisherID(publisherID string) bool {
	// Regular expression pattern for Publisher ID validation (lowercase letters only)
	pattern := "^[a-z][a-z0-9-]*$"
	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)
	// Check if the string matches the pattern
	return regex.MatchString(publisherID)
}

func DbPublisherToApiPublisher(publisher *ent.Publisher, public bool) *drip.Publisher {
	members := make([]drip.PublisherMember, 0)

	if publisher.Edges.PublisherPermissions != nil {
		for _, permission := range publisher.Edges.PublisherPermissions {
			if permission.Edges.User != nil {
				member := drip.PublisherMember{}
				// If the data is not public, include sensitive information.
				if !public {
					member.User = &drip.PublisherUser{
						Id:    ToStringPointer(permission.Edges.User.ID),
						Email: ToStringPointer(permission.Edges.User.Email),
						Name:  ToStringPointer(permission.Edges.User.Name),
					}
					member.Role = ToStringPointer(string(permission.Permission))
				} else {
					member.User = &drip.PublisherUser{
						Name: ToStringPointer(permission.Edges.User.Name),
					}
				}

				members = append(members, member)
			}
		}
	}

	return &drip.Publisher{
		Description:    &publisher.Description,
		Id:             &publisher.ID,
		Logo:           &publisher.LogoURL,
		Name:           &publisher.Name,
		SourceCodeRepo: &publisher.SourceCodeRepo,
		Support:        &publisher.SupportEmail,
		Website:        &publisher.Website,
		CreatedAt:      &publisher.CreateTime,
		Members:        &members,
		Status:         DbPublisherStatusToApiPublisherStatus(publisher.Status),
	}
}

func ToStringPointer(s string) *string {
	return &s
}

func DbPublisherStatusToApiPublisherStatus(status schema.PublisherStatusType) *drip.PublisherStatus {
	var apiStatus drip.PublisherStatus
	switch status {
	case schema.PublisherStatusTypeActive:
		apiStatus = drip.PublisherStatusActive
	case schema.PublisherStatusTypeBanned:
		apiStatus = drip.PublisherStatusBanned
	default:
		apiStatus = ""
	}
	return &apiStatus
}
