## 模块系统

akashaProxy 提供了一个模块系统，允许用户通过脚本扩展功能。模块位于运行目录下的 `module` 文件夹中（通常为 `/data/clash/module/`）

### 目录结构

每个模块都是 `module` 目录下的一个子文件夹。模块的目录名即为模块的标识符。

一个典型的模块包含以下文件：

```text
module_name/
├── module.prop       # 模块元数据（可选）
├── post-fs-data.sh   # 启动前脚本
├── start.sh          # 启动后脚本
├── stop.sh           # 停止脚本
└── bin/              # (可选) 存放二进制文件或其他资源
```

akashaProxy 在启动和停止过程中会按顺序执行模块脚本：

1. **启动前**:
   - 服务启动时，首先执行所有模块的 `post-fs-data.sh`
   - 系统会等待post-fs-data.sh执行完毕。请勿在此脚本中执行耗时操作

2. **启动后**:
   - 代理内核成功启动后，执行所有模块的 `start.sh`

3. **停止**:
   - 服务停止时（停止内核后），执行所有模块的 `stop.sh`

### 脚本说明

- 所有脚本均使用 busybox ash 执行