import Blowfish from "egoroof-blowfish";

const bf = new Blowfish(
  "thisissecretkey",
  Blowfish.MODE.ECB,
  Blowfish.PADDING.NULL
);

const parseSubdomain = (data: Buffer) => {
  const httpRequest = data.toString();
  const lines = httpRequest.split("\n").map((x) => x.trim());
  const headers: { [_: string]: string } = lines
    .slice(1)
    .reduce((acc, line) => {
      const [key, value] = line.split(": ");
      acc[key] = value;
      return acc;
    }, {} as { [_: string]: string });
  if (!headers.Host || !headers.Host.includes(".")) {
    return { subdomain: "", err: true };
  }
  const subdomain = headers.Host.split(".")[0] as string;
  return { subdomain, err: false };
};

export const parseAppAuthority = (data: Buffer) => {
  let { subdomain, err } = parseSubdomain(data);
  if (err) {
    return { appHost: "", appPort: 80, err: true };
  }
  let decodedSubdomain = Buffer.from(subdomain, "hex");
  const decipheredSubdomain = bf.decode(decodedSubdomain, Blowfish.TYPE.STRING);
  const appAuthority = decipheredSubdomain.split(":");
  const appHost = appAuthority[0];
  const appPort = parseInt(appAuthority[1]);
  err = false;
  return { appHost, appPort, err };
};
