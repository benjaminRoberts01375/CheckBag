import React, { ReactNode } from "react";
import { Context, ContextType } from "./context-object";
import { useState } from "react";
import Service from "./types/service.tsx";
import { CookieKeys, Timescale } from "./types/strings";

interface Props {
	children: ReactNode;
}

export const ContextProvider: React.FC<Props> = ({ children }) => {
	const [services, setServices] = useState<Service[]>(new Array<Service>());
	const [timescale, setTimescale] = useState<Timescale>("hour"); // TODO: I'm not a huge fan of this being here, but I'm short on time

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
			(async () => {
				try {
					const url = new URL("/api/service-data", window.location.origin);
					url.searchParams.set("time-step", time_step);
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
					const newServices: Service[] = (await response.json()).map((serviceData: any) =>
						Service.fromJSON(serviceData),
					);
					setServices(existingServices => {
						var finalServices: Service[] = [];
						for (const newService of newServices) {
							var finalService = existingServices.find(
								existingService => existingService.id === newService.id,
							);

							// If the service doesn't exist, add it
							if (finalService === undefined) {
								newService.clientID = crypto.randomUUID();
								newService.enabled = true;
								finalServices.push(newService);
								break;
							}
							// If the service does exist, update it
							switch (time_step) {
								case "hour":
									finalService.hour = newService.hour;
									finalServices.push(finalService);
									break;
								case "day":
									finalService.day = newService.day;
									finalServices.push(finalService);
									break;
								case "month":
									finalService.month = newService.month;
									finalServices.push(finalService);
									break;
								case "year":
									finalService.year = newService.year;
									finalServices.push(finalService);
									break;
							}
						}
						return finalServices;
					});
				} catch (error) {
					console.error("Error fetching initial data:", error);
				}
			})();
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

	// Toggle the service's enabled state and update setServices
	function serviceToggle(serviceID: string): void {
		setServices(services => {
			return services.map(existingService => {
				if (existingService.clientID === serviceID) {
					return { ...existingService, enabled: !existingService.enabled };
				}
				return existingService;
			});
		});
	}

	const contextShape: ContextType = {
		services,
		timescale,
		setTimescale,
		serviceAdd,
		cookieGet,
		requestServiceData,
		passwordReset,
		serviceToggle,
	};

	return <Context.Provider value={contextShape}>{children}</Context.Provider>;
};
