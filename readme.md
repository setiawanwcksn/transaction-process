This service accepts bank statement files (CSV), parses each row, and builds an in-memory report.
Any rows marked as PENDING is pushed to a worker queue for background processing.

The whole system runs fully in memory, with no external dependencies.

```
[CSV Upload]
      │
      ▼
  HTTP /upload
      │
      ▼
   Parser ──► Storage (SUCCESS/FAILED)
      │
      └────► Bus (PENDING)
                 │
                 ▼
            Worker pool
                 │ (retry + backoff)
                 ▼
           Storage (resolved SUCCESS)
```

Components
- API — receives upload requests and starts parsing in a goroutine
- Parser — reads CSV, stores successful rows, sends pending to bus
- Storage — holds all results in memory (balance, lines, issues)
- Bus — lightweight channel for storing events
- Worker — retries PENDING, logs/finalizes FAILED

🎯 Key Behaviors
Parsing
- SUCCESS → stored immediately, affects balance
- FAILED → logged as issue
- PENDING → sent to worker for retry

Worker
-Introduces delay with backoff on pending items
- Each item processed only once (idempotent tracking)
- Does not block API

Storage
- Balance and issue tracking updated as items settle
- All data is isolated per upload ID

Pros:
- Easy to run locally — no database needed
- Event flow makes slow work asynchronous
- Worker does not block client requests
- Simple API with predictable behavior

Cons:
- Data disappears on restart (memory-only)
- No horizontal scaling without shared state
- Worker retry logic is basic (not durable queue)
- File parsing is single-threaded

What could be improved
- Move storage to Redis or SQL
- Use a durable queue (Kafka / RabbitMQ)
- Add metrics & tracing
- Add auth / rate limiting

How to Run
``` 
make build
make run 
```

Server will start on:
```
http://localhost:8080
```

Example Requests:

upload file:
```
curl -X POST      -F "file=@big_test.csv"      http://localhost:8080/statements
```

check progress:
```
curl "http://localhost:8080/progress?upload_id=<id>"
```

check balance:
```
curl "http://localhost:8080/balance?upload_id=<id>"
```

view issues:
```
curl "http://localhost:8080/issues?upload_id=<id>&limit=<...>&offset=<...>"
```
