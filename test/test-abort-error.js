import { create } from '../requester.js';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));

const requester = create({
  configPath: join(__dirname, '../config-chrome141.json'),
  timeout: 60000,
});

// 测试取消请求
async function testAbort() {
  console.log('=== 测试取消请求 ===\n');
  
  const controller = new AbortController();
  
  // 1秒后取消请求
  setTimeout(() => {
    console.log('取消请求...');
    controller.abort();
  }, 1000);
  
  try {
    const response = await requester.get('https://httpbin.org/delay/5', {
      signal: controller.signal,
    });
    console.log('请求成功:', response.status);
  } catch (error) {
    console.log('错误码:', error.code); // ERR_CANCELED
    console.log('错误信息:', error.message); // Request aborted
  }
}

// 测试超时
async function testTimeout() {
  console.log('\n=== 测试超时 ===\n');
  
  try {
    const response = await requester.get('https://httpbin.org/delay/10', {
      timeout: 2000, // 2秒超时
    });
    console.log('请求成功:', response.status);
  } catch (error) {
    console.log('错误码:', error.code); // ECONNABORTED
    console.log('错误信息:', error.message); // Request timeout
  }
}

// 测试 4xx 错误
async function test4xxError() {
  console.log('\n=== 测试 4xx 错误 ===\n');
  
  try {
    const response = await requester.get('https://httpbin.org/status/404');
    console.log('请求成功:', response.status);
  } catch (error) {
    console.log('错误码:', error.code); // ERR_BAD_REQUEST
    console.log('错误信息:', error.message); // Request failed with status code 404
    console.log('响应状态:', error.response.status); // 404
  }
}

// 测试 5xx 错误
async function test5xxError() {
  console.log('\n=== 测试 5xx 错误 ===\n');
  
  try {
    const response = await requester.get('https://httpbin.org/status/500');
    console.log('请求成功:', response.status);
  } catch (error) {
    console.log('错误码:', error.code); // ERR_BAD_RESPONSE
    console.log('错误信息:', error.message); // Request failed with status code 500
    console.log('响应状态:', error.response.status); // 500
  }
}

// 测试自定义 validateStatus
async function testCustomValidateStatus() {
  console.log('\n=== 测试自定义 validateStatus ===\n');
  
  try {
    // 允许 404 状态码
    const response = await requester.get('https://httpbin.org/status/404', {
      validateStatus: (status) => status === 404 || (status >= 200 && status < 300),
    });
    console.log('请求成功:', response.status); // 404
    console.log('允许 404 通过验证');
  } catch (error) {
    console.log('错误:', error.message);
  }
}

// 运行测试
(async () => {
  await testAbort();
  await testTimeout();
  await test4xxError();
  await test5xxError();
  await testCustomValidateStatus();
})();
