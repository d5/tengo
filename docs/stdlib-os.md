# Module - "os"

```golang
os := import("os")
```

## Constants

- `o_rdonly`
- `o_wronly`
- `o_rdwr`
- `o_append`
- `o_create`
- `o_excl`
- `o_sync`
- `o_trunc`
- `mode_dir`
- `mode_append`
- `mode_exclusive`
- `mode_temporary`
- `mode_symlink`
- `mode_device`
- `mode_named_pipe`
- `mode_socket`
- `mode_setuid`
- `mode_setgui`
- `mode_char_device`
- `mode_sticky`
- `mode_irregular`
- `mode_type`
- `mode_perm`
- `seek_set`
- `seek_cur`
- `seek_end`
- `path_separator`
- `path_list_separator`
- `dev_null`

## Functions

- `args() => [string]`: returns command-line arguments, starting with the
  program name.
- `chdir(dir string) => error`: changes the current working directory to the
  named directory.
- `chmod(name string, mode int) => error`: changes the mode of the named file
  to mode.
- `chown(name string, uid int, gid int) => error`: changes the numeric uid and
  gid of the named file.
- `clearenv()`: deletes all environment variables.
- `environ() => [string]`: returns a copy of strings representing the
  environment.
- `exit(code int)`: causes the current program to exit with the given status
  code.
- `expand_env(s string) => string`: replaces ${var} or $var in the string
  according to the values of the current environment variables.
- `getegid() => int`: returns the numeric effective group id of the caller.
- `getenv(key string) => string`: retrieves the value of the environment
  variable named by the key.
- `geteuid() => int`: returns the numeric effective user id of the caller.
- `getgid() => int`: returns the numeric group id of the caller.
- `getgroups() => [int]/error`: returns a list of the numeric ids of groups
  that the caller belongs to.
- `getpagesize() => int`: returns the underlying system's memory page size.
- `getpid() => int`: returns the process id of the caller.
- `getppid() => int`: returns the process id of the caller's parent.
- `getuid() => int`: returns the numeric user id of the caller.
- `getwd() => string/error`: returns a rooted path name corresponding to the
  current directory.
- `hostname() => string/error`: returns the host name reported by the kernel.
- `lchown(name string, uid int, gid int) => error`: changes the numeric uid
  and gid of the named file.
- `link(oldname string, newname string) => error`: creates newname as a hard
  link to the oldname file.
- `lookup_env(key string) => string/false`: retrieves the value of the
  environment variable named by the key.
- `mkdir(name string, perm int) => error`: creates a new directory with the
  specified name and permission bits (before umask).
- `mkdir_all(name string, perm int) => error`: creates a directory named path,
  along with any necessary parents, and returns nil, or else returns an error.
- `read_file(name string) => bytes/error`: reads the contents of a file into
  a byte array
- `readlink(name string) => string/error`: returns the destination of the
  named symbolic link.
- `remove(name string) => error`: removes the named file or (empty) directory.
- `remove_all(name string) => error`: removes path and any children it
  contains.
- `rename(oldpath string, newpath string) => error`: renames (moves) oldpath
  to newpath.
- `setenv(key string, value string) => error`: sets the value of the
  environment variable named by the key.
- `stat(filename string) => FileInfo/error`: returns a file info structure
  describing the file
- `symlink(oldname string newname string) => error`: creates newname as a
  symbolic link to oldname.
- `temp_dir() => string`: returns the default directory to use for temporary
  files.
- `truncate(name string, size int) => error`: changes the size of the named
  file.
- `unsetenv(key string) => error`: unsets a single environment variable.
- `create(name string) => File/error`: creates the named file with mode 0666
  (before umask), truncating it if it already exists.
- `open(name string) => File/error`: opens the named file for reading. If
  successful, methods on the returned file can be used for reading; the
  associated file descriptor has mode O_RDONLY.
- `open_file(name string, flag int, perm int) => File/error`: is the
  generalized open call; most users will use Open or Create instead. It opens
  the named file with specified flag (O_RDONLY etc.) and perm (before umask),
  if applicable.
- `find_process(pid int) => Process/error`: looks for a running process by its
  pid.
- `start_process(name string, argv [string], dir string, env [string]) => Process/error`:
  starts a new process with the program, arguments and attributes specified by
  name, argv and attr. The argv slice will become os.Args in the new process,
  so it normally starts with the program name.
- `exec_look_path(file string) => string/error`: searches for an executable
  named file in the directories named by the PATH environment variable.
- `exec(name string, args...) => Command/error`: returns the Command to execute
  the named program with the given arguments.

## File

```golang
file := os.create("myfile")
file.write_string("some data")
file.close()
```

- `chdir() => true/error`: changes the current working directory to the file,
- `chown(uid int, gid int) => true/error`: changes the numeric uid and gid of
  the named file.
- `close() => error`: closes the File, rendering it unusable for I/O.
- `name() => string`: returns the name of the file as presented to Open.
- `readdirnames(n int) => [string]/error`: reads and returns a slice of names
  from the directory.
- `sync() => error`: commits the current contents of the file to stable storage.
- `write(bytes) => int/error`: writes len(b) bytes to the File.
- `write_string(string) => int/error`: is like 'write', but writes the contents
  of string s rather than a slice of bytes.
- `read(bytes) => int/error`: reads up to len(b) bytes from the File.
- `stat() => FileInfo/error`: returns a file info structure describing the file
- `chmod(mode int) => error`: changes the mode of the file to mode.
- `seek(offset int, whence int) => int/error`: sets the offset for the next
  Read or Write on file to offset, interpreted according to whence: 0 means
  relative to the origin of the file, 1 means relative to the current offset,
  and 2 means relative to the end.

## Process

```golang
proc := start_process("app", ["arg1", "arg2"], "dir", [])
proc.wait()
```

- `kill() => error`: causes the Process to exit immediately.
- `release() => error`: releases any resources associated with the process,
  rendering it unusable in the future.
- `signal(signal int) => error`: sends a signal to the Process.
- `wait() => ProcessState/error`: waits for the Process to exit, and then
  returns a ProcessState describing its status and an error, if any.

## ProcessState

```golang
proc := start_process("app", ["arg1", "arg2"], "dir", [])
stat := proc.wait()
pid := stat.pid()
```

- `exited() => bool`: reports whether the program has exited.
- `pid() => int`: returns the process id of the exited process.
- `string() => string`: returns a string representation of the process.
- `success() => bool`: reports whether the program exited successfully, such as
  with exit status 0 on Unix.

```golang
cmd := exec.command("echo", ["foo", "bar"])
output := cmd.output()
```

## FileInfo

- `name`: name of the file the info describes
- `mtime`: time the file was last modified
- `size`: file size in bytes
- `mode`: file permissions as in int, comparable to octal permissions
- `directory`: boolean indicating if the file is a directory

## Command

- `combined_output() => bytes/error`: runs the command and returns its combined
  standard output and standard error.
- `output() => bytes/error`: runs the command and returns its standard output.
- `run() => error`: starts the specified command and waits for it to complete.
- `start() => error`: starts the specified command but does not wait for it to
  complete.
- `wait() => error`: waits for the command to exit and waits for any copying to
  stdin or copying from stdout or stderr to complete.
- `set_path(path string)`: sets the path of the command to run.
- `set_dir(dir string)`: sets the working directory of the process.
- `set_env(env [string])`: sets the environment of the process.
- `process() => Process`: returns the underlying process, once started.
