# Benchmarks
The following is based on [uber-go/zap benchmark suite](https://github.com/uber-go/zap/tree/master/benchmarks)
and compares the usage of slogutil with Zap and log/slog. 

**NOTE:** these benchmarks are for comparative analysis only.

## Running
You can run the benchmarks by navigating into this benchmark directory and then running the following:
```text
    $ go test -bench=. -benchmem -v
```

## Example Output
Expect different results to the following depending on the environment the benchmarks are run in. These results
have been reformatted and trimmed for clarity.

```text
go test -bench=. -benchmem -v                                                                                                                                                     1 ↵
goos: darwin
goarch: arm64
pkg: github.com/nickbryan/slogutil/benchmarks
cpu: Apple M2
BenchmarkDisabledWithoutFields
    scenario_bench_test.go:14: Logging at a disabled level without any structured context.
BenchmarkDisabledWithoutFields/Zap-8    	                    1000000000	         0.8416 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledWithoutFields/Zap.Check-8         	            1000000000	         0.6763 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledWithoutFields/Zap.Sugar-8         	             197632594	         6.204  ns/op	      16 B/op	       1 allocs/op
BenchmarkDisabledWithoutFields/Zap.SugarFormatting-8         	  24661567	        44.83   ns/op	     136 B/op	       6 allocs/op
BenchmarkDisabledWithoutFields/slog-8                        	1000000000	         0.8537 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledWithoutFields/slog.LogAttrs-8               	1000000000	         0.8473 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledWithoutFields/slogmem-8                     	1000000000	         0.8405 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledWithoutFields/slogmem.LogAttrs-8            	1000000000	         0.8366 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledWithoutFields/slogctx-8                     	1000000000	         1.010  ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledWithoutFields/slogctx.LogAttrs-8            	1000000000	         1.008  ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledWithoutFields/slogutiljsonlogger-8          	   2371167	       507.6    ns/op	     770 B/op	      10 allocs/op
BenchmarkDisabledWithoutFields/slogutiljsonlogger.LogAttrs-8 	1000000000	         1.013  ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext
    scenario_bench_test.go:128: Logging at a disabled level with some accumulated context.
BenchmarkDisabledAccumulatedContext/Zap-8                    	    1000000000	         0.9030 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/Zap.Check-8              	    1000000000	         0.6796 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/Zap.Sugar-8              	     175562364	         6.747  ns/op	      16 B/op	       1 allocs/op
BenchmarkDisabledAccumulatedContext/Zap.SugarFormatting-8    	      24053679	        47.56   ns/op	     136 B/op	       6 allocs/op
BenchmarkDisabledAccumulatedContext/slog-8                   	    1000000000	         0.8817 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/slog.LogAttrs-8          	    1000000000	         0.8897 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/slogmem-8                	    1000000000	         0.8901 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/slogmem.LogAttrs-8       	    1000000000	         0.8877 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/slogctx-8                	    1000000000	         1.047  ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/slogctx.LogAttrs-8       	    1000000000	         1.049  ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/slogutiljsonlogger-8     	    1000000000	         1.084  ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAccumulatedContext/slogutiljsonlogger.LogAttrs-8   1000000000	         1.053  ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAddingFields
    scenario_bench_test.go:242: Logging at a disabled level, adding context at each log site.
BenchmarkDisabledAddingFields/Zap-8                                    6552414	       180.4    ns/op	     800 B/op	       5 allocs/op
BenchmarkDisabledAddingFields/Zap.Check-8                           1000000000	         0.6826 ns/op	       0 B/op	       0 allocs/op
BenchmarkDisabledAddingFields/Zap.Sugar-8                             22369417	        48.93   ns/op	     136 B/op	       6 allocs/op
BenchmarkDisabledAddingFields/slog-8                                  22811067	        49.54   ns/op	     136 B/op	       6 allocs/op
BenchmarkDisabledAddingFields/slog.LogAttrs-8                          9199884	       132.0    ns/op	     512 B/op	       5 allocs/op
BenchmarkDisabledAddingFields/slogmem-8                               23727606	        49.00   ns/op	     136 B/op	       6 allocs/op
BenchmarkDisabledAddingFields/slogmem.LogAttrs-8                       9037459	       132.2    ns/op	     512 B/op	       5 allocs/op
BenchmarkDisabledAddingFields/slogctx-8                               23515214	        49.44   ns/op	     136 B/op	       6 allocs/op
BenchmarkDisabledAddingFields/slogctx.LogAttrs-8                       9010521	       127.4    ns/op	     512 B/op	       5 allocs/op
BenchmarkDisabledAddingFields/slogutiljsonlogger-8                    23266390	        49.28   ns/op	     136 B/op	       6 allocs/op
BenchmarkDisabledAddingFields/slogutiljsonlogger.LogAttrs-8            9037641	       130.9    ns/op	     512 B/op	       5 allocs/op
BenchmarkWithoutFields
    scenario_bench_test.go:347: Logging without any structured context.
BenchmarkWithoutFields/Zap-8                                              	23825655	        60.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithoutFields/Zap.Check-8                                        	24326110	        55.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithoutFields/Zap.CheckSampled-8                                 	34241444	        34.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithoutFields/Zap.Sugar-8                                        	16025140	        73.00 ns/op	      16 B/op	       1 allocs/op
BenchmarkWithoutFields/Zap.SugarFormatting-8                              	  727348	      1642.00 ns/op	    1919 B/op	      58 allocs/op
BenchmarkWithoutFields/stdlib.Println-8                                   	10619124	       113.3  ns/op	      16 B/op	       1 allocs/op
BenchmarkWithoutFields/stdlib.Printf-8                                    	  936378	      1299.00 ns/op	    1274 B/op	      57 allocs/op
BenchmarkWithoutFields/slog-8                                             	 7475978	       152.8  ns/op	       0 B/op	       0 allocs/op
BenchmarkWithoutFields/slog.LogAttrs-8                                    	 7485754	       159.8  ns/op	       0 B/op	       0 allocs/op
BenchmarkWithoutFields/slogmem-8                                          	 5292988	       267.1  ns/op	     507 B/op	       1 allocs/op
BenchmarkWithoutFields/slogmem.LogAttrs-8                                 	 5017994	       218.1  ns/op	     438 B/op	       1 allocs/op
BenchmarkWithoutFields/slogctx-8                                          	 6059479	       188.1  ns/op	      64 B/op	       1 allocs/op
BenchmarkWithoutFields/slogctx.LogAttrs-8                                 	 6265539	       188.7  ns/op	      64 B/op	       1 allocs/op
BenchmarkWithoutFields/slogutiljsonlogger-8                               	 2313895	       520.2  ns/op	     770 B/op	      10 allocs/op
BenchmarkWithoutFields/slogutiljsonlogger.LogAttrs-8                      	 2271818	       526.6  ns/op	     770 B/op	      10 allocs/op
BenchmarkAccumulatedContext
    scenario_bench_test.go:492: Logging with some accumulated context.
BenchmarkAccumulatedContext/Zap-8                                         	19919850	        65.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccumulatedContext/Zap.Check-8                                   	22762281	        61.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccumulatedContext/Zap.CheckSampled-8                            	33184693	        34.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccumulatedContext/Zap.Sugar-8                                   	16269604	        76.88 ns/op	      16 B/op	       1 allocs/op
BenchmarkAccumulatedContext/Zap.SugarFormatting-8                         	  704304	      1558.00 ns/op	    1922 B/op	      58 allocs/op
BenchmarkAccumulatedContext/slog-8                                        	 7304456	       164.1  ns/op	       0 B/op	       0 allocs/op
BenchmarkAccumulatedContext/slog.LogAttrs-8                               	 6567868	       164.7  ns/op	       0 B/op	       0 allocs/op
BenchmarkAccumulatedContext/slogmem-8                                     	 2031447	       596.8  ns/op	    1366 B/op	      13 allocs/op
BenchmarkAccumulatedContext/slogmem.LogAttrs-8                            	 1984200	       591.1  ns/op	    1375 B/op	      13 allocs/op
BenchmarkAccumulatedContext/slogutiljsonlogger-8                          	 2293304	       524.7  ns/op	     770 B/op	      10 allocs/op
BenchmarkAccumulatedContext/slogutiljsonlogger.LogAttrs-8                 	 2308488	       522.1  ns/op	     770 B/op	      10 allocs/op
BenchmarkAddingFields
    scenario_bench_test.go:601: Logging with additional context at each log site.
BenchmarkAddingFields/Zap-8                                               	 1834888	       653.1  ns/op	     803 B/op	       5 allocs/op
BenchmarkAddingFields/Zap.Check-8                                         	 1777706	       680.7  ns/op	     803 B/op	       5 allocs/op
BenchmarkAddingFields/Zap.CheckSampled-8                                  	11999574	        99.86 ns/op	      86 B/op	       0 allocs/op
BenchmarkAddingFields/Zap.Sugar-8                                         	 1295394	       906.0  ns/op	    1624 B/op	      10 allocs/op
BenchmarkAddingFields/slog-8                                              	  502731	      2188.00 ns/op	    3174 B/op	      41 allocs/op
BenchmarkAddingFields/slog.LogAttrs-8                                     	  496743	      2286.00 ns/op	    3553 B/op	      40 allocs/op
BenchmarkAddingFields/slogmem-8                                           	 1355288	       913.4  ns/op	    2526 B/op	      22 allocs/op
BenchmarkAddingFields/slogmem.LogAttrs-8                                  	 1329951	       939.7  ns/op	    2909 B/op	      21 allocs/op
BenchmarkAddingFields/slogctx-8                                           	  379437	      3010.00 ns/op	    5222 B/op	      57 allocs/op
BenchmarkAddingFields/slogctx.LogAttrs-8                                  	  326156	      3542.00 ns/op	    6316 B/op	      65 allocs/op
BenchmarkAddingFields/slogutiljsonlogger-8                                	  377458	      2989.00 ns/op	    5222 B/op	      57 allocs/op
BenchmarkAddingFields/slogutiljsonlogger.LogAttrs-8                       	  329028	      3543.00 ns/op	    6316 B/op	      65 allocs/op
PASS
ok  	github.com/nickbryan/slogutil/benchmarks	96.698s
```
