# Blitz ⚡

A simple, fast, and secure command-line tool for transferring files.

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Table of Contents

- [Key Features](#key-features)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Contributing](#contributing)

## Key Features

*   ✨ **Simple & Intuitive:** A clean command-line interface designed for ease of use.
*   📚 **Built-in Help:** Simply type `blitz` or `blitz --help` to view all available commands and options.
*   🔒 **Secure & Fast:** All transfers are encrypted using TLS to keep your data safe while being sent at high speed.

## Usage

Getting help is easy. Just run the main command:
```powershell
blitz --help
```

### Sending a File

To send a file, use the `send` command:
```powershell
# Example
blitz send --file "path/to/your/document.pdf"
```

## Roadmap

Here are the features and updates planned for the future.

- [ ] Add OAuth 2.0 to validate users via their Gmail identity.
- [ ] Host the public-facing server to enable transfers without self-hosting.
- [ ] Add support for transferring entire directories.

## Contributing

Contributions are welcome! If you have a suggestion or find a bug, please feel free to open an issue or submit a pull request. 
