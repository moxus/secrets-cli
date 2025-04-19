# secrets-cli

A command-line tool to manage encrypted key-value secrets using different storage backends (sqlite, jsonfile, mongodb-placeholder).

## Install

install the latest binary with this command (only macos and linux)

```bash
curl -sSL https://raw.githubusercontent.com/moxus/secrets-cli/refs/heads/main/install.sh | bash
```

## Environment Variable

### `SECRETS_ENCRYPTION_KEY`

- **Required**: Yes
- **Description**: The encryption key used for encrypting and decrypting secrets.
- **Format**: Base64-encoded 32-byte value.
- **Example**:

  ```sh
  export SECRETS_ENCRYPTION_KEY="your-base64-encoded-key-here"
  ```

- If the key is shorter than 32 bytes, it will be zero-padded. If longer, it will be truncated. The variable must be set.

## Configuration File

### `~/.secrets-cli.json`

- **Purpose**: Provides default values for backend configuration parameters.
- **Location**: User's home directory.
- **Example**:

  ```json
  {
    "backend_type": "sqlite",
    "sqlite_db_path": "/Users/youruser/secrets.db",
    "json_file_path": "/Users/youruser/secrets.json",
    "mongo_uri": "",
    "mongo_database": "",
    "mongo_collection": ""
  }
  ```

- **Fields**:
  - `backend_type`: `"sqlite"`, `"jsonfile"`, or `"mongodb-placeholder"`
  - `sqlite_db_path`: Path to SQLite database file
  - `json_file_path`: Path to JSON file for secrets
  - `mongo_uri`: MongoDB connection URI
  - `mongo_database`: MongoDB database name
  - `mongo_collection`: MongoDB collection name

## Program Parameters

All parameters can be set via command-line flags, which override config file values.

### Global Flags

- `--backend`  
  Storage backend type (`sqlite`, `jsonfile`, `mongodb-placeholder`)

- `--sqlite-db`  
  SQLite database file path

- `--json-file`  
  JSON file path

- `--mongo-uri`  
  MongoDB connection URI

- `--mongo-db`  
  MongoDB database name

- `--mongo-collection`  
  MongoDB collection name

### Commands

- `create [key] [value] [--update]`  
  Create a new secret. Use `--update` to update if the key exists.

- `read [key]`  
  Read and decrypt a secret by key.

- `delete [key]`  
  Delete a secret by key.

- `list`  
  List all secret keys.

## Example Usage

```sh
export SECRETS_ENCRYPTION_KEY="your-base64-key"
secrets-cli create mykey myvalue
secrets-cli read mykey
secrets-cli list
secrets-cli delete mykey
```

## Generate Command

The `generate` (alias: `gen`) command creates a random password of a specified length and stores it as a secret under the given key.

### Usage

```sh
secrets-cli generate [key] [length] [flags]
```

- `[key]`: The key under which to store the generated password.
- `[length]`: The length of the password to generate.

### Flags

- `-u`, `--uppercase`  
  Include uppercase letters (A-Z) in the generated password.

- `-l`, `--lowercase`  
  Include lowercase letters (a-z) in the generated password.

- `-n`, `--numbers`  
  Include numbers (0-9) in the generated password.

- `--update`  
  Update the secret if it already exists for the given key.

### Example

Generate a 20-character password with uppercase, lowercase, and numbers, and store it under the key `db_password`:

```sh
secrets-cli generate db_password 20 --uppercase --lowercase --numbers
```

If you want to overwrite an existing secret:

```sh
secrets-cli generate db_password 20 --uppercase --lowercase --numbers --update
```
