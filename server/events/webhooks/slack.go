// Copyright 2017 HootSuite Media Inc.
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Modified hereafter by contributors to runatlantis/atlantis.

package webhooks

import (
	"regexp"

	"fmt"

	"github.com/runatlantis/atlantis/server/logging"
)

// SlackWebhook sends webhooks to Slack.
type SlackWebhook struct {
	Client         SlackClient
	WorkspaceRegex *regexp.Regexp
	BranchRegex    *regexp.Regexp
	Channel        string
}

func NewSlack(wr *regexp.Regexp, br *regexp.Regexp, channel string, client SlackClient) (*SlackWebhook, error) {
	if err := client.AuthTest(); err != nil {
		return nil, fmt.Errorf("testing slack authentication: %s. Verify your slack-token is valid", err)
	}

	return &SlackWebhook{
		Client:         client,
		WorkspaceRegex: wr,
		BranchRegex:    br,
		Channel:        channel,
	}, nil
}

// Send sends the webhook to Slack if workspace and branch matches their respective regex.
func (s *SlackWebhook) Send(log logging.SimpleLogging, applyResult ApplyResult) error {
	if !s.WorkspaceRegex.MatchString(applyResult.Workspace) || !s.BranchRegex.MatchString(applyResult.Pull.BaseBranch) {
		return nil
	}
	return s.Client.PostMessage(s.Channel, applyResult)
}
