import { createContext } from "react";
import Service from "./types/service.tsx";
import { CookieKeys, Timescale } from "./types/strings";

export interface ContextType {
	services: Service[];
	timescale: Timescale;
	setTimescale: (timescale: Timescale) => void;
	serviceAdd: (service: Service) => void;
	serviceDelete: (serviceID: string) => void;
	serviceUpdate: (service: Service) => void;
	requestServiceData: () => void;
	cookieGet: (key: CookieKeys) => string | undefined;
	passwordReset: (newPassword: string) => void;
	serviceToggle: (serviceID: string) => void;
}

// Create the context with a default value
export const Context = createContext<ContextType | undefined>(undefined);
