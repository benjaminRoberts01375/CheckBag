import React, { ReactNode } from "react";
import { Context, ContextType, CookieKeys } from "./context-object";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import User from "./types/user";

interface Props {
	children: ReactNode;
}

export const ContextProvider: React.FC<Props> = ({ children }) => {
	const [user, setUser] = useState<User | undefined>(undefined);
	const navigate = useNavigate();

	function cookieGet(key: CookieKeys): string | undefined {
		const cookieString = document.cookie.split("; ").find(cookie => cookie.startsWith(`${key}=`));

		if (cookieString) {
			return decodeURIComponent(cookieString.split("=")[1]);
		}
		return undefined;
	}

	function userSignUp(
		username: string,
		password: string,
		first_name: string,
		last_name: string,
	): void {
		(async () => {
			try {
				const response = await fetch("/api/user-sign-up", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify({
						username: username,
						password: password,
						first_name: first_name,
						last_name: last_name,
					}),
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to sign up user");
				}
				const rawData = await response.json();
				setUser(rawData.user);
				console.log("Successfully signed up user");
				navigate("/check-email");
			} catch (error) {
				console.error("Error signing up user:", error);
			}
		})();
	}

	function userLogin(username: string, password: string): void {
		(async () => {
			try {
				const response = await fetch("/api/user-sign-in", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify({
						username: username,
						password: password,
					}),
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to log in user");
				}
				const rawData = await response.json();
				setUser(rawData.user);
				console.log("Successfully logged in user");
			} catch (error) {
				console.error("Error logging in user:", error);
			}
		})();
	}

	function userLoginJWT(): void {
		(async () => {
			try {
				const response = await fetch("/api/user-sign-in-jwt", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to log in user");
				}
				const rawData = await response.json();
				setUser(rawData.user);
				console.log("Successfully logged in user");
			} catch (error) {
				console.error("Error logging in user:", error);
			}
		})();
	}

	function userLogout(): void {
		(async () => {
			try {
				await fetch("/api/user-logout", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					credentials: "include",
				});

				// Delete the cookie after the user logs out
				document.cookie =
					"session-token=; Max-Age=0; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
				setUser(undefined);
			} catch (error) {
				console.error("Error deleting gift:", error);
			}
		})();
		navigate("/login");
	}

	function userRequestData(): void {
		(async () => {
			try {
				const response = await fetch("/api/user-get-data", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to fetch user data");
				}
				const rawData = await response.json();
				setUser(rawData.user);
				console.log("Successfully fetched user data:");
			} catch (error) {
				console.error("Error fetching user data:", error);
			}
		})();
	}

	function passwordReset(newPassword: string): void {
		(async () => {
			try {
				const response = await fetch("/api/user-reset-password", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify(newPassword),
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to reset password");
				}
				console.log("Successfully reset password");
			} catch (error) {
				console.error("Error resetting password:", error);
			}
		})();
	}

	const contextShape: ContextType = {
		user,
		cookieGet,
		userSignUp,
		userLogin,
		userLoginJWT,
		userLogout,
		userRequestData,
		passwordReset,
	};

	return <Context.Provider value={contextShape}>{children}</Context.Provider>;
};
