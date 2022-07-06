# vatz-plugin-near
> This is Plugins for VATZ that check status of NEAR protocol.
>
> Every individual plugins running to retrieve a single metric value for its purpose. 
> Please, be careful with Plugin Port number

## 1. Plugins list 
> Total 8 single module plugins.
> 
> - Plugins for <b> Machine State </b> starts with port number from `9001` to `9999`
> - Plugins for <b> Protocol State </b> state starts with port number from `10001` to `11000` 


### 1.1. Machine state plugins
- `machine-status-cpu`: Plugin that retrieves machine's CPU usage info
   - default address: localhost
   - default port: 9091
   - method name: `GetMachineCPUUsage`
- `machine-status-disk`: Plugin that retrieves machine's expected path's disk usage info
   - default address: localhost
   - default port: 9092
   - method name: `GetMachineDiskUsage`
- `machine-status-memory`: Plugin that retrieves machine's memory usage info
   - default address: localhost
   - default port: 9093
   - method name: `GetMachineMemoryUsage`
    
### 1.2. Protocol state plugins
- `near-metric-alive`: Plugin that retrieves info that NEAR node is alive.
  - default address: localhost
  - default port: 10001
  - method name: `NearGetAlive`
- `near-metric-block-height`: Plugin that retrieves NEAR block height within 5 seconds
   - default address: localhost
   - default port: 10002
   - method name: `NearGetBlockHeight`
- `near-metric-chunk-produce-rate`: Plugin that retrieves NEAR current Chunk produce rate per epoch
   - default address: localhost
   - default port: 10003
   - method name: `NearGetChunkProduceRate`
- `near-metric-number-of-peer`: Plugin that retrieves NEAR current connected peers status
   - default address: localhost
   - default port: 10004
   - method name: `NearGetNumberOfPeer`
- `near-metric-up-time`: Plugin that retrieves NEAR node's average uptime during epoch
   - default address: localhost
   - default port: 10005
   - method name: `NearGetUptime`

## 2. How to run the plugins

### 2.1 Build first with following command
  ```
  make build
  ```
### 2.2. Run all plugins with following command
  ```
  make start
  ```
### 2.3. Stop all plugins
  ```
  make stop
  ```
### 2.4. Remove all plugin's build
  ```
  make clean
  ```

## 3. How to develop Plugins








## 4. How to manage vatz-plugin-near