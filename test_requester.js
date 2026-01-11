import { create } from './requester.js';

// Create requester instance
const requester = create({
  // binDir: './bin',  // default
  // binaryPath: './bin/custom_binary',  // manual override
  configPath: './config.json',  // default
  timeout: 30000,
});

async function testBasicRequest() {
  console.log('=== Basic GET Request ===');
  try {
    const response = await requester.get('https://tls.peet.ws/api/all', {
      responseType: 'json',
    });
    console.log('Status:', response.status);
    console.log(JSON.stringify(response.data,null,2))
  } catch (error) {
    console.error('Error:', error.message);
  }
}

async function testPostRequest() {
  console.log('\n=== POST Request ===');
  try {
    const response = await requester.post(
      'https://httpbin.org/post',
      { test: 'data', foo: 'bar' },
      { responseType: 'json' }
    );
    console.log('Status:', response.status);
    console.log('Response data:', response.data.json);
  } catch (error) {
    console.error('Error:', error.message);
  }
}

async function testStreamRequest() {
  console.log('\n=== Stream Request ===');
  try {
    let chunks = 0;
    await requester.stream({
      method: 'GET',
      url: 'https://httpbin.org/stream/5',
      onData: ({ data }) => {
        chunks++;
        process.stdout.write('.');
      },
    });
    console.log(`\nReceived ${chunks} chunks`);
  } catch (error) {
    console.error('Error:', error.message);
  }
}

async function testWithHeaders() {
  console.log('\n=== Request with Custom Headers ===');
  try {
    const response = await requester.get('https://httpbin.org/headers', {
      headers: {
        'User-Agent': 'CustomAgent/1.0',
        'X-Custom-Header': 'test-value',
      },
      responseType: 'json',
    });
    console.log('Status:', response.status);
    console.log('Request headers seen by server:', response.data.headers);
  } catch (error) {
    console.error('Error:', error.message);
  }
}

async function testTimeout() {
  console.log('\n=== Timeout Test ===');
  try {
    await requester.get('https://httpbin.org/delay/10', {
      timeout: 2000, // 2 seconds
    });
  } catch (error) {
    console.log('Expected timeout error:', error.message);
  }
}

// Run all tests
(async () => {
  await testBasicRequest();
  await testPostRequest();
  await testStreamRequest();
  await testWithHeaders();
  await testTimeout();
  
  console.log('\n=== All tests completed ===');
})();
