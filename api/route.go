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

package api

// ==========================
//
// Route
//
// ==========================

// Route --
type Route struct {
	Service
	id       string
	services []Service
}

// NewRoute --
func NewRoute(id string) *Route {
	return &Route{
		id:       id,
		services: make([]Service, 0),
	}
}

// ID --
func (route *Route) ID() string {
	return route.id
}

// AddService --
func (route *Route) AddService(service Service) {
	if service != nil {
		route.services = append(route.services, service)
	}
}

// Start --
func (route *Route) Start() {
	StartServices(route.services)
}

// Stop --
func (route *Route) Stop() {
	StopServices(route.services)
}
