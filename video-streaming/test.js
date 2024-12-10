const http3 = require('http3');

http3
  .request('https://localhost:8080')
  .on('response', (response) => {
    console.log(`Status: ${response.statusCode}`);
    response.pipe(process.stdout);
  })
  .end();

