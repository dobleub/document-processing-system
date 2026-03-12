import dotenv from 'dotenv';
import path from 'path';

// Load .env.test with override: true so test values take precedence over
// .env.local / .env that Next.js's jest preset auto-loads before this runs.
dotenv.config({ path: path.resolve(__dirname, '.env.test'), override: true });

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter() {
    return {
      push: jest.fn(),
      back: jest.fn(),
      forward: jest.fn(),
    };
  },
  usePathname() {
    return '/';
  },
  useSearchParams() {
    return new URLSearchParams();
  },
}));
