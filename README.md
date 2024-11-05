# Implementing websocket communication client-server with Go and React

// as a test first I'm implementing a simple index.html frontend with JS to establish the connection and just keep pushing numbers to the page

Requirement: react frontend connects to go backend and they establish a websocket connection; the server sends a random number in a 1-10 range every second; the client can tell the server to pause/resume these updates; if disconnection occurs, by any means, there should be reconnection attempts

- note - not sure how I can test random disconnects, maybe I find a way to simulate it? maybe through browser devtools throttle for the frontend, not sure for the backend though; since both programs run locally, it's not possible to disconnect from wifi for a while to simulate network problems

- using gorilla/websocket here; it was archived a while ago and then reopened by another maintainer; last commit is 3 months ago at the time of writing this (Nov 2024), so it looks like it's still going
- react will publish the incoming numbers to the screen and maybe will later build a real-time line chart with these numbers, simulating a real-time analytics scenario
- weird behaviour comes when using goroutines to handle this continuous writing on a separate thread; I also had situations when adding the time.Sleep(time.Second) delay would freeze the entire server, even though there was no client connected, just 1 trying to connect; there is no immediate reason for this behaviour I could find
