export const makeRequest = async (API_URL: string, API_KEY: string, method: string = "GET", body?: any) => {
  const response = await fetch(`${API_URL}`, {
    method,
    headers: {
      Authorization: `Bearer ${API_KEY}`,
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : undefined,
  });
  return response.json();
};
