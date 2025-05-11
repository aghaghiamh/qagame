# Accessing the Visual pprof Diagram
1. **Download the pprof stats**:
    ```bash
    curl http://pprof_host:pprof_port/debug/pprof/{goroutine} --output ./pprof_stats/0001_stat.o
    ```

2. **Visualize the pprof stats**:
    ```bash
    go tool pprof -http=localhost:8077 ./0001_stat.o
    ``` 