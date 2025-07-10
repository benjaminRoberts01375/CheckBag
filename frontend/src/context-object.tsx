import { createContext } from "react";
import User from "./types/user";

export type CookieKeys = "session-token";

export interface ContextType {
	user: User | undefined;
	userSignUp: (username: string, password: string, first_name: string, last_name: string) => void;
	userLogin: (username: string, password: string) => void;
	userLoginJWT: () => void;
	userLogout: () => void;
	userRequestData: () => void;
	cookieGet: (key: CookieKeys) => string | undefined;
	passwordReset: (newPassword: string) => void;
	emailResetRequest: (newEmail: string) => void;
	forgotPasswordRequest: (email: string) => void;
	forgotPasswordCheckValid: (token: string) => void;
	forgotPasswordConfirmation: (token: string) => void;
}

// Create the context with a default value
export const Context = createContext<ContextType | undefined>(undefined);
