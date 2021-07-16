# GONEWS
A golang NNTP (Network News Transfer Protocol) server library and implementation

## About
This library provides the means for creating NNTP servers with customizable storage, message filtering, and authentication backends.
A useable reference implementation is also provided under cmd/gonews.go

## Features
TBD

## TODO
- RFC3977 Compliance
- Look at the possibilty of using a reactor pattern with a worker queue instead of a goroutine per connection
- Find better way to support CAPABILITIES
- Ability to support more fine grained authentication (restrict command use, group visibility, etc...)
- Commands
    - DATE
    - HDR
    - HELP
    - IHAVE
    - LIST
    - LISTGROUP
    - MODE
    - NEWGROUPS
    - NEWNEWS
    - OVER

## License
Both the nntp library and server are provided under the MIT license