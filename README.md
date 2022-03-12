# mque

Simple message queue application.

## Getting Started

```bash
# Go to the root project directory and run
go run .

# Open another terminal screen and run
# this will run the consumer 
cd ./examples/singlecon
go run .

# Open another terminal screen and run
# this will run the publisher
# you can write some texts to terminal and hit enter
# you will see the messages in consumer terminal
cd ./examples/publisher
go run .
```