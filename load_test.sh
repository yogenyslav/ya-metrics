PARALLEL=10
RATE_LIMIT=100
DURATION=10s

hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m GET http://localhost:8080/ > reports/list_metrics.txt &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m POST -d '{"id":"counter_metric","type":"counter","delta":1}' http://localhost:8080/update/ > reports/update_counter_raw.txt && sleep 1s &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m POST -d '{"id":"gauge_metric","type":"gauge","value":123.45}' http://localhost:8080/update/ > reports/update_gauge_raw.txt && sleep 1s &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m POST -d '{"id":"counter_metric","type":"counter"}' http://localhost:8080/value/ > reports/update_counter_json.txt && sleep 1s &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m POST -d '{"id":"gauge_metric","type":"gauge"}' http://localhost:8080/value/ > reports/update_gauge_json.txt && sleep 1s &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m POST -d '[{"id":"counter_metric","type":"counter","delta":1},{"id":"gauge_metric","type":"gauge","value":123.45}]' http://localhost:8080/updates/ > reports/update_metrics_batch.txt && sleep 1s &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m GET http://localhost:8080/value/counter/counter_metric > reports/get_counter_raw.txt && sleep 1s &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m GET http://localhost:8080/value/gauge/gauge_metric > reports/get_gauge_raw.txt && sleep 1s &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m POST -d '{"id":"counter_metric","type":"counter","delta":5}' http://localhost:8080/update/ > reports/update_counter_json.txt && sleep 1s &
hey -z "$DURATION" -c "$PARALLEL" -q "$RATE_LIMIT" -m POST -d '{"id":"gauge_metric","type":"gauge","value":678.90}' http://localhost:8080/update/ > reports/update_gauge_json.txt && sleep 1s &
echo "waiting for load test to finish"
sleep "$DURATION"
echo "load test finished"

mkdir -p profiles
touch profiles/$1.pprof
curl -o profiles/$1.pprof http://localhost:8080/debug/pprof/allocs