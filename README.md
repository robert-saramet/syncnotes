# syncnotes
- Note syncing git utility for lazy people
- Only Github is currently supported
- Needs SSH on private repositories
- Pull requests welcome!

## Installation
- Binary:
  - download from [releases](https://github.com/robert-saramet/syncnotes/releases/)
  - save it somewhere in your `PATH`
- **With Go**:
  - `$ go install github.com/robert-saramet/syncnotes@latest`
- From source:
  - `$ git clone https://github.com/robert-saramet/syncnotes`
  - `$ go install syncnotes`

## Commands
- `syncnotes config` *set repo and directory*
- `syncnotes push [dir]` *upload specified directory*
- `syncnotes` *sync default notes directory*

## To-Do
- [ ] Support classic auth for private repos
- [ ] Handle errors properly
- [ ] Handle git conflicts
- [ ] Create background service
- [ ] Add GitHub webhook support
- [ ] Remove configuration file
