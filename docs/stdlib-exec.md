# Module - "exec"

```golang
exec := import("exec")
```

## Module Functions

- `look_path(file string) => string/error`: port of `exec.LookPath` function
- `command(name string, args array(string)) => `Cmd/error`: port of `exec.Command` function

## Cmd Functions

```golang
cmd := exec.command("echo", ["foo", "bar"])
output := cmd.output()
```

- `combined_output() => bytes/error`: port of `exec.Cmd.CombinedOutput` function
- `output() => bytes/error`: port of `exec.Cmd.Output` function
- `combined_output() => bytes/error`: port of `exec.Cmd.CombinedOutput` function
- `run() => error`: port of `exec.Cmd.Run` function
- `start() => error`: port of `exec.Cmd.Start` function
- `wait() => error`: port of `exec.Cmd.Wait` function
- `set_path(path string)`: sets `Path` of `exec.Cmd`
- `set_dir(dir string)`: sets `Dir` of `exec.Cmd`
- `set_env(env array(string))`: sets `Env` of `exec.Cmd`
- `process() => Process`: returns Process (`Process` of `exec.Cmd`)