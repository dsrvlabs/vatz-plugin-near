@rootwarp @xellos00

Interface Type -> module 정리
CLI Command
(TBD)
config
rest of tasks
3.2 Health checker
3.3 Monitoring



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

<b> Special Conditions </b>

1. Those following Plugins must declare network as below if you are running a VATZ on Testnet
- near-metric-chunk-produce-rate
- near-metric-uptime
> ```./near-metric-chunk-produce-rate --network testnet```

### 2.1 Build first with following command
```
~$ make build
Build All Plugins
===================
=> building machine-status-cpu
=> building machine-status-disk
=> building machine-status-memory
=> building near-metric-alive
=> building near-metric-block-height
=> building near-metric-chunk-produce-rate
=> building near-metric-number-of-peer
=> building near-metric-uptime
===================
All Build Finished
```

### 2.2. Run all plugins with following command
  ```
~$ make start
Start All Plugins
===================
=> Starting Plugins machine-status-cpu
=> Starting Plugins machine-status-disk
=> Starting Plugins machine-status-memory
=> Starting Plugins near-metric-alive
=> Starting Plugins near-metric-block-height
=> Starting Plugins near-metric-chunk-produce-rate
=> Starting Plugins near-metric-number-of-peer
=> Starting Plugins near-metric-uptime
===================
All Plugins are started!
```
### 2.3. Stop all plugins
```
~$ make stop 
Stopping All Plugins
===================
=> Stopping Plugins: machine-status-cpu in PID: 98633
=> Stopping Plugins: machine-status-disk in PID: 98636
=> Stopping Plugins: machine-status-memory in PID: 98639
=> Stopping Plugins: near-metric-alive in PID: 98642
=> Stopping Plugins: near-metric-block-height in PID: 98645
=> Stopping Plugins: near-metric-chunk-produce-rate in PID: 98648
=> Stopping Plugins: near-metric-number-of-peer in PID: 98651
=> Stopping Plugins: near-metric-uptime in PID: 98654
===================
All Plugins has stopped
```
### 2.4. Remove all plugin's build
```
~$ make clean
Cleaning All Plugins
===================
=> cleaning machine-status-cpu
=> cleaning machine-status-disk
=> cleaning machine-status-memory
=> cleaning near-metric-alive
=> cleaning near-metric-block-height
=> cleaning near-metric-chunk-produce-rate
=> cleaning near-metric-number-of-peer
=> cleaning near-metric-uptime
===================
All Plugins Cleaned
```

## 3. How to develop Plugins
For Plugins Developments
1. Create a new plugin folder with plugin's name under plugins.
2. Please follow the plugins name rule for two types
- For Machine
  - machine-{action_target}-{source}
- For protocol
  - protocol-{action_target}-{source}
  
3.Check VATZ SDK reference for development at https://github.com/dsrvlabs/vatz/tree/main/sdk


## 4. How to manage vatz-plugin-near
(TBD)