import * as documentProcessing from './documentProcessing';

// Mock makeRequest
jest.mock('./makeRequest', () => ({
  makeRequest: jest.fn(),
}));

import { makeRequest } from './makeRequest';

describe('documentProcessing', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    process.env.NEXT_PUBLIC_API_URL = 'http://api.backend.test';
    process.env.NEXT_PUBLIC_API_URL_FRONT = 'http://api.frontend.test';
    process.env.NEXT_PUBLIC_API_AUTH_TOKEN = 'test-token';
  });

  describe('getAllProcess', () => {
    it('should call backend API URL by default', async () => {
      (makeRequest as jest.Mock).mockResolvedValue({ processes: [] });

      await documentProcessing.getAllProcess();

      expect(makeRequest).toHaveBeenCalledWith('http://api.backend.test/process/list', 'test-token');
    });

    it('should call frontend API URL when stage is front', async () => {
      (makeRequest as jest.Mock).mockResolvedValue({ processes: [] });

      await documentProcessing.getAllProcess('front');

      expect(makeRequest).toHaveBeenCalledWith('http://api.frontend.test/process/list', 'test-token');
    });
  });

  describe('startProcess', () => {
    it('should start process with POST method on backend', async () => {
      (makeRequest as jest.Mock).mockResolvedValue({ id: 'proc-1' });

      await documentProcessing.startProcess();

      expect(makeRequest).toHaveBeenCalledWith(
        'http://api.backend.test/process/start',
        'test-token',
        'POST'
      );
    });

    it('should start process on frontend when stage is front', async () => {
      (makeRequest as jest.Mock).mockResolvedValue({ id: 'proc-2' });

      await documentProcessing.startProcess('front');

      expect(makeRequest).toHaveBeenCalledWith(
        'http://api.frontend.test/process/start',
        'test-token',
        'POST'
      );
    });
  });

  describe('stopProcess', () => {
    it('should stop process with processId', async () => {
      (makeRequest as jest.Mock).mockResolvedValue({ stopped: true });

      await documentProcessing.stopProcess('proc-123');

      expect(makeRequest).toHaveBeenCalledWith(
        'http://api.backend.test/process/stop/proc-123',
        'test-token',
        'POST'
      );
    });
  });

  describe('statusProcess', () => {
    it('should fetch status for processId', async () => {
      const mockStatus = { id: 'proc-1', status: 'running' };
      (makeRequest as jest.Mock).mockResolvedValue(mockStatus);

      const result = await documentProcessing.statusProcess('proc-1');

      expect(makeRequest).toHaveBeenCalledWith(
        'http://api.backend.test/process/status/proc-1',
        'test-token'
      );
      expect(result).toEqual(mockStatus);
    });
  });

  describe('resultProcess', () => {
    it('should fetch results for processId', async () => {
      const mockResults = { summary: 'test summary', wordCount: 1000 };
      (makeRequest as jest.Mock).mockResolvedValue(mockResults);

      const result = await documentProcessing.resultProcess('proc-1');

      expect(makeRequest).toHaveBeenCalledWith(
        'http://api.backend.test/process/results/proc-1',
        'test-token'
      );
      expect(result).toEqual(mockResults);
    });
  });

  describe('ProcessStatus interface', () => {
    it('should define ProcessStatus type with required fields', () => {
      const status: documentProcessing.ProcessStatus = {
        id: '123',
        status: 'running',
        error: '',
        started_at: '2024-01-01T10:00:00Z',
        estimated_completion: '2024-01-01T11:00:00Z',
        files_processed: ['file1.txt'],
        files_to_process: ['file2.txt'],
        completed_at: '',
      };

      expect(status.id).toBe('123');
      expect(status.status).toBe('running');
      expect(status.files_processed).toContain('file1.txt');
    });
  });
});
