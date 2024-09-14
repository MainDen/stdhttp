# STDHTTP

STDHTTP is a Go application that provides an HTTP interface for standard streams. It allows redirecting `stdout` and `stderr` of running processes to specified HTTP endpoints and also manages the processes themselves using a built-in broker.

## Features

- Redirect `stdout` and `stderr` to specified URLs via HTTP.
- Manage processes started with `stdhttp run` through commands like `stdhttp list` and `stdhttp kill`.
- Start a debug server to monitor and manage processes in debug mode.

## Installation

You can download pre-built binaries for your operating system from the [GitHub releases page](https://github.com/MainDen/stdhttp/releases).

### Steps:

1. Go to the [Releases](https://github.com/MainDen/stdhttp/releases) page.
2. Download the appropriate archive for your OS:
   - For Linux: `stdhttp-linux-amd64.tar.gz`
   - For Windows: `stdhttp-windows-amd64.tar.gz`
3. Extract the archive:
   - For Linux:
     ```bash
     tar -xvf stdhttp-linux-amd64.tar.gz
     ```
   - For Windows, extract the ZIP using your preferred tool (e.g., WinRAR, 7-Zip or built-in Windows Explorer functionality).
4. After extracting, you should have the following binaries:
   - **Linux**: `stdhttp`
   - **Windows**: `stdhttp.exe` and `stdhttpd.exe`
5. (Optional) Move the binary to a directory in your `$PATH`:
   - For Linux:
     ```bash
     sudo mv stdhttp /usr/local/bin/
     ```
   - For Windows, add the binary path to your system environment variables.

## Usage

The STDHTTP application has three main commands: `run`, `broker`, and `debug`.

### Running a command with `stdout` and `stderr` redirected to HTTP

You can redirect the standard streams `stdout` and `stderr` of a command to specified HTTP URLs using the `--stdout-url` and `--stderr-url` options:

```bash
stdhttp run --stdout-url URL --stderr-url URL COMMAND [ARG ...]
```

Example:

```bash
stdhttp run --stdout-url http://localhost:8080/stdout --stderr-url http://localhost:8080/stderr ls -la
```

### Managing processes via the broker
If you start the broker with the `stdhttp broker` command, you can monitor and manage the processes you run:

Start the broker:

```bash
stdhttp broker
```

List active processes (wait for a 10 seconds after starting the broker):

```bash
stdhttp list
```

Run a command:

```bash
stdhttp run COMMAND [ARG ...]
```

Kill a process by its PID or kill processes by a pattern:

```bash
stdhttp kill {PID|PATTERN}
```

If broker goes down, you can restart it with the `stdhttp broker` command. Also if you want to see the broker in the list of processes, you can run it with `stdhttp run` command:
   
```bash
stdhttp run stdhttp broker
```

Or you can run the broker in the background:

```bash
stdhttpd run stdhttp broker
```

### Debug Mode
You can also start a debug server to track requests in debug mode. To do this, use the command:

```bash
stdhttp debug
```

Once the debug server is running, you can run commands with the --debug flag to send requests to the debug server:

```bash
stdhttp run --debug COMMAND [ARG ...]
```

## Background Mode

You can run `stdhttp` in the background on both Linux and Windows.

### Linux

To run `stdhttp` in the background on Linux, you can use the `&` symbol after the command.

Example:

```bash
stdhttp run --stdout-url http://localhost:8080/stdout --stderr-url http://localhost:8080/stderr COMMAND [ARG ...] &
```

You can also use standard tools like `nohup` or `screen` for more advanced background management if needed.

### Windows

On Windows, to run `stdhttp` in background mode, you should use the `stdhttpd.exe` binary instead of `stdhttp.exe`. This daemon version (`stdhttpd.exe`) will run the process in the background without holding up the terminal.

Example:

```bash
stdhttpd.exe run --stdout-url http://localhost:8080/stdout --stderr-url http://localhost:8080/stderr COMMAND [ARG ...]
```

This allows you to start processes and redirect their output without occupying the terminal window.

## Configuration

STDHTTP uses options and environment variables for configuration. For more details, run `stdhttp --help`.

## License
STDHTTP is distributed under the BSD 3-Clause License. For more details, see the `LICENSE.md` file.
