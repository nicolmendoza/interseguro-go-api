import { spawn } from 'node:child_process';
import process from 'node:process';

const child = spawn('go', ['run', '.'], {
  stdio: 'inherit',
  shell: process.platform === 'win32',
  env: {
    ...process.env,
    PORT: process.env.PORT || '3000',
    JWT_SECRET: process.env.JWT_SECRET || 'local-development-secret',
    NODE_API_URL: process.env.NODE_API_URL || 'http://localhost:3001',
  },
});

child.on('exit', (code) => {
  process.exit(code ?? 0);
});
