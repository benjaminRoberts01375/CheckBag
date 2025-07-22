import { createContext } from "react";
import Service from "./types/service.tsx";

export type CookieKeys = "session-token";

export interface ContextType {
	services: Service[];
	serviceAdd: (service: Service) => void;
	requestServiceData: () => void;
	cookieGet: (key: CookieKeys) => string | undefined;
	passwordReset: (newPassword: string) => void;
	serviceToggle: (serviceID: string) => void;
}

// Create the context with a default value
export const Context = createContext<ContextType | undefined>(undefined);
