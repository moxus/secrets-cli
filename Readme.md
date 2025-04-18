# secrets-cli

A command-line tool to manage encrypted key-value secrets using different storage backends (sqlite, jsonfile, mongodb-placeholder).

## Environment Variable

### `SECRETBOX_KEY`

- **Required**: Yes
- **Description**: The encryption key used for encrypting and decrypting secrets.
- **Format**: Base64-encoded 32-byte value.
- **Example**:

  ```sh
  export SECRETBOX_KEY="your-base64-encoded-key-here"
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
export SECRETBOX_KEY="your-base64-key"
secrets-cli create mykey myvalue
secrets-cli read mykey
secrets-cli list
secrets-cli delete mykey
```
