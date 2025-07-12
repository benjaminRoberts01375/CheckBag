import React, { ReactNode } from "react";
import { Context, ContextType, CookieKeys } from "./context-object";
import { useState } from "react";
import Service from "./types/service.tsx";

interface Props {
	children: ReactNode;
}

export const ContextProvider: React.FC<Props> = ({ children }) => {

	function cookieGet(key: CookieKeys): string | undefined {
		const cookieString = document.cookie.split("; ").find(cookie => cookie.startsWith(`${key}=`));

		if (cookieString) {
			return decodeURIComponent(cookieString.split("=")[1]);
		}
		return undefined;
	}

	function requestInitialData(): void {
		(async () => {
			try {

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
				setServices(rawData.services);
			} catch (error) {
				console.error("Error fetching initial data:", error);
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
		userLoginJWT,
		userLogout,
		userRequestData,
		passwordReset,
	};

	return <Context.Provider value={contextShape}>{children}</Context.Provider>;
};
