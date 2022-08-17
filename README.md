# logctx

Package logctx provides a way to decorate structured log entries with metadata
added to a `context.Context`.

---

Currently, this works with Zap only. But it would be trivial to support other
structured logging libraries.

WithMeta creates a new context which contains a hash table of arbitrary
metadata strings which can later be easily added to a structured log entry.

For example, if you want to decorate a call tree with some data:

```go
func DoBusinessLogic(ctx context.Context, userID string) error {
    ctx = logctx.WithMeta(ctx, logctx.Meta{"user_id": userID})
    GetResource(ctx, ...)
}
```

You can wrap contexts with this helper as much as you want:

```go
func DoBusinessLogic(ctx context.Context, userID string) error {
    ctx = logctx.WithMeta(ctx, logctx.Meta{"user_id": userID})
    GetResource(ctx, ...)
}

func GetResource(ctx context.Context, ...) {
    ctx = logctx.WithMeta(ctx, logctx.Meta{"something_else": xyz})
    CallAnotherThing(ctx, ...)
}
```

Then, when you need to log it out, use `logctx.Zap`.

Zap will wrap your Zap log fields with any available metadata from the given
context. Any context returned from calls to `WithMeta` will work in this
function and provide a "context" field to the log entry. If the given context
was not decorated with `WithMeta` then this function does nothing and just
passes your fields unmodified.

It's best used directly in a zap log call, with the spread operator:

```go
func (s *service) DoBusinessLogic(ctx context.Context, userID string) error {
    s.l.Info("i am doing the thing", logctx.Zap(
        zap.String("event_specific", "information"),
    )...)
}
```

In this example, assuming a function higher up in the call chain used the
`WithMeta` to add a `user_id`, the log entry for this will be:

```json
{
    "level": "info",
    "msg": "i am doing the thing",
    "context": {
        "user_id": "the_user_id"
    }
}
```
