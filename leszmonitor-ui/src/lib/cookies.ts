export const setCookie = (name: string, value: string, options?: { maxAge?: number; path?: string; secure?: boolean; sameSite?: string }): void => {
  let cookieString = `${encodeURIComponent(name)}=${encodeURIComponent(value)}`;

  if (options?.maxAge) {
    cookieString += `; Max-Age=${options.maxAge}`;
  }

  if (options?.path) {
    cookieString += `; Path=${options.path}`;
  }

  if (options?.secure) {
    cookieString += "; Secure";
  }

  if (options?.sameSite) {
    cookieString += `; SameSite=${options.sameSite}`;
  }

  document.cookie = cookieString;
};

export const getCookie = (name: string): string | null => {
  const cookies = document.cookie.split("; ");
  for (const cookie of cookies) {
    const [cookieName, cookieValue] = cookie.split("=");
    if (decodeURIComponent(cookieName) === name) {
      return decodeURIComponent(cookieValue);
    }
  }
  return null;
};

export const deleteCookie = (name: string): void => {
  setCookie(name, "", { maxAge: -1, path: "/" });
};

