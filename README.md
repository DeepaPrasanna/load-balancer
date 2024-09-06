# load-balancer
Implementation of load balancer in go using round robin as the load balancing algorithm. It also health checks the servers before forwarding requests to it. The health checks are sent periodically. If no health check period is provided, the default value is 10s. Custom value can be provided by a flag healthCheckPeriod.
For example, 
  ```
    go run --healthCheckPeriod=20

  ```

### Run the below steps to run the load balancer

1. start the three web servers as follows:
    ```
        python3 -m http.server 8080 --directory server-1
   ```
   ```
        python3 -m http.server 8081 --directory server-2
   ```

   ```
        python3 -m http.server 8082 --directory server-3
   ```
2. ```
      cd load-balancer
      go run .
   ```

3. Run  ```curl://localhost:8000 ``` or if you want to make concurrent requests run
   ```
      curl --parallel --parallel-immediate --parallel-max 3 --config urls.txt
   ```



     

