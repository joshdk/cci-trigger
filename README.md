[![License](https://img.shields.io/github/license/joshdk/cci-trigger.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/joshdk/cci-trigger?status.svg)](https://godoc.org/github.com/joshdk/cci-trigger)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshdk/cci-trigger)](https://goreportcard.com/report/github.com/joshdk/cci-trigger)
[![CircleCI](https://circleci.com/gh/joshdk/cci-trigger.svg?&style=shield)](https://circleci.com/gh/joshdk/cci-trigger/tree/master)

# CircleCI Trigger

⚙️ Trigger CircleCI builds programmatically

## Installing

You can fetch this tool by running the following

```bash
go get -u github.com/joshdk/cci-trigger
```

## Usage

### Setup

You must first obtain a API key by [going to your account's API page](https://circleci.com/account/api), and creating a new token. Export this API token as `CIRCLE_TOKEN` into your working environment.

```bash
export CIRCLE_TOKEN='cf1...d7c'
```

If you have an enterprise instance of CircleCI, you can export `CIRCLE_HOST` with the hostname of your internal instance.

```bash
export CIRCLE_HOST='circleci.example.com'
```

### Build head of default branch

Starts a build on the HEAD of the default branch. This branch is _typically_ master, and can usually be customized in your VCS platform.

```
$ cci-trigger username/project
https://circleci.com/gh/username/project/123
```

### Build tag

Starts a build on the given tag.

```
$ cci-trigger username/project --tag <TAG>
https://circleci.com/gh/username/project/123
```

### Build ref

Starts a build on the given VCS ref.

```
$ cci-trigger username/project --ref <REF>
https://circleci.com/gh/username/project/123
```

### Build head of branch

Starts a build on the HEAD of the given branch.

```
$ cci-trigger username/project --branch <BRANCH>
https://circleci.com/gh/username/project/123
```

### Build branch at ref

Starts a build on the given branch at the given ref.

```
$ cci-trigger username/project --branch <BRANCH> --ref <REF>
https://circleci.com/gh/username/project/123
```

### Rebuild build number

Restarts a build on the given build number.

```
$ cci-trigger username/project --build <BUILD>
https://circleci.com/gh/username/project/123
```

### Rebuild build number with SSH

Restarts a build on the given build number, and enables SSH.

```
$ cci-trigger username/project --build <BUILD> --ssh
https://circleci.com/gh/username/project/123
```

## Issues

If you find a bug in `cci-trigger` or need additional features, please feel free to [open an issue](https://github.com/joshdk/cci-trigger/issues/new) or [submit a pull request](https://github.com/joshdk/cci-trigger/pulls).

## License

This library is distributed under the [MIT License](https://opensource.org/licenses/MIT), see LICENSE.txt for more information.
