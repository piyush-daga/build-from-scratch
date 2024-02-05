package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		counter := 0

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "generate_ok"
		// One soln - using uuids, not the best approach in high throughput situations
		// body["id"] = fmt.Sprintf("%s-%s", msg.Dest, uuid.New().String())

		// Only uuid also works as a soln
		// body["id"] = uuid.New().String()

		// Try to fail by using just a local counter?

		counter++
		// One approach
		// body["id"] = fmt.Sprintf("%s-%s-%d", msg.Dest, time.Now().UTC(), counter)
		// But tha above apporach has a problem of not having ordering of messages by time, so it makes sense to put the
		// time component first
		body["id"] = fmt.Sprintf("%s-%s-%d", time.Now().UTC(), msg.Dest, counter)

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
