# Fingerprint Requester - Node.js Wrapper

Axios-like HTTP client with custom TLS fingerprint support.

## Installation

```bash
npm install
```

## Quick Start

```javascript
import { create } from './requester.js';

const requester = create();

// GET request
const response = await requester.get('https://api.example.com/data');
console.log(response.data);

// POST request
const result = await requester.post('https://api.example.com/users', {
  name: 'John',
  email: 'john@example.com'
});
```

## Configuration

```javascript
const requester = create({
  binDir: './bin',                    // Binary directory (default: './bin')
  binaryPath: './custom/path/binary', // Manual binary path (optional)
  configPath: './config.json',        // TLS config file (default: './config.json')
  timeout: 30000,                     // Default timeout in ms (default: 30000)
});
```

## API Methods

### GET Request
```javascript
const response = await requester.get(url, config);
```

### POST Request
```javascript
const response = await requester.post(url, data, config);
```

### Other Methods
```javascript
await requester.put(url, data, config);
await requester.delete(url, config);
await requester.patch(url, data, config);
await requester.head(url, config);
await requester.options(url, config);
```

### Generic Request
```javascript
const response = await requester.request({
  method: 'POST',
  url: 'https://api.example.com',
  headers: { 'Authorization': 'Bearer token' },
  data: { key: 'value' },
  timeout: 5000,
  responseType: 'json', // 'text' or 'json'
});
```

### Stream Request
```javascript
await requester.stream({
  method: 'GET',
  url: 'https://api.example.com/stream',
  onData: ({ data }) => {
    console.log('Chunk received:', data);
  },
});
```

## Response Object

```javascript
{
  data: '...',           // Response body
  status: 200,           // HTTP status code
  statusText: 'OK',      // HTTP status text
  headers: {},           // Response headers
  config: {}             // Request config
}
```

## Platform Detection

The wrapper automatically detects your platform and architecture:
- **Windows AMD64**: `fingerprint_windows_amd64.exe`
- **Linux AMD64**: `fingerprint_linux_amd64`
- **Android ARM64**: `fingerprint_android_arm64`

If auto-detection fails, manually specify the binary:
```javascript
const requester = create({
  binaryPath: './bin/fingerprint_custom'
});
```

## Examples

### With Custom Headers
```javascript
const response = await requester.get('https://api.example.com', {
  headers: {
    'User-Agent': 'MyApp/1.0',
    'Authorization': 'Bearer token123'
  }
});
```

### With Timeout
```javascript
const response = await requester.get('https://api.example.com', {
  timeout: 5000 // 5 seconds
});
```

### JSON Response
```javascript
const response = await requester.get('https://api.example.com/json', {
  responseType: 'json'
});
console.log(response.data.field);
```

### Error Handling
```javascript
try {
  const response = await requester.get('https://api.example.com');
  console.log(response.data);
} catch (error) {
  console.error('Request failed:', error.message);
}
```

## Run Tests

```bash
npm test
```
