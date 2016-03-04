// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main_test

import (
	"fmt"
	"net/http"

	. "gopkg.in/check.v1"

	. "github.com/ubuntu-core/snappy/cmd/snap"
)

func (s *SnapSuite) TestDisconnectHelp(c *C) {
	msg := `Usage:
  snap.test [OPTIONS] disconnect [<snap>:<plug>] [<snap>:<slot>]

The disconnect command disconnects a plug from a slot.
It may be called in the following ways:

$ snap disconnect <snap>:<plug> <snap>:<slot>

Disconnects the specific plug from the specific slot.

$ snap disconnect <snap>:<slot>

Disconnects any previously connected plugs from the provided slot.

$ snap disconnect <snap>

Disconnects all plugs from the provided snap.

Help Options:
  -h, --help               Show this help message
`
	rest, err := Parser().ParseArgs([]string{"disconnect", "--help"})
	c.Assert(err.Error(), Equals, msg)
	c.Assert(rest, DeepEquals, []string{})
}

func (s *SnapSuite) TestDisconnectExplicitEverything(c *C) {
	s.RedirectClientToTestServer(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.Method, Equals, "POST")
		c.Check(r.URL.Path, Equals, "/2.0/interfaces")
		c.Check(DecodedRequestBody(c, r), DeepEquals, map[string]interface{}{
			"action": "disconnect",
			"plugs": []interface{}{
				map[string]interface{}{
					"snap": "producer",
					"plug": "plug",
				},
			},
			"slots": []interface{}{
				map[string]interface{}{
					"snap": "consumer",
					"slot": "slot",
				},
			},
		})
		fmt.Fprintln(w, `{"type":"sync", "result":{}}`)
	})
	rest, err := Parser().ParseArgs([]string{"disconnect", "producer:plug", "consumer:slot"})
	c.Assert(err, IsNil)
	c.Assert(rest, DeepEquals, []string{})
	c.Assert(s.Stdout(), Equals, "")
	c.Assert(s.Stderr(), Equals, "")
}

func (s *SnapSuite) TestDisconnectEverythingFromSpecificSlot(c *C) {
	s.RedirectClientToTestServer(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.Method, Equals, "POST")
		c.Check(r.URL.Path, Equals, "/2.0/interfaces")
		c.Check(DecodedRequestBody(c, r), DeepEquals, map[string]interface{}{
			"action": "disconnect",
			"plugs": []interface{}{
				map[string]interface{}{
					"snap": "",
					"plug": "",
				},
			},
			"slots": []interface{}{
				map[string]interface{}{
					"snap": "consumer",
					"slot": "slot",
				},
			},
		})
		fmt.Fprintln(w, `{"type":"sync", "result":{}}`)
	})
	rest, err := Parser().ParseArgs([]string{"disconnect", "consumer:slot"})
	c.Assert(err, IsNil)
	c.Assert(rest, DeepEquals, []string{})
	c.Assert(s.Stdout(), Equals, "")
	c.Assert(s.Stderr(), Equals, "")
}

func (s *SnapSuite) TestDisconnectEverythingFromSpecificSnap(c *C) {
	s.RedirectClientToTestServer(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.Method, Equals, "POST")
		c.Check(r.URL.Path, Equals, "/2.0/interfaces")
		c.Check(DecodedRequestBody(c, r), DeepEquals, map[string]interface{}{
			"action": "disconnect",
			"plugs": []interface{}{
				map[string]interface{}{
					"snap": "",
					"plug": "",
				},
			},
			"slots": []interface{}{
				map[string]interface{}{
					"snap": "consumer",
					"slot": "",
				},
			},
		})
		fmt.Fprintln(w, `{"type":"sync", "result":{}}`)
	})
	rest, err := Parser().ParseArgs([]string{"disconnect", "consumer"})
	c.Assert(err, IsNil)
	c.Assert(rest, DeepEquals, []string{})
	c.Assert(s.Stdout(), Equals, "")
	c.Assert(s.Stderr(), Equals, "")
}
