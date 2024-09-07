# Cellular Automata

This project implements a parallelized cellular automaton simulation with support for multiple neighborhood types and video output.

## Requirements

- Go 1.16 or later
- FFmpeg (for video output)

## Building

To build the project into a binary, run the following command from the project root:

```
go build -o cellular-automata ./cmd/ca
```

This will create an executable named `cellular-automata` (or `cellular-automata.exe` on Windows) in the current directory.

## Usage

To run the cellular automaton simulation:

```
./cellular-automata [flags]
```

Flags:
- `-size int`: Size of the cellular automata grid (default 100)
- `-steps int`: Number of steps to run (default 100)
- `-video`: Save output as video (default false)

Examples:
- Run with default settings:
  ```
  ./cellular-automata
  ```
- Run with a 200x200 grid for 150 steps and save as video:
  ```
  ./cellular-automata -size 200 -steps 150 -video
  ```

## Project Structure

- `cmd/ca/main.go`: Main entry point of the application
- `internal/automaton/automaton.go`: Implementation of the cellular automaton
- `internal/video/video.go`: Video creation utilities
