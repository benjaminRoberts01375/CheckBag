export const Timescales = ["hour", "day", "month", "year"] as const;
export type Timescale = (typeof Timescales)[number];
export type CookieKeys = "session-token";
export type CommunicationProtocols = "http" | "https";
