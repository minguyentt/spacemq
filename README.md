# spacemq

## TODO

### internal
- finish the internal workflow for Dequeue Command
- get more in tune how you wanna serialize and deserialize the data

- implement queue acknowledgement
- queue completions
- dead-to-letter queues
- maybe retries?

#### Stretch goals
- refactoring
- Scheduler

### external/public
- build a simple taskqueue client for simulation/debugging
- intuitive custom error type handling

- implement healthchecker to periodically ping the taskqueue
- figure out how to implement heartbeater

- NOTE: run tests or check how the data is returned from the scripts

## Tools

github.com/json-iterator/go - 5x faster for encoding/decoding json compared to std lib

