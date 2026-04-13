# dolitri

CLI tool that exports invoices and expense reports from a [Dolibarr](https://www.dolibarr.org/) instance for a given quarter into a ZIP archive.

## Usage

```
dolitri <quarter>
dolitri -version
```

`<quarter>` format: `yyyy-T1` through `yyyy-T4`

**Examples:**

```sh
dolitri 2025-T1   # Jan–Mar 2025
dolitri 2025-T4   # Oct–Dec 2025
dolitri -version  # print version and exit
```

The output is written to `<quarter>.zip` in the current directory with the following structure:

```
2025-T1.zip
├── facturas emitidas/
│   └── <invoice files>
└── <expense-report-ref>/
    ├── <expense report PDF>
    └── facturas/
        └── <uploaded expense receipts>
```

## Configuration

Credentials can be provided via a config file or environment variables.

### `~/.doliconf`

```ini
url    = https://your-dolibarr-instance.example.com
apikey = your_api_key_here
```

Lines starting with `#` and blank lines are ignored.

### Environment variables

| Variable     | Description                        |
|--------------|------------------------------------|
| `DOLAPIURL`  | Base URL of the Dolibarr instance  |
| `DOLAPIKEY`  | Dolibarr REST API key              |

Environment variables take precedence over `~/.doliconf`.

## Installation

Download the binary for your platform from the [Releases](../../releases) page.

| File | Platform |
|------|----------|
| `dolitri-linux-amd64` | Linux x86-64 |
| `dolitri-linux-arm64` | Linux ARM 64-bit |
| `dolitri-linux-armv6` | Raspberry Pi (armv6) |
| `dolitri-darwin-amd64` | macOS Intel |
| `dolitri-darwin-arm64` | macOS Apple Silicon |
| `dolitri-windows-amd64.exe` | Windows x86-64 |

Make the binary executable (Linux/macOS):

```sh
chmod +x dolitri-linux-amd64
```
