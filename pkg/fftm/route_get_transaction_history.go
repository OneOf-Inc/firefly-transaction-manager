// Copyright © 2023 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fftm

import (
	"net/http"

	"github.com/hyperledger/firefly-common/pkg/ffapi"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly-transaction-manager/internal/persistence"
	"github.com/hyperledger/firefly-transaction-manager/internal/tmmsgs"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
)

var getTransactionHistory = func(m *manager) *ffapi.Route {
	route := &ffapi.Route{
		Name:   "getTransactionHistory",
		Path:   "/transactions/{transactionId}/history",
		Method: http.MethodGet,
		PathParams: []*ffapi.PathParam{
			{Name: "transactionId", Description: tmmsgs.APIParamTransactionID},
		},
		Description:     tmmsgs.APIEndpointGetTransactionHistory,
		JSONInputValue:  nil,
		JSONOutputValue: func() interface{} { return []*apitypes.TXHistoryRecord{} },
		JSONOutputCodes: []int{http.StatusOK},
	}
	if m.richQueryEnabled {
		route.FilterFactory = persistence.TXHistoryFilters
		route.JSONHandler = func(r *ffapi.APIRequest) (output interface{}, err error) {
			return r.FilterResult(m.persistence.RichQuery().ListTransactionHistory(r.Req.Context(), r.PP["transactionId"], r.Filter))
		}
	} else {
		route.JSONHandler = func(r *ffapi.APIRequest) (output interface{}, err error) {
			return nil, i18n.NewError(r.Req.Context(), tmmsgs.MsgOpNotSupportedWithoutRichQuery)
		}
	}
	return route
}
