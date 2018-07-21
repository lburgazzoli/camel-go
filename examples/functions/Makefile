# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: clean build
build: 
		go build -buildmode=plugin -o filter4.so  function_filter4.go
		go build -buildmode=plugin -o set_header.so function_set_header.go
		go build -buildmode=plugin -o trace.so function_trace.go
clean: 
		go clean
		rm -f *.so