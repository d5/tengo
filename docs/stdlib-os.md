# Module - "os"

```golang
os := import("os")
```

## Module Variables

- `o_rdonly`: equivalent of Go's `os.O_RDONLY`
- `o_wronly`: equivalent of Go's `os.O_WRONLY`
- `o_rdwr`: equivalent of Go's `os.O_RDWR`
- `o_append`: equivalent of Go's `os.O_APPEND`
- `o_create`: equivalent of Go's `os.O_CREATE`
- `o_excl`: equivalent of Go's `os.O_EXCL`
- `o_sync`: equivalent of Go's `os.O_SYNC`
- `o_trunc`: equivalent of Go's `os.O_TRUNC`
- `mode_dir`: equivalent of Go's `os.ModeDir`
- `mode_append`: equivalent of Go's `os.ModeAppend`
- `mode_exclusive`: equivalent of Go's `os.ModeExclusive`
- `mode_temporary`: equivalent of Go's `os.ModeTemporary`
- `mode_symlink`: equivalent of Go's `os.ModeSymlink`
- `mode_device`: equivalent of Go's `os.ModeDevice`
- `mode_named_pipe`: equivalent of Go's `os.ModeNamedPipe`
- `mode_socket`: equivalent of Go's `os.ModeSocket`
- `mode_setuid`: equivalent of Go's `os.ModeSetuid`
- `mode_setgui`: equivalent of Go's `os.ModeSetgid`
- `mode_char_device`: equivalent of Go's `os.ModeCharDevice`
- `mode_sticky`: equivalent of Go's `os.ModeSticky`
- `mode_irregular`: equivalent of Go's `os.ModeIrregular`
- `mode_type`: equivalent of Go's `os.ModeType`
- `mode_perm`: equivalent of Go's `os.ModePerm`
- `seek_set`: equivalent of Go's `os.SEEK_SET`
- `seek_cur`: equivalent of Go's `os.SEEK_CUR`
- `seek_end`: equivalent of Go's `os.SEEK_END`
- `path_separator`: equivalent of Go's `os.PathSeparator`
- `path_list_separator`: equivalent of Go's `os.PathListSeparator`
- `dev_null`: equivalent of Go's `os.DevNull`

## Module Functions

- `args() => array(string)`: returns `os.Args`
- `chdir(dir string) => error`: port of `os.Chdir` function
- `chmod(name string, mode int) => error `: port of Go's `os.Chmod` function
- `chown(name string, uid int, gid int) => error `: port of Go's `os.Chown` function
- `clearenv() `: port of Go's `os.Clearenv` function
- `environ() => array(string) `: port of Go's `os.Environ` function
- `executable() => string/error`: port of Go's `os.Executable()` function
- `exit(code int) `: port of Go's `os.Exit` function
- `expand_env(s string) => string `: port of Go's `os.ExpandEnv` function
- `getegid() => int `: port of Go's `os.Getegid` function
- `getenv(s string) => string `: port of Go's `os.Getenv` function
- `geteuid() => int `: port of Go's `os.Geteuid` function
- `getgid() => int `: port of Go's `os.Getgid` function
- `getgroups() => array(string)/error `: port of Go's `os.Getgroups` function
- `getpagesize() => int `: port of Go's `os.Getpagesize` function
- `getpid() => int `: port of Go's `os.Getpid` function
- `getppid() => int `: port of Go's `os.Getppid` function
- `getuid() => int `: port of Go's `os.Getuid` function
- `getwd() => string/error `: port of Go's `os.Getwd` function
- `hostname() => string/error `: port of Go's `os.Hostname` function
- `lchown(name string, uid int, gid int) => error `: port of Go's `os.Lchown` function
- `link(oldname string, newname string) => error `: port of Go's `os.Link` function
- `lookup_env(key string) => string/false`: port of Go's `os,LookupEnv` function
- `mkdir(name string, perm int) => error `: port of Go's `os.Mkdir` function
- `mkdir_all(name string, perm int) => error `: port of Go's `os.MkdirAll` function
- `readlink(name string) => string/error `: port of Go's `os.Readlink` function
- `remove(name string) => error `: port of Go's `os.Remove` function
- `remove_all(name string) => error `: port of Go's `os.RemoveAll` function
- `rename(oldpath string, newpath string) => error `: port of Go's `os.Rename` function
- `setenv(key string, value string) => error `: port of Go's `os.Setenv` function
- `symlink(oldname string newname string) => error `: port of Go's `os.Symlink` function
- `temp_dir() => string `: port of Go's `os.TempDir` function
- `truncate(name string, size int) => error `: port of Go's `os.Truncate` function
- `unsetenv(key string) => error `: port of Go's `os.Unsetenv` function
- `user_cache_dir() => string/error `: port of Go's `os.UserCacheDir` function
- `create(name string) => File/error`: port of Go's `os.Create` function
- `open(name string) => File/error`: port of Go's `os.Open` function
- `open_file(name string, flag int, perm int) => File/error`: port of Go's `os.OpenFile` function
- `find_process(pid int) => Process/error`: port of Go's `os.FindProcess` function
- `start_process(name string, argv array(string), dir string, env array(string)) => Process/error`: port of Go's `os.StartProcess` function

## File Functions

```golang
file := os.create("myfile")
file.write_string("some data")
file.close()
```

- `chdir() => true/error`: port of `os.File.Chdir` function
- `chown(uid int, gid int) => true/error`: port of `os.File.Chown` function
- `close() => error`: port of `os.File.Close` function
- `name() => string`: port of `os.File.Name` function
- `readdirnames() => array(string)/error`: port of `os.File.Readdirnames` function
- `sync() => error`: port of `os.File.Sync` function
- `write(bytes) => int/error`: port of `os.File.Write` function
- `write_string(string) => int/error`: port of `os.File.WriteString` function
- `read(bytes) => int/error`: port of `os.File.Read` function
- `chmod(mode int) => error`: port of `os.File.Chmod` function
- `seek(offset int, whence int) => int/error`: port of `os.File.Seek` function

## Process Functions

```golang
proc := start_process("app", ["arg1", "arg2"], "dir", [])
proc.wait()
```

- `kill() => error`: port of `os.Process.Kill` function
- `release() => error`: port of `os.Process.Release` function
- `signal(signal int) => error`: port of `os.Process.Signal` function
- `wait() => ProcessState/error`: port of `os.Process.Wait` function

## ProcessState Functions

```golang
proc := start_process("app", ["arg1", "arg2"], "dir", [])
stat := proc.wait()
pid := stat.pid()
```

- `exited() => bool`: port of `os.ProcessState.Exited` function
- `pid() => int`: port of `os.ProcessState.Pid` function
- `string() => string`: port of `os.ProcessState.String` function
- `success() => bool`: port of `os.ProcessState.Success` function