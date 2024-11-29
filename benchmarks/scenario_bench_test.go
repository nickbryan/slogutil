package benchmarks

import (
	"context"
	"log"
	"log/slog"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func BenchmarkDisabledWithoutFields(b *testing.B) {
	b.Logf("Logging at a disabled level without any structured context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if m := logger.Check(zap.InfoLevel, getMessage(0)); m != nil {
					m.Write()
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("slog", func(b *testing.B) {
		logger := newDisabledSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogmem", func(b *testing.B) {
		logger := newDisabledSlogUtilInMem()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogmem.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilInMem()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogctx", func(b *testing.B) {
		logger := newDisabledSlogUtilCtx()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogctx.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilCtx()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogutiljsonlogger", func(b *testing.B) {
		logger := newSlogUtilJSONLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogutiljsonlogger.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilJSONLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
}

func BenchmarkDisabledAccumulatedContext(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if m := logger.Check(zap.InfoLevel, getMessage(0)); m != nil {
					m.Write()
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("slog", func(b *testing.B) {
		logger := newDisabledSlog(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlog(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogmem", func(b *testing.B) {
		logger := newDisabledSlogUtilInMem(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogmem.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilInMem(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogctx", func(b *testing.B) {
		logger := newDisabledSlogUtilCtx(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogctx.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilCtx(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogutiljsonlogger", func(b *testing.B) {
		logger := newDisabledSlogUtilJSONLogger(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogutiljsonlogger.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilJSONLogger(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
}

func BenchmarkDisabledAddingFields(b *testing.B) {
	b.Logf("Logging at a disabled level, adding context at each log site.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeFields()...)
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if m := logger.Check(zap.InfoLevel, getMessage(0)); m != nil {
					m.Write(fakeFields()...)
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("slog", func(b *testing.B) {
		logger := newDisabledSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
	})
	b.Run("slogmem", func(b *testing.B) {
		logger := newDisabledSlogUtilInMem()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slogmem.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilInMem()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
	})
	b.Run("slogctx", func(b *testing.B) {
		logger := newDisabledSlogUtilCtx()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slogctx.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilCtx()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
	})
	b.Run("slogutiljsonlogger", func(b *testing.B) {
		logger := newDisabledSlogUtilJSONLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slogutiljsonlogger.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlogUtilJSONLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
	})
}

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
					ce.Write()
				}
			}
		})
	})
	b.Run("Zap.CheckSampled", func(b *testing.B) {
		logger := newSampledLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if ce := logger.Check(zap.InfoLevel, getMessage(i)); ce != nil {
					ce.Write()
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("stdlib.Println", func(b *testing.B) {
		logger := log.New(&zaptest.Discarder{}, "", log.LstdFlags)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Println(getMessage(0))
			}
		})
	})
	b.Run("stdlib.Printf", func(b *testing.B) {
		logger := log.New(&zaptest.Discarder{}, "", log.LstdFlags)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Printf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("slog", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogmem", func(b *testing.B) {
		logger := newSlogUtilInMem()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogmem.LogAttrs", func(b *testing.B) {
		logger := newSlogUtilInMem()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogctx", func(b *testing.B) {
		logger := newSlogUtilCtx()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogctx.LogAttrs", func(b *testing.B) {
		logger := newSlogUtilCtx()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogutiljsonlogger", func(b *testing.B) {
		logger := newSlogUtilJSONLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogutiljsonlogger.LogAttrs", func(b *testing.B) {
		logger := newSlogUtilJSONLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
}

func BenchmarkAccumulatedContext(b *testing.B) {
	b.Logf("Logging with some accumulated context.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
					ce.Write()
				}
			}
		})
	})
	b.Run("Zap.CheckSampled", func(b *testing.B) {
		logger := newSampledLogger(zap.DebugLevel).With(fakeFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if ce := logger.Check(zap.InfoLevel, getMessage(i)); ce != nil {
					ce.Write()
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).With(fakeFields()...).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("slog", func(b *testing.B) {
		logger := newSlog(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newSlog(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogmem", func(b *testing.B) {
		logger := newSlogUtilInMem(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogmem.LogAttrs", func(b *testing.B) {
		logger := newSlogUtilInMem(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
	b.Run("slogutiljsonlogger", func(b *testing.B) {
		logger := newSlogUtilJSONLogger(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slogutiljsonlogger.LogAttrs", func(b *testing.B) {
		logger := newSlogUtilJSONLogger(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
}

func BenchmarkAddingFields(b *testing.B) {
	b.Logf("Logging with additional context at each log site.")
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeFields()...)
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
					ce.Write(fakeFields()...)
				}
			}
		})
	})
	b.Run("Zap.CheckSampled", func(b *testing.B) {
		logger := newSampledLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if ce := logger.Check(zap.InfoLevel, getMessage(i)); ce != nil {
					ce.Write(fakeFields()...)
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("slog", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
	})
	b.Run("slogmem", func(b *testing.B) {
		logger := newSlogUtilInMem()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slogmem.LogAttrs", func(b *testing.B) {
		logger := newSlogUtilInMem()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
	})
	b.Run("slogctx", func(b *testing.B) {
		logger := newSlogUtilCtx()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slogctx.LogAttrs", func(b *testing.B) {
		logger := newSlogUtilJSONLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
	})
	b.Run("slogutiljsonlogger", func(b *testing.B) {
		logger := newSlogUtilCtx()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slogutiljsonlogger.LogAttrs", func(b *testing.B) {
		logger := newSlogUtilJSONLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
	})
}
