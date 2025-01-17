package identity

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

var _ json.Marshaler = &UserAssignedList{}

type UserAssignedList struct {
	Type        Type     `json:"type"`
	IdentityIds []string `json:"userAssignedIdentities"`
}

func (s *UserAssignedList) MarshalJSON() ([]byte, error) {
	// we use a custom marshal function here since we can only send the Type / UserAssignedIdentities field
	identityType := TypeNone
	userAssignedIdentityIds := []string{}

	if s != nil {
		if s.Type == TypeUserAssigned {
			identityType = TypeUserAssigned
			userAssignedIdentityIds = s.IdentityIds
		}
	}

	out := map[string]interface{}{
		"type":                   string(identityType),
		"userAssignedIdentities": userAssignedIdentityIds,
	}
	return json.Marshal(out)
}

// ExpandUserAssignedList expands the schema input into a UserAssignedList struct
func ExpandUserAssignedList(input []interface{}) (*UserAssignedList, error) {
	identityType := TypeNone
	identityIds := make([]string, 0)

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})
		typeRaw := raw["type"].(string)
		if typeRaw == string(TypeUserAssigned) {
			identityType = TypeUserAssigned
		}

		identityIdsRaw := raw["identity_ids"].(*schema.Set).List()
		for _, v := range identityIdsRaw {
			identityIds = append(identityIds, v.(string))
		}
	}

	if len(identityIds) > 0 && identityType != TypeUserAssigned {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q", string(TypeUserAssigned))
	}

	return &UserAssignedList{
		Type:        identityType,
		IdentityIds: identityIds,
	}, nil
}

// FlattenUserAssignedList turns a UserAssignedList into a []interface{}
func FlattenUserAssignedList(input *UserAssignedList) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	input.Type = normalizeType(input.Type)

	if input.Type != TypeUserAssigned {
		return &[]interface{}{}, nil
	}

	identityIds := make([]string, 0)
	for _, raw := range input.IdentityIds {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", raw, err)
		}
		identityIds = append(identityIds, id.ID())
	}

	return &[]interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
		},
	}, nil
}
