# idle

Command line utility to lower process priorities to IDLE (lowest priority)

Supports Linux, MacOS and MS Windows.

## Usage

```
usage: idle process [process ...]

Changes priority to Idle for matching processes.

Process is treated as a case-sensitive substring to match running processes against. When prefixed and suffixed (surrounded by) '/' characters, the pattern is interpreted as a regular expression.

Options:
  -V      Show version of this program
  -e str  Exclude specific pids from being idled. Argument is specified as a csv string: pid[,...]
  -i      Use case insensitive process name matching
  -p int  Process poll interval (in milliseconds) (default 500)
  -w      Watch processes and lower the priority of any matches the process expression(s)

```

## Examples

Lower CPU priority of all currently running processes called `ffmpeg`:

    idle ffmpeg

Lower priority of all currently running processes called `ffmpeg` except for a given one having PID 123:

    idle -e 123 ffmpeg

Lower priority of two different programs:

    idle ffmpeg gzip

Use case insensitive process name matching:

    idle -i fFMpeg.Exe

Lower priority based on a regular expression:

    idle /^(SlowSof|ff).+/

And the same with case insensitive matching:

    idle -i /^(SlowSof|ff).+/

Continuously watch for new processes and lower their priorities as they are spawned:

    idle -w ffmpeg


## Gotchas

By default the process name matching is case sensitive and absolute. That means on windows in particular, you have to provide the right extension as part of the process pattern. I.e. `ffmpeg.exe` since that's the name of the process. On windows, check the Task explorer if you don't know the proper name of the executable. On Mac or Linux, use `ps` instead.

## Other

The process polling flag controls how fast the program will react when watching for new processes. The default of 0.5 seconds ought to be fast enough for most uses, and incurs negligible CPU usage. But if you want it to act instantly, you can lower this value. It will however consume a bit more CPU if you lower it too much. 

Additionally, there are a few milliseconds of latency incurred during each poll loop, so 1 millisecond will likely just result in a tight spin-loop of a few milliseconds, but it's up to you to decide on this trade-off between responsiveness and CPU usage of the program.
