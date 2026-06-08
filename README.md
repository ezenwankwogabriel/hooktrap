# Hooktrap

A lightweight CLI tool for capturing, inspecting, and replaying webhook requests locally.

Built for developers who are tired of missing webhook payloads during local development.

## Installation

```bash
go install github.com/ezenwankwogabriel/hooktrap@latest
```

## Usage

Start catching webhooks:
```bash
hooktrap
hooktrap --port 9000
```

List captured requests:
```bash
hooktrap list
```

Replay a request:
```bash
hooktrap replay 2
hooktrap replay 2 --target http://localhost:9000
```

## Roadmap
- [ ] Public URL via tunnel integration
- [ ] Terminal UI with request inspector
- [ ] Webhook signature verification (Stripe, GitHub)
- [ ] Export requests as curl commands

## License
MIT