# gh-pull-requests-slack-reminder

A CLI application that notifies to Slack with pull requests that have a specific label, like `in-review`.


## Installation

```sh
go install github.com/nekonenene/gh-pull-requests-slack-reminder@v1
```

## Usage

First, you need to get GitHub API Token to control your repository, please see [here](https://docs.github.com/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).  
The token needs the read-access permission of pull requests and issues.

### Example

```sh
gh-pull-requests-slack-reminder --token 123456789abcd123456789abcd --owner nekonenene --repo my-repository-name --label-name "in-review" --webhook-url https://hooks.slack.com/services/XXXXXXXX --avoid-weekend
```

### Parameters

| Parameter | Description | Required? |
|:---:|:---:|:---:|
|-token| GitHub API Token | YES |
|-owner| Owner name of the repository (e.g. octocat) | YES |
|-repo| Repository name (e.g. hello-world) | YES |
|-label-name| Label name related to target pull requests (e.g. in-review) | YES |
|-webhook-url| URL of Slack Incoming Webhook | YES |
|-avoid-weekend| If set, don't send notifications on weekends |  |
|-avoid-jp-holidays| If set, don't send notifications on Japanese holidays |  |
|-enterprise-url| URL of GitHub Enterprise (e.g. https://github.your.domain ) |  |
|-dry-run| If set, don't send notifications to Slack |  |

And this command shows all parameters:

```sh
gh-pull-requests-slack-reminder --help
```

### Run with GitHub Actions

An example file of [GitHub Actions](https://docs.github.com/actions/using-workflows/workflow-syntax-for-github-actions):

```yml
name: Pull Request Review Reminder
on:
  schedule:
    # Exec on 0:30 and 6:30 UTC every day
    - cron: "30 0,6 * * *"

jobs:
  gh-pull-requests-slack-reminder:
    name: gh-pull-requests-slack-reminder
    runs-on: ubuntu-latest
    env:
      TZ: Asia/Tokyo
    steps:
      - uses: actions/setup-go@v4
        with:
          go-version: ^1.21.4
      - name: Install gh-pull-requests-slack-reminder
        run: go install github.com/nekonenene/gh-pull-requests-slack-reminder@v1
      - name: Run gh-pull-requests-slack-reminder
        run: >
          gh-pull-requests-slack-reminder
            --token ${{ secrets.GITHUB_TOKEN }}
            --owner ${{ github.repository_owner }}
            --repo ${{ github.event.repository.name }}
            --label-name "in-review"
            --webhook-url "https://hooks.slack.com/services/XXXXXXXX"
            --avoid-weekend
```


## License

[MIT](https://choosealicense.com/licenses/mit/)
