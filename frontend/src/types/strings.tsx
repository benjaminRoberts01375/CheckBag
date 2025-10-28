export const Timescales = ["hour", "day", "month", "year"] as const;
export type Timescale = (typeof Timescales)[number];
export type CookieKeys = "checkbag-session-token";
export const CommunicationProtocols = ["http", "https"];
export type CommunicationProtocol = (typeof CommunicationProtocols)[number];
