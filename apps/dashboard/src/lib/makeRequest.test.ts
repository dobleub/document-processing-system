import { makeRequest } from './makeRequest';

describe('makeRequest', () => {
  beforeEach(() => {
    global.fetch = jest.fn();
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it('should make a GET request with authorization header', async () => {
    const mockResponse = { success: true, data: [] };
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      json: async () => mockResponse,
    });

    const result = await makeRequest('http://api.test/list', 'secret-token');

    expect(global.fetch).toHaveBeenCalledWith('http://api.test/list', {
      method: 'GET',
      headers: {
        Authorization: 'Bearer secret-token',
        'Content-Type': 'application/json',
      },
      body: undefined,
    });
    expect(result).toEqual(mockResponse);
  });

  it('should make a POST request with body', async () => {
    const mockResponse = { id: '123', status: 'pending' };
    const payload = { name: 'test' };
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      json: async () => mockResponse,
    });

    const result = await makeRequest('http://api.test/create', 'secret-token', 'POST', payload);

    expect(global.fetch).toHaveBeenCalledWith('http://api.test/create', {
      method: 'POST',
      headers: {
        Authorization: 'Bearer secret-token',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });
    expect(result).toEqual(mockResponse);
  });

  it('should handle fetch errors gracefully', async () => {
    (global.fetch as jest.Mock).mockRejectedValueOnce(new Error('Network error'));

    await expect(
      makeRequest('http://api.test/list', 'secret-token')
    ).rejects.toThrow('Network error');
  });

  it('should support Custom HTTP methods', async () => {
    const mockResponse = { deleted: true };
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      json: async () => mockResponse,
    });

    await makeRequest('http://api.test/process/123', 'secret-token', 'DELETE');

    expect(global.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({ method: 'DELETE' })
    );
  });
});
