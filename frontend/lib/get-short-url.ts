import { BASE_URL } from "./utils";

export default async function getShortUrl(url: string): Promise<string> {
  const apiResponse = await fetch(`${BASE_URL}/api/shorten`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-API-Key": "somekey",
    },
    body: JSON.stringify({ original: url }),
  });

  if (!apiResponse.ok) {
    const errorMessage = await apiResponse.json();

    //filter known api error codes
    if ([400, 401, 404, 429, 500].includes(apiResponse.status))
      throw new Error(errorMessage);
    else {
      process.env.NODE_ENV != "production" &&
        console.log(
          "Request to shorten " + url + " failed. API response: " + apiResponse,
        );

      throw new Error(
        "Oops, something broke on our end. Please try again later",
      );
    }
  }

  const shortUrl = await apiResponse.json().then((res) => res.shortCode);
  return shortUrl;
}
