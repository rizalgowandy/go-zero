package redis

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	red "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	ztrace "github.com/zeromicro/go-zero/core/trace"
)

func TestHookProcessCase1(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})

	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx, err := durationHook.BeforeProcess(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, durationHook.AfterProcess(ctx, red.NewCmd(context.Background())))
	assert.False(t, strings.Contains(buf.String(), "slow"))
	assert.Equal(t, "redis", ctx.Value(spanKey).(interface{ Name() string }).Name())
}

func TestHookProcessCase2(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})

	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx, err := durationHook.BeforeProcess(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "redis", ctx.Value(spanKey).(interface{ Name() string }).Name())

	time.Sleep(slowThreshold.Load() + time.Millisecond)

	assert.Nil(t, durationHook.AfterProcess(ctx, red.NewCmd(context.Background(), "foo", "bar")))
	assert.True(t, strings.Contains(buf.String(), "slow"))
	assert.True(t, strings.Contains(buf.String(), "trace"))
	assert.True(t, strings.Contains(buf.String(), "span"))
}

func TestHookProcessCase3(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	assert.Nil(t, durationHook.AfterProcess(context.Background(), red.NewCmd(context.Background())))
	assert.True(t, buf.Len() == 0)
}

func TestHookProcessCase4(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx := context.WithValue(context.Background(), startTimeKey, "foo")
	assert.Nil(t, durationHook.AfterProcess(ctx, red.NewCmd(context.Background())))
	assert.True(t, buf.Len() == 0)
}

func TestHookProcessPipelineCase1(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx, err := durationHook.BeforeProcessPipeline(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "redis", ctx.Value(spanKey).(interface{ Name() string }).Name())

	assert.Nil(t, durationHook.AfterProcessPipeline(ctx, []red.Cmder{
		red.NewCmd(context.Background()),
	}))
	assert.False(t, strings.Contains(buf.String(), "slow"))
}

func TestHookProcessPipelineCase2(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})

	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx, err := durationHook.BeforeProcessPipeline(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "redis", ctx.Value(spanKey).(interface{ Name() string }).Name())

	time.Sleep(slowThreshold.Load() + time.Millisecond)

	assert.Nil(t, durationHook.AfterProcessPipeline(ctx, []red.Cmder{
		red.NewCmd(context.Background(), "foo", "bar"),
	}))
	assert.True(t, strings.Contains(buf.String(), "slow"))
	assert.True(t, strings.Contains(buf.String(), "trace"))
	assert.True(t, strings.Contains(buf.String(), "span"))
}

func TestHookProcessPipelineCase3(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	assert.Nil(t, durationHook.AfterProcessPipeline(context.Background(), []red.Cmder{
		red.NewCmd(context.Background()),
	}))
	assert.True(t, buf.Len() == 0)
}

func TestHookProcessPipelineCase4(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx := context.WithValue(context.Background(), startTimeKey, "foo")
	assert.Nil(t, durationHook.AfterProcessPipeline(ctx, []red.Cmder{
		red.NewCmd(context.Background()),
	}))
	assert.True(t, buf.Len() == 0)
}

func TestHookProcessPipelineCase5(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx := context.WithValue(context.Background(), startTimeKey, "foo")
	assert.Nil(t, durationHook.AfterProcessPipeline(ctx, nil))
	assert.True(t, buf.Len() == 0)
}
