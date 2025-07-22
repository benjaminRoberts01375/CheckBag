import React, { ReactNode } from "react";
import { Context, ContextType, CookieKeys } from "./context-object";
import { useState } from "react";
import Service from "./types/service.tsx";

interface Props {
	children: ReactNode;
}

export const ContextProvider: React.FC<Props> = ({ children }) => {
	const [services, setServices] = useState<Service[]>(new Array<Service>());

	function cookieGet(key: CookieKeys): string | undefined {
		const cookieString = document.cookie.split("; ").find(cookie => cookie.startsWith(`${key}=`));

		if (cookieString) {
			return decodeURIComponent(cookieString.split("=")[1]);
		}
		return undefined;
	}

	/**Get a specific service's data from the backend.
	Overwrites the service's data if it already exists while maintaining the client ID.
	Creates a new service if it doesn't exist.
	*/
	function requestServiceData(): void {
		const time_steps: string[] = ["hour", "day", "month", "year"];
		time_steps.forEach(time_step => {
			async () => {
				try {
					const url = new URL("/api/service-data", window.location.origin);
					url.searchParams.set("time-step", time_step);
					console.log("Final URL:", url.toString());
					const response = await fetch(url.toString(), {
						method: "GET",
						headers: {
							"Content-Type": "application/json",
						},
						credentials: "include",
					});

					if (!response.ok) {
						throw new Error("Failed to fetch initial data");
					}
					const rawData: Service[] = await response.json();

					if (services.length === 0) {
						for (let i = 0; i < rawData.length; i++) {
							rawData[i].clientID = crypto.randomUUID(); // Generate new client ID
						}
						setServices(rawData);
						return;
					}

					setServices(oldServices => {
						const originalUniqueServices = oldServices.filter(existingService => {
							// Remove duplicate services, maintaining client ID
							for (let i = 0; i < rawData.length; i++) {
								if (existingService.id === rawData[i].id) {
									rawData[i].clientID = existingService.clientID; // Maintain client ID
									return false;
								}
								rawData[i].clientID = crypto.randomUUID(); // Generate new client ID
							}
							return true;
						});
						return [...originalUniqueServices, ...rawData]; // Add new services
					});

					console.log("Successfully fetched initial data");
				} catch (error) {
					console.error("Error fetching initial data:", error);
				}
			};
		});
	}

	function serviceAdd(service: Service): void {
		const updatedServices = [...services, service];
		setServices(updatedServices);
		serverUpdateServices(updatedServices);
	}

	/**A helper function to set the services on the server. */
	async function serverUpdateServices(servicesToSend: Service[] = services) {
		console.log("Sending: " + JSON.stringify(servicesToSend));
		try {
			const response = await fetch("/api/services-set", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(servicesToSend),
				credentials: "include",
			});

			if (!response.ok) {
				throw new Error("Failed to add service");
			}
			console.log("Successfully added service");
		} catch (error) {
			console.error("Error adding service:", error);
		}
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
		services,
		serviceAdd,
		cookieGet,
		requestServiceData,
		passwordReset,
	};

	return <Context.Provider value={contextShape}>{children}</Context.Provider>;
};
