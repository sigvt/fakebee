# Fake Bee

FakeBee is a CLI tool that generates fake Youtube Live events, and can produce them either to standard output or to a Kafka broker.

## Usage

```bash
make
./bin/fb start -b kafka --broker <broker addr>:9092
```

If you need more info about commands, use the `-h` flag.

## Configuration

Add the file `fb_config.json` to your `HOME` directory, or to the current working directory.

You can find an example in [fb_config.json](fb_config.json)

### Docker

```bash
docker build --tag fakebee .
docker run fakebee
```

## Worker hierarchy

Each `EventWorker` produces one kind of **events** to one specific **Kafka topic**.

```mermaid
erDiagram
  Queen ||--|{ EventWorker: manages

  EventWorker {
   string Topic
   int Interval
   int BacklogSize
   string OriginChannelId
   string OriginVideoId
  }
```

An `origin` is simply an entity representing a fake source of live chat events (a Youtube Live **stream** from a specific **channel**).

The `topic` can be one of:

- `chats`
- `superchats`
- `superstickers`
- `memberships`
- `milestones`
- `banactions`
- `deleteactions`

`fakebee` will create a corresponding Kafka `topic` if it does not exist in the cluster.
