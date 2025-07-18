# Prayers CLI

A command-line tool to check Islamic prayer times by city
<br>

<div align="center">
  <img src="image-3.png" alt="Prayers CLI Preview" width="500"/>
</div>

## Features

- Prayer times by city
- Shows upcoming prayer
- Time-based ASCII art (Morning, Night)
- Cozy, eye-comfortable terminal style
- Remembers your last configured city

## Preview

## Prerequisites

Make sure you have Go installed.

To verify Go is installed:

```bash
go version
```

## Installation

```bash
# Clone the repo
git clone https://github.com/OmarAshour02/prayer-cli.git
cd prayer-cli
go build -o prayers
```

### Option 1: Global install

```bash
# Move the binary to a directory in your PATH
sudo mv prayers /usr/local/bin/
```

### Option 2: Run locally

```bash
cd prayer-cli
./prayers
```

## Usage

After installation (assuming global installation), run:

```bash
prayers
```

This will fetch prayer times using default city "Makkah", unless you configure your selected city

### Set or change your saved city

Use the --city or -c flag to set (or change) your default city:

```bash
prayers --city "Dublin"
```

Or shorthand:

```bash
prayers -c "Cairo"
```

Once set, the city will be remembered for future runs.
