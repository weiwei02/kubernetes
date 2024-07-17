
- [v1.24.17-baidu-0716]
  - [feature] 增加 apiserver `watch_wait_seconds` 和 `resource_version` 度量
  - [feature] 支持环境变量 `APISERVER_GET_CURRENT_RESOURCE_VERSION_INTERVAL` 设置定时从 etcd 获取当前资源版本周期任务时间间隔
  - [feature] 新特性： 支持 etcd WatchProgress 进度查询和通知特性，优化对更新频率极低对象的 rv 匹配逻辑。