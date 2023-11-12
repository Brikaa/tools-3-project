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

  const response = await fetch('/api/' + endpoint, requestOptions);
  if (response.status === 500) {
    alert('An internal error has occurred');
  } else if (response.status === 400) {
    const body = await response.text();
    alert(`Bad request: ${body}`);
  }
  return response;
};

export const isSuccessResponse = (res: Response) => {
  return res.status < 400;
};
