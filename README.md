# Implementing websocket communication client-server with Go and React

// as a test first I'm implementing a simple index.html frontend with JS to establish the connection and just keep pushing numbers to the page

Requirement: react frontend connects to go backend and they establish a websocket connection; the server sends a random number in a 1-10 range every second; the client can tell the server to pause/resume these updates; if disconnection occurs, by any means, there should be reconnection attempts

- note - not sure how I can test random disconnects, maybe I find a way to simulate it? maybe through browser devtools throttle for the frontend, not sure for the backend though; since both programs run locally, it's not possible to disconnect from wifi for a while to simulate network problems

- using gorilla/websocket here; it was archived a while ago and then reopened by another maintainer; last commit is 3 months ago at the time of writing this (Nov 2024), so it looks like it's still going
- react will publish the incoming numbers to the screen and maybe will later build a real-time line chart with these numbers, simulating a real-time analytics scenario
- weird behaviour comes when using goroutines to handle this continuous writing on a separate thread; I also had situations when adding the time.Sleep(time.Second) delay would freeze the entire server, even though there was no client connected, just 1 trying to connect; there is no immediate reason for this behaviour I could find
  - I could not reproduce this behaviour after the fact, right now 2 goroutines handle reading and writing in parallel with a pauseChannel controlled by the reader based on messages from the client
- I would like more information from the gorilla/websocket documentation to help with achieving various goals; there are a few examples on their github, though learning from them comes from reverse-engineering the code since there are no explained intentions or best practices
- websockets turn out to be not that hard to implement once the library/package is understood; would need more testing in a production environment to cover accidental disconnections though; those should be initiated by the client; if we're talking JS/TS, most likely socket.onclose could reinitialize the connection
- I leave as a future TODO testing regarding efficiency, or at least understanding the inner workings of the gorilla/websocket package to be able to write better code
- there is also the _websocket.html_ statically served file, coming from code example from the Internet; it's an echo app (what the client sends gets repeated back)

A few notes on what I've learned

- read docs carefully regarding package limitations; in this case, it is stated that concurrent read/writes on same connection from multiple threads/goroutines is not supported; people online tend to synchronize them via mutexes; haven't tested, should work in theory; I would rather suggest a thread to read/write and spawning other goroutines for handling further processing if necessary
- careful with retring reads/writes; React dev server (at least) always executes things twice, so it would initiate 2 websockets and immediately kill one; by retrying to read from the dead connection, the Go server would panic; this first left the impression I could not read and write from the same connection at the same time, which would defeat the purpose of websockets
