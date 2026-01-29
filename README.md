# ya-metrics

Репозиторий для выполнения задания трека "Сервис сбора метрик и алертинга".

## Результат профайлинга после оптимизации gzip.Writer

```bash
-> % go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
File: main
Type: alloc_space
Time: 2026-01-29 02:25:58 MSK
Showing nodes accounting for -49147.89MB, 97.11% of 50609.61MB total
Dropped 222 nodes (cum <= 253.05MB)
      flat  flat%   sum%        cum   cum%
-40582.88MB 80.19% 80.19% -49382.07MB 97.57%  compress/flate.NewWriter (inline)
-8565.01MB 16.92% 97.11% -8565.01MB 16.92%  compress/flate.(*compressor).initDeflate (inline)
      -1MB 0.002% 97.11% -49751.96MB 98.31%  github.com/yogenyslav/ya-metrics/internal/server.NewServer.WithLogging.func1.1
    0.50MB 0.00099% 97.11% -11433.54MB 22.59%  github.com/yogenyslav/ya-metrics/internal/server/handler.(*Handler).GetMetricJSON
    0.50MB 0.00099% 97.11% -50128.07MB 99.05%  net/http.(*conn).serve
         0     0% 97.11%  -276.93MB  0.55%  bufio.(*Writer).Flush
         0     0% 97.11% -8799.19MB 17.39%  compress/flate.(*compressor).init
         0     0% 97.11% -23756.27MB 46.94%  compress/gzip.(*Writer).Close
         0     0% 97.11% -49382.07MB 97.57%  compress/gzip.(*Writer).Write
         0     0% 97.11% -49778.49MB 98.36%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 97.11% -25943.16MB 51.26%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 97.11% -49713.43MB 98.23%  github.com/yogenyslav/ya-metrics/internal/server.NewServer.WithCompression.func2.1
         0     0% 97.11% -11549.27MB 22.82%  github.com/yogenyslav/ya-metrics/internal/server/handler.(*Handler).GetMetricRaw
         0     0% 97.11% -2935.33MB  5.80%  github.com/yogenyslav/ya-metrics/internal/server/handler.(*Handler).ListMetrics
         0     0% 97.11% -23756.27MB 46.94%  github.com/yogenyslav/ya-metrics/internal/server/middleware.(*compressionResponseWriter).Close
         0     0% 97.11% -25854.31MB 51.09%  github.com/yogenyslav/ya-metrics/internal/server/middleware.(*compressionResponseWriter).Write
         0     0% 97.11% -25942.66MB 51.26%  github.com/yogenyslav/ya-metrics/internal/server/middleware.WithSignature.func1.1
         0     0% 97.11%  -255.40MB   0.5%  io.Copy (inline)
         0     0% 97.11%  -255.40MB   0.5%  io.CopyN
         0     0% 97.11%  -255.40MB   0.5%  io.copyBuffer
         0     0% 97.11%  -255.40MB   0.5%  io.discard.ReadFrom
         0     0% 97.11%  -276.93MB  0.55%  net/http.(*chunkWriter).Write
         0     0% 97.11%  -275.93MB  0.55%  net/http.(*chunkWriter).writeHeader
         0     0% 97.11%  -285.95MB  0.57%  net/http.(*response).finishRequest
         0     0% 97.11% -49751.96MB 98.31%  net/http.HandlerFunc.ServeHTTP
         0     0% 97.11% -49778.49MB 98.36%  net/http.serverHandler.ServeHTTP
         0     0% 97.11%  -354.04MB   0.7%  sync.(*Pool).Get
```
