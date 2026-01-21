import { spawn } from 'child_process';
import { existsSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { platform, arch } from 'os';

const __dirname = dirname(fileURLToPath(import.meta.url));

class FingerprintRequester {
  constructor(options = {}) {
    this.binDir = options.binDir || join(__dirname, 'bin');
    this.binaryPath = options.binaryPath || this._detectBinary();
    this.configPath = options.configPath || join(__dirname, 'config.json');
    this.defaults = {
      timeout: options.timeout || 30, // seconds
      proxy: options.proxy || null,
    };
    this.activeProcesses = new Set();
  }

  _detectBinary() {
    const platformMap = {
      win32: 'windows',
      linux: 'linux',
      android: 'android',
      darwin: 'linux', // fallback to linux for macOS
    };

    const archMap = {
      x64: 'amd64',
      arm64: 'arm64',
    };

    const os = platformMap[platform()];
    const cpuArch = archMap[arch()];

    if (!os || !cpuArch) {
      throw new Error(`Unsupported platform: ${platform()} ${arch()}`);
    }

    const ext = platform() === 'win32' ? '.exe' : '';
    const binaryName = `fingerprint_${os}_${cpuArch}${ext}`;
    const binaryPath = join(this.binDir, binaryName);

    if (!existsSync(binaryPath)) {
      throw new Error(`Binary not found: ${binaryPath}`);
    }

    return binaryPath;
  }

  async request(config) {
    const {
      method = 'GET',
      url,
      headers = {},
      data = '',
      timeout,
      proxy,
      responseType = 'text',
      onDownloadProgress,
      validateStatus = (status) => status >= 200 && status < 300,
      signal,
    } = config;

    if (!url) {
      throw new Error('URL is required');
    }

    const requestPayload = {
      method: method.toUpperCase(),
      url,
      headers,
      body: typeof data === 'string' ? data : JSON.stringify(data),
      config_path: this.configPath,
    };

    // Add timeout if specified (in seconds)
    const timeoutSec = timeout || this.defaults.timeout;
    if (timeoutSec) {
      requestPayload.timeout = {
        connect: timeoutSec,
        read: timeoutSec,
      };
    }

    // Add proxy if specified
    const proxyUrl = proxy || this.defaults.proxy;
    if (proxyUrl) {
      const proxyType = proxyUrl.startsWith('socks') ? 'socks5' : 'http';
      requestPayload.proxy = {
        enabled: true,
        type: proxyType,
        url: proxyUrl,
      };
    }

    return new Promise((resolve, reject) => {
      const proc = spawn(this.binaryPath);
      this.activeProcesses.add(proc);
      let headersParsed = false;
      let responseHeaders = {};
      let responseStatus = 200;
      let responseStatusText = 'OK';
      let headerBuffer = '';
      let bodyChunks = [];
      let totalLoaded = 0;
      let stderrData = '';

      const timeoutId = setTimeout(() => {
        proc.kill();
        const error = new Error('Request timeout');
        error.code = 'ECONNABORTED';
        error.config = config;
        reject(error);
      }, timeoutSec * 1000);

      // Support request cancellation
      if (signal) {
        signal.addEventListener('abort', () => {
          proc.kill();
          clearTimeout(timeoutId);
          const error = new Error('Request aborted');
          error.code = 'ERR_CANCELED';
          error.config = config;
          reject(error);
        });
      }

      proc.stdout.on('data', (chunk) => {
        if (!headersParsed) {
          headerBuffer += chunk.toString();
          const headerEndIndex = headerBuffer.indexOf('\r\n\r\n');
          
          if (headerEndIndex !== -1) {
            // Parse headers
            const headerPart = headerBuffer.substring(0, headerEndIndex);
            const bodyPart = headerBuffer.substring(headerEndIndex + 4);
            
            const lines = headerPart.split('\r\n');
            const statusMatch = lines[0].match(/HTTP\/[\d.]+ (\d+) (.+)/);
            responseStatus = statusMatch ? parseInt(statusMatch[1]) : 200;
            responseStatusText = statusMatch ? statusMatch[2] : 'OK';
            
            for (let i = 1; i < lines.length; i++) {
              const [key, ...valueParts] = lines[i].split(': ');
              if (key) responseHeaders[key.toLowerCase()] = valueParts.join(': ');
            }
            
            headersParsed = true;
            
            // Clear timeout for streaming responses
            clearTimeout(timeoutId);
            
            // 先调用 validateStatus 让外部知道状态码
            if (validateStatus) {
              validateStatus(responseStatus);
            }
            
            // Process body part after headers
            if (bodyPart) {
              bodyChunks.push(Buffer.from(bodyPart));
              totalLoaded += bodyPart.length;
              if (onDownloadProgress) {
                onDownloadProgress({
                  loaded: totalLoaded,
                  total: parseInt(responseHeaders['content-length']) || 0,
                  chunk: bodyPart,
                  status: responseStatus,
                });
              }
            }
          }
        } else {
          // Headers already parsed, process body chunks
          const chunkStr = chunk.toString();
          bodyChunks.push(chunk);
          totalLoaded += chunk.length;
          if (onDownloadProgress) {
            onDownloadProgress({
              loaded: totalLoaded,
              total: parseInt(responseHeaders['content-length']) || 0,
              chunk: chunkStr,
              status: responseStatus,
            });
          }
        }
      });

      proc.stderr.on('data', (chunk) => {
        stderrData += chunk.toString();
      });

      proc.on('close', (code) => {
        clearTimeout(timeoutId);
        this.activeProcesses.delete(proc);

        if (code !== 0) {
          let errorInfo = { error: `Process exited with code ${code}`, error_type: 'UNKNOWN_ERROR' };
          if (stderrData) {
            try {
              errorInfo = JSON.parse(stderrData);
            } catch (e) {
              errorInfo.error = stderrData;
            }
          }
          const error = new Error(errorInfo.error);
          error.code = code === 3 ? 'ECONNABORTED' : code === 4 ? 'ERR_CONFIG' : 'ERR_NETWORK';
          error.errorType = errorInfo.error_type;
          error.exitCode = code;
          error.config = config;
          return reject(error);
        }

        const body = Buffer.concat(bodyChunks).toString();
        let parsedData = body;
        
        if (responseType === 'json') {
          try {
            parsedData = JSON.parse(body);
          } catch (e) {
            // keep as text if parse fails
          }
        }

        const response = {
          data: parsedData,
          status: responseStatus,
          statusText: responseStatusText,
          headers: responseHeaders,
          config,
        };
        
        if (!validateStatus(responseStatus)) {
          const error = new Error(`Request failed with status code ${responseStatus}`);
          error.code = responseStatus >= 500 ? 'ERR_BAD_RESPONSE' : 'ERR_BAD_REQUEST';
          error.response = response;
          error.config = config;
          return reject(error);
        }
        
        resolve(response);
      });

      proc.on('error', (err) => {
        clearTimeout(timeoutId);
        const error = new Error(`Failed to spawn process: ${err.message}`);
        error.code = 'ERR_SPAWN';
        error.config = config;
        reject(error);
      });

      proc.stdin.write(JSON.stringify(requestPayload));
      proc.stdin.end();
    });
  }

  async get(url, config = {}) {
    return this.request({ ...config, method: 'GET', url });
  }

  async post(url, data, config = {}) {
    return this.request({ ...config, method: 'POST', url, data });
  }

  async put(url, data, config = {}) {
    return this.request({ ...config, method: 'PUT', url, data });
  }

  async delete(url, config = {}) {
    return this.request({ ...config, method: 'DELETE', url });
  }

  async patch(url, data, config = {}) {
    return this.request({ ...config, method: 'PATCH', url, data });
  }

  async head(url, config = {}) {
    return this.request({ ...config, method: 'HEAD', url });
  }

  async options(url, config = {}) {
    return this.request({ ...config, method: 'OPTIONS', url });
  }

  stream(config) {
    // Ensure onDownloadProgress is set for streaming
    const onProgress = config.onData || config.onDownloadProgress;
    if (!onProgress) {
      console.warn('[stream] No onData or onDownloadProgress callback provided');
    }
    return this.request({
      ...config,
      onDownloadProgress: onProgress || (() => {}),
    });
  }

  close() {
    this.activeProcesses.forEach(proc => proc.kill());
    this.activeProcesses.clear();
  }
}

export function create(options) {
  return new FingerprintRequester(options);
}

export default {
  create,
  FingerprintRequester,
};
