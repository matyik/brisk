console.log('Booting API...');

const PORT = 3000;

Brisk.serve(PORT, (req) => {
  console.info(`[${req.method}] ${req.url}`);

  if (req.url === '/health') {
    return {
      status: 200,
      body: JSON.stringify({ status: 'healthy', engine: 'brisk' }),
    };
  }

  if (req.url === '/echo' && req.method === 'POST') {
    return {
      status: 201,
      body: `You sent: ${req.body}`,
    };
  }

  return {
    status: 404,
    body: '404 - Route not found',
  };
});
