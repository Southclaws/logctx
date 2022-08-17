package logctx_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Southclaws/logctx"
)

func testLogger() (*zap.Logger, *bytes.Buffer) {
	buf := bytes.NewBuffer(nil)
	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(buf), zap.LevelEnablerFunc(func(level zapcore.Level) bool { return true })))
	return logger, buf
}

func TestContext(t *testing.T) {
	a := assert.New(t)
	logger, buf := testLogger()

	// make a new context
	root := context.Background()

	// add some metadata to this context, a user ID
	ctx := logctx.WithMeta(root, map[string]string{"user_id": "southclaws"})

	logger.Info("test context", logctx.Zap(ctx)...)

	a.Contains(buf.String(), `"context":{"user_id":"southclaws"}`)
}

func TestContextNested(t *testing.T) {
	a := assert.New(t)
	logger, buf := testLogger()

	// make a new context
	root := context.Background()

	// add some metadata to this context, a user ID
	ctx1 := logctx.WithMeta(root, map[string]string{"user_id": "southclaws"})
	ctx2 := logctx.WithMeta(ctx1, map[string]string{"deal_id": "xyz"})
	ctx3 := logctx.WithMeta(ctx2, map[string]string{"commitment_id": "123"})

	logger.Info("test context", logctx.Zap(ctx3)...)

	a.Contains(buf.String(), `"user_id":"southclaws"`)
	a.Contains(buf.String(), `"deal_id":"xyz"`)
	a.Contains(buf.String(), `"commitment_id":"123"`)
}

func TestContextNestedOverwrite(t *testing.T) {
	a := assert.New(t)
	logger, buf := testLogger()

	// make a new context
	root := context.Background()

	// add some metadata to this context, a user ID
	ctx1 := logctx.WithMeta(root, map[string]string{"user_id": "southclaws"})
	ctx2 := logctx.WithMeta(ctx1, map[string]string{"deal_id": "xyz"})
	ctx3 := logctx.WithMeta(ctx2, map[string]string{"deal_id": "overwrite context metadata"})

	logger.Info("test context", logctx.Zap(ctx3)...)

	a.Contains(buf.String(), `"user_id":"southclaws"`)
	a.Contains(buf.String(), `"deal_id":"overwrite context metadata"`)
}

func TestContextEmpty(t *testing.T) {
	a := assert.New(t)

	// make a new context
	root := context.Background()

	// test logging out the metadata
	logger, buf := testLogger()

	logger.Info("test context", logctx.Zap(root, zap.String("message", "hello"))...)

	// no context key means no log field
	a.NotContains(buf.String(), `"context"`)
	a.Contains(buf.String(), "message")
	a.Contains(buf.String(), "hello")
}
