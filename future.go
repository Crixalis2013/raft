// Copyright 2017 The raft Authors
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

package raft

type deferError struct {
	errCh chan error
}

func (d *deferError) init() {
	d.errCh = make(chan error, 1)
}

func (d *deferError) respond(err error) {
	d.errCh <- err
	close(d.errCh)
}

func (d *deferError) error() <-chan error {
	return d.errCh
}

type Future struct {
	deferError
	respCh chan interface{}
}

func newFuture() *Future {
	f := &Future{
		respCh: make(chan interface{}, 1),
	}
	f.init()
	return f
}

func (f *Future) respond(resp interface{}, err error) {
	if err == nil {
		f.respCh <- resp
		close(f.respCh)
	} else {
		f.deferError.respond(err)
	}
}

func (f *Future) Response() (resp interface{}, err error) {
	select {
	case err = <-f.error():
		return
	case resp = <-f.respCh:
		return
	}
}
func (f *Future) AsyncResponse() (respCh <-chan interface{}, errCh <-chan error) {
	return f.respCh, f.errCh
}
