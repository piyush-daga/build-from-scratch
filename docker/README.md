## Docker

Build a small docker clone that can fetch images from the docker registry and run them.

In doing so, we implement:
- **Chroot**: To isolate the execution environment, and not allow the container to view the host fs
- **Namespaces**: Isolate the container from seeing the parent's running processes. Any new process inside the container starts with a PID of 1
- Wiring up the std(in/out/err) to the parent
- Passing through the exit codes from child to parent

