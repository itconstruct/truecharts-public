package main

import (
    "fmt"
    "os"
    "time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"

    "github.com/truecharts/public/clustertool/cmd"
    "github.com/truecharts/public/clustertool/embed"
    "github.com/truecharts/public/clustertool/pkg/helper"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
    ctrllogzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var version = "dev"

func main() {
    // Configure zerolog
    zerolog.DurationFieldUnit = time.Second

    var zerologLevel zerolog.Level
    var zapLevel zapcore.Level

    // Switch-case for setting the global log level
    switch os.Getenv("DEBUG") {
    case "trace":
        zerologLevel = zerolog.TraceLevel
        zapLevel = zapcore.DebugLevel // zap does not have a Trace level, use Debug as equivalent
    case "debug":
        zerologLevel = zerolog.DebugLevel
        zapLevel = zapcore.DebugLevel
    case "warn":
        zerologLevel = zerolog.WarnLevel
        zapLevel = zapcore.WarnLevel
    case "error":
        zerologLevel = zerolog.ErrorLevel
        zapLevel = zapcore.ErrorLevel
    case "fatal":
        zerologLevel = zerolog.FatalLevel
        zapLevel = zapcore.FatalLevel
    case "panic":
        zerologLevel = zerolog.PanicLevel
        zapLevel = zapcore.PanicLevel
    default:
        zerologLevel = zerolog.InfoLevel
        zapLevel = zapcore.InfoLevel
    }

    // Set zerolog level
    zerolog.SetGlobalLevel(zerologLevel)
    log.Logger = log.Output(zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: time.RFC3339, // Keep this for the timestamp format
        NoColor:    false,        // Set to true if you prefer no color
        FormatLevel: func(i interface{}) string {
            return fmt.Sprintf("\033[1;34m[%s]\033[0m", i) // Blue color for level
        },
        FormatMessage: func(i interface{}) string {
            return fmt.Sprintf("\033[1;32m%s\033[0m", i) // Green color for message
        },
    })

    // Configure zap logger
    zapConfig := zap.NewProductionConfig()
    zapConfig.Level = zap.NewAtomicLevelAt(zapLevel)
    zapLogger, err := zapConfig.Build()
    if err != nil {
        panic(err)
    }
    defer zapLogger.Sync()

    // Set controller-runtime logger to use zap
    ctrlLogger := ctrllogzap.New(ctrllogzap.UseDevMode(true), ctrllogzap.Level(zapLevel))
    ctrllog.SetLogger(ctrlLogger)
    fmt.Printf("\n%s\n", helper.Logo)
    fmt.Printf("---\nClustertool Version: %s\n---\n", version)

    embed.AllToCache()

    helper.CheckSystemTime()
    helper.CheckReqDomains()

    err = cmd.Execute()
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to execute command")
    }
}
