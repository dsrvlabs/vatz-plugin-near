# vatz-plugin-near
Plugins for near protocol.
> Every individual plugins running to retrieve a single metric value for its purpose.

## 1. Plugins list 

### 1.1. Machine state plugins
- machine-status-cpu : localhost:9091
- machine-status-disk : localhost:9092
- machine-status-memory : localhost:9093 

### 1.2. Protocol state plugins
- near-metric-blockheight : localhost:10001
- near-metric-up : localhost:10002
  

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
.Build first with following command
  ```
  make stop
  ```

## 3. How to add Makefile for Run all plugins
(TBD)