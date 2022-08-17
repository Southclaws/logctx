// Package logctx provides a way to decorate structured log entries with
// metadata added to a `context.Context`.
package logctx

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var contextKey = struct{}{}

// Meta is a simple wrapper around a basic hash table that can be serialised
// into a log field for zap.
type Meta map[string]string

// MarshalLogObject implements zapcore.ObjectMarshaler
func (m Meta) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range m {
		enc.AddString(k, v)
	}
	return nil
}

// WithMeta creates a new context which contains a hash table of arbitrary
// metadata strings which can later be easily added to a structured log entry.
//
// For example, if you want to decorate a call tree with some data:
//
//    func DoBusinessLogic(ctx context.Context, userID string) error {
//        ctx = logctx.WithMeta(ctx, logctx.Meta{"user_id": userID})
//        GetResource(ctx, ...)
//    }
//
// You can wrap contexts with this helper as much as you want:
//
//    func DoBusinessLogic(ctx context.Context, userID string) error {
//        ctx = logctx.WithMeta(ctx, logctx.Meta{"user_id": userID})
//        GetResource(ctx, ...)
//    }
//
//    func GetResource(ctx context.Context, ...) {
//        ctx = logctx.WithMeta(ctx, logctx.Meta{"something_else": xyz})
//        CallAnotherThing(ctx, ...)
//    }
//
// Then, when you need to log it out, use `logctx.Zap`.
//
func WithMeta(ctx context.Context, data Meta) context.Context {
	// We don't need to stack metadata, just update/overwrite any existing keys.
	if existing, ok := ctx.Value(contextKey).(Meta); existing != nil && ok {
		for k, v := range data {
			existing[k] = v
		}

		return context.WithValue(ctx, contextKey, existing)
	}

	return context.WithValue(ctx, contextKey, data)
}

// Zap will wrap your Zap log fields with any available metadata from the given
// context. Any context returned from calls to `WithMeta` will work in this
// function and provide a "context" field to the log entry. If the given context
// was not decorated with `WithMeta` then this function does nothing and just
// passes your fields unmodified.
//
// It's best used directly in a zap log call, with the spread operator:
//
//    func (s *service) DoBusinessLogic(ctx context.Context, userID string) error {
//        s.l.Info("i am doing the thing", logctx.Zap(
//            zap.String("event_specific", "information"),
//        )...)
//    }
//
// In this example, assuming a function higher up in the call chain used the
// `WithMeta` to add a `user_id`, the log entry for this will be:
//
//     {
//         "level": "info",
//         "msg": "i am doing the thing",
//         "context": {
//             "user_id": "the_user_id"
//         }
//     }
//
func Zap(ctx context.Context, fields ...zapcore.Field) []zapcore.Field {
	value := ctx.Value(contextKey)
	if value == nil {
		return fields
	}

	casted, ok := value.(Meta)
	if !ok {
		return fields
	}

	return append(fields, zap.Object("context", casted))
}
