/*
 * Echotron
 * Copyright (C) 2022 The Echotron Devs
 *
 * Echotron is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Echotron is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package echotron

import (
	"encoding/json"
	"fmt"
)

// ChatAdministratorRights represents the rights of an administrator in a chat.
type ChatAdministratorRights struct {
	IsAnonymous          bool `json:"is_anonymous"`
	CanManageChat        bool `json:"can_manage_chat"`
	CanDeleteMessages    bool `json:"can_delete_messages"`
	CanManageVideo_chats bool `json:"can_manage_video_chats"`
	CanRestrictMembers   bool `json:"can_restrict_members"`
	CanPromoteMembers    bool `json:"can_promote_members"`
	CanChangeInfo        bool `json:"can_change_info"`
	CanInviteUsers       bool `json:"can_invite_users"`
	CanPostMessages      bool `json:"can_post_messages,omitempty"`
	CanEditMessages      bool `json:"can_edit_messages,omitempty"`
	CanPinMessages       bool `json:"can_pin_messages,omitempty"`
}

// SetMyDefaultAdministratorRightsOptions contains the optional parameters used by
// the SetMyDefaultAdministratorRights method.
type SetMyDefaultAdministratorRightsOptions struct {
	Rights      ChatAdministratorRights `query:"rights"`
	ForChannels bool                    `query:"for_channels"`
}

// GetMyDefaultAdministratorRightsOptions contains the optional parameters used by
// the GetMyDefaultAdministratorRights method.
type GetMyDefaultAdministratorRightsOptions struct {
	ForChannels bool `query:"for_channels"`
}

// SetMyDefaultAdministratorRights is used to change the default administrator rights
// requested by the bot when it's added as an administrator to groups or channels.
// These rights will be suggested to users, but they are are free to modify the list
// before adding the bot.
func (a API) SetMyDefaultAdministratorRights(opts SetMyDefaultAdministratorRightsOptions) (res APIResponseBool, err error) {
	var url = fmt.Sprintf(
		"%ssetMyDefaultAdministratorRights?%s",
		a.base,
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return
	}

	if err = json.Unmarshal(cnt, &res); err != nil {
		return
	}

	err = check(res)
	return
}

// GetMyDefaultAdministratorRights is used to get the current default administrator rights of the bot.
func (a API) GetMyDefaultAdministratorRights(opts GetMyDefaultAdministratorRightsOptions) (res APIResponseChatAdministratorRights, err error) {
	var url = fmt.Sprintf(
		"%sgetMyDefaultAdministratorRights?%s",
		a.base,
		querify(opts),
	)

	cnt, err := sendGetRequest(url)
	if err != nil {
		return
	}

	if err = json.Unmarshal(cnt, &res); err != nil {
		return
	}

	err = check(res)
	return
}