import { createContext } from "react";

export type CookieKeys = "session-token";

export interface ContextType {
	userRequestData: () => void;
	cookieGet: (key: CookieKeys) => string | undefined;
	passwordReset: (newPassword: string) => void;
}

// Create the context with a default value
export const Context = createContext<ContextType | undefined>(undefined);
