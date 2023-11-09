import { API_URL } from './config';
import { BadRequestResponse } from './types';

export const sendRequest = async (
  token: string | null,
  method: 'POST' | 'GET' | 'PUT' | 'DELETE',
  endpoint: string,
  body: any = null
) => {
  const headers = new Headers();
  headers.append('Content-Type', 'application/json');
  if (token != null) {
    headers.append('Authorization', 'Basic ' + token);
  }

  const requestOptions =
    body === null ? { method, headers } : { method, headers, body: JSON.stringify(body) };

  const response = await fetch(API_URL + endpoint, requestOptions);
  if (response.status === 500) {
    alert('An internal error has occurred');
  } else if (response.status === 400) {
    const body: BadRequestResponse = await response.json();
    alert(`Bad request: ${body.message}`);
  }
  return response;
};

export const isSuccessResponse = (res: Response) => {
  return res.status < 400;
};
