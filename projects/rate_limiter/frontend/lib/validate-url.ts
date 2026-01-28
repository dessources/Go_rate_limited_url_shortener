import { ALLOWED_PROTOCOLS, MAX_URL_LENGTH } from "./utils";

export default function validateURL(url: string) {
  let parsedURL;
  try {
    parsedURL = new URL(url);
  } catch {
    throw new Error(
      "Invalid URL. Please provide a valid link e.g. https://example.com/very/long/url/...",
    );
  }

  if (url.length > MAX_URL_LENGTH) {
    throw new Error(
      "Your link is too long. Max length is " + MAX_URL_LENGTH + " characters.",
    );
  }

  if (!ALLOWED_PROTOCOLS.includes(parsedURL.protocol)) {
    throw new Error(
      "Your link uses an invalid protocol. Please provide a link starting with " +
        ALLOWED_PROTOCOLS.split(",").join(":// ") +
        ":// .",
    );
  }
}
