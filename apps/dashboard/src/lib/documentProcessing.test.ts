// Mock makeRequest
jest.mock('./makeRequest', () => ({
  makeRequest: jest.fn(),
}));

const loadModuleUnderTest = async () => {
  process.env.NEXT_PUBLIC_API_URL = 'http://api.backend.test';
  process.env.NEXT_PUBLIC_API_URL_FRONT = 'http://api.frontend.test';
  process.env.NEXT_PUBLIC_API_AUTH_TOKEN = 'test-token';

  jest.resetModules();

  const documentProcessing = await import('./documentProcessing');
  const makeRequestModule = await import('./makeRequest');

  return {
    documentProcessing,
    makeRequest: makeRequestModule.makeRequest as jest.Mock,
  };
};

describe('documentProcessing', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('getAllProcess', () => {
    it('should call backend API URL by default', async () => {
      const { documentProcessing, makeRequest } = await loadModuleUnderTest();
      makeRequest.mockResolvedValue({ processes: [] });

      await documentProcessing.getAllProcess();

      expect(makeRequest).toHaveBeenCalledWith('http://api.backend.test/process/list', 'test-token');
    });

    it('should call frontend API URL when stage is front', async () => {
      const { documentProcessing, makeRequest } = await loadModuleUnderTest();
      makeRequest.mockResolvedValue({ processes: [] });

      await documentProcessing.getAllProcess('front');

      expect(makeRequest).toHaveBeenCalledWith('http://api.frontend.test/process/list', 'test-token');
    });
  });

  describe('startProcess', () => {
    it('should start process with POST method on backend', async () => {
      const { documentProcessing, makeRequest } = await loadModuleUnderTest();
      makeRequest.mockResolvedValue({ id: 'proc-1' });

      await documentProcessing.startProcess();

      expect(makeRequest).toHaveBeenCalledWith(
        'http://api.backend.test/process/start',
        'test-token',
        'POST'
      );
    });

    it('should start process on frontend when stage is front', async () => {
      const { documentProcessing, makeRequest } = await loadModuleUnderTest();
      makeRequest.mockResolvedValue({ id: 'proc-2' });

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
      const { documentProcessing, makeRequest } = await loadModuleUnderTest();
      makeRequest.mockResolvedValue({ stopped: true });

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
      const { documentProcessing, makeRequest } = await loadModuleUnderTest();
      const mockStatus = { id: 'proc-1', status: 'running' };
      makeRequest.mockResolvedValue(mockStatus);

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
      const { documentProcessing, makeRequest } = await loadModuleUnderTest();
      const mockResults = { summary: 'test summary', wordCount: 1000 };
      makeRequest.mockResolvedValue(mockResults);

      const result = await documentProcessing.resultProcess('proc-1');

      expect(makeRequest).toHaveBeenCalledWith(
        'http://api.backend.test/process/results/proc-1',
        'test-token'
      );
      expect(result).toEqual(mockResults);
    });
  });
});
