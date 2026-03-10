import { makeRequest } from "./makeRequest";

const API_KEY: string = process.env.NEXT_PUBLIC_API_AUTH_TOKEN || "";
const API_URL: string = process.env.NEXT_PUBLIC_API_URL || "";
const API_URL_FRONT: string = process.env.NEXT_PUBLIC_API_URL_FRONT || "";

// Type for process status
export interface ProcessStatus {
  id: string;
  status: string;
  error: string;
  started_at: string;
  estimated_completion: string;
  files_processed: string[];
  files_to_process: string[];
  completed_at: string;
}

// REST API version (original)
export const getAllProcess = async (stage: string = "back") => {
  const api_url = stage === "front" ? API_URL_FRONT : API_URL;
  return makeRequest(`${api_url}/process/list`, API_KEY);
};

export const startProcess = async (stage: string = "back") => {
  const api_url = stage === "front" ? API_URL_FRONT : API_URL;
  return makeRequest(`${api_url}/process/start`, API_KEY, "POST");
}

export const stopProcess = async (processId: string, stage: string = "back") => {
  const api_url = stage === "front" ? API_URL_FRONT : API_URL;
  return makeRequest(`${api_url}/process/stop/${processId}`, API_KEY, "POST");
}

export const statusProcess = async (processId: string, stage: string = "back") => {
  const api_url = stage === "front" ? API_URL_FRONT : API_URL;
  return makeRequest(`${api_url}/process/status/${processId}`, API_KEY);
}

export const resultProcess = async (processId: string, stage: string = "back") => {
  const api_url = stage === "front" ? API_URL_FRONT : API_URL;
  return makeRequest(`${api_url}/process/results/${processId}`, API_KEY);
}
