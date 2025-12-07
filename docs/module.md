
## Module System

akashaProxy provides a module system that allows users to extend functionality through scripts. Modules are located in the `module` folder under the runtime directory (typically `/data/clash/module/`)

### Directory Structure

Each module is a subfolder within the `module` directory. The directory name of the module serves as its identifier.

A typical module consists of the following files:

```text
module_name/
├── module.prop       # Module metadata (optional)
├── post-fs-data.sh   # Pre-start script
├── start.sh          # Post-startup script
├── stop.sh           # Stop script
└── bin/              # (Optional) Storage for binaries or other resources
```

During the startup and shutdown process, AkashaProxy executes module scripts in sequence:

1. **Before Startup**:
   - When the service starts, it first executes the `post-fs-data.sh` script for all modules
   - The system will wait for the post-fs-data.sh script to complete. Do not perform time-consuming operations in this script

2. **After startup**:
   After the proxy kernel is successfully booted, execute the `start.sh` script for all modules

3. **Stop**:
   - When the service stops (after stopping the kernel), execute the `stop.sh` script for all modules

### Script Instructions

- All scripts are executed using busybox ash