import { spawn } from 'node:child_process';
import process from 'node:process';

const child = spawn('go', ['run', '.'], {
  stdio: 'inherit',
  shell: process.platform === 'win32',
  env: {
    ...process.env,
    PORT: process.env.PORT || '3000',
    JWT_SECRET: process.env.JWT_SECRET || 'local-development-secret',
    FRONTEND_URL: process.env.FRONTEND_URL || 'http://localhost:3002',
    RATE_LIMIT_MAX: process.env.RATE_LIMIT_MAX || '100',
    RATE_LIMIT_WINDOW_SECONDS: process.env.RATE_LIMIT_WINDOW_SECONDS || '60',
    NODE_API_URL: process.env.NODE_API_URL || 'http://localhost:3001',
  },
});

child.on('exit', (code) => {
  process.exit(code ?? 0);
});
