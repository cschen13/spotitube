const message =
  "An error occurred while fetching data from the server. Please try again later.";

export interface IApiResponse<T> {
  status: number;
  value?: T;
  error?: {
    message: string;
  };
}

async function parse<T>(response: Response): Promise<IApiResponse<T>> {
  const result = { status: response.status } as IApiResponse<T>;

  try {
    if (response.ok) {
      const data = await response.json();
      result.value = data;
    } else {
      const errorText = await response.text();
      console.error(errorText);
      result.error = { message: errorText };
    }
  } catch (err) {
    console.error(err);
    result.error = { message };
  }

  return result;
}

export default async function request<T>(
  url: string,
  includeCredentials: boolean = false,
  method: string = "GET",
  body?: any,
  headers?: any
): Promise<IApiResponse<T>> {
  const options: RequestInit = {
    body,
    headers: new Headers({ "content-type": "application/json", ...headers }),
    method
  };

  if (includeCredentials) {
    options.credentials = "include";
  }

  try {
    const response = await fetch(url, options);
    return parse<T>(response);
  } catch (err) {
    console.error(err);
    return {
      error: { message }
    } as IApiResponse<T>;
  }
}
