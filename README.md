# Project Overview

This project includes a Go application that processes shell scripts and generates corresponding Python scripts. The Python scripts execute `curl` commands with dynamic parameters.

## Project Structure

- `main.go`: The main Go application that reads a shell script, extracts `curl` commands, and generates a Python script.
- `py.tmpl`: The template used for generating Python scripts.

## Prerequisites

- Go 1.16 or later
- Python 3.6 or later

## Usage

1. **Build the Go application:**

    ```sh
    go build -o curl2py main.go
    ```

2. **Run the Go application with a shell script as input:**

   ```sh
   ./curl2py <inputfile>.sh
   ```

   This will generate a Python script named `<inputfile>_gen.py`.

3. **Execute the generated Python script:**

   ```sh
   python3 <inputfile>_gen.py
   ```

## Input File Format

The input file should be a shell script containing `curl` commands. Each `curl` command can include placeholders for dynamic parameters, enclosed in double curly braces `{{ }}`. These placeholders will be replaced with actual values when the corresponding Python function is generated. The generated Python function will have a name derived from the comment preceding the `curl` command and parameters corresponding to the placeholders.

### Example

Given a shell script `example.sh` with the following content:

```sh
# RunJob:
curl 'https://example.com/api?param1={{param1}}&param2={{param2}}' -H 'Authorization: Bearer {{token}}'
```
