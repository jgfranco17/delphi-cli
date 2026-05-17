# Delphi CLI

A small Go CLI (`delphi`) that provides utilities for Claude Code workflows.

## Commands

### `delphi statusline`

Reads Claude Code's JSON hook payload from stdin and renders a three-line,
color-coded statusline:

```console
Using claude-opus-4-5 in ~/dev/my-project
Usage: ▓▓▓▓░░░░░░ 42% | ~¥150 equiv | 5h: 30% 7d: 15%
Git: main ✓ clean
```

#### Flags

Delphi CLI supports the following flags for `statusline`, which can be combined as needed:

| Flag | Default | Description |
| ------ | --------- | ------------- |
| `--all` | `false` | Show all statusline sections |
| `--cwd` | `false` | Show current working directory |
| `--git` | `false` | Show git information |
| `--stats` | `false` | Show usage statistics |
| `--disable-colors` | `false` | Disable ANSI color output formatting |

**Example** (Claude Code `statusLine` config):

```json
{
    "statusLine": {
        "type": "command",
        "command": "delphi statusline --all",
        "padding": 2
    },
}
```

## Development

Requires [just](https://just.systems).

| Command | Description |
| --------- | ------------- |
| `just build` | Build the binary for the current platform |
| `just install-local` | Build and install to `~/.claude/tools/bin/` |
| `just uninstall-local` | Remove the installed binary |
| `just test` | Run tests with race detection and coverage |
| `just tidy` | Sync Go modules |
| `just update-deps` | Upgrade all dependencies |
