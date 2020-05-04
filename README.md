# PlanBot

PlanBot is a friendly helper for the MyHomeworkSpace development team! You may have seen it pop up on issues and unassign people who haven't been working on them. It also nags you on slack if you have an issue assigned to you that you haven't made any progress on. That's all it does for now, but hopefully it will do more in the future.

The bot is run daily on a cronjob.

## Development Setup
Copy the [`config.sample.toml`](./config.sample.toml) file and name the copy `config.toml`. Here is a breif description of each of the configuration options.

### `github`
| Field | Key | Description |
|---|---|---|
| Organazation | `organization` | The GitHub organization that PlanBot should operate in |
| Repository | `repo` | The GitHub repository that PlanBot should operate in |
| Private Key | `privatekey` | Path to the Private Key file that GitHub provides for your app. You can generate a private key in the "General" tab of your app settings on GitHub. |
| App ID | `appId` | The App ID that GitHub generates for your app. You can find it under the "About" section of your app's settings page (under the "General" tab) |
| Installation ID | `installationId` | The Installation ID that GitHub generates for the specific installation of the app to your organization. You can find this in the URL of the app install page, which can be reached by clicking the gear in the "Install App" tab of your app settings. In the URL, the Installation ID is found at the end, for example, in https://github.com/apps/mhs-planbot/installations/3423413, `3423413` is the installation ID. |

### `unassign`
| Field | Key | Description |
|---|---|---|
| Days until warning | `daysUntilWarning` | The number of days until the slack warning that your issue will soon be unassigned |
| Days until unassign | `daysUntilUnassign` | The number of days until the issue is unassigned due to inactivity |

### `slack`
| Field | Key | Description |
|---|---|---|
| Bot token | `token` | The bot OAuth token from Slack. It should begin with `xoxb-`. |
| Error log channel | `C260FU7HQ` | The slack channel to post errors in. |

Note that the Slack bot token requires the following OAuth permissions:
- `chat:write`
- `chat:write.public`
- `im:write`

### `users`
This should be a map of GitHub usernames to Slack User IDs so that the bot can associate a user's slack account with their GitHub account.

For example:

```toml
[users]
willbarkoff = "U0EGTKEFM"
thatoddmailbox = "U0EGDRXEK"
amazansky = "U0G281Q8L"
```