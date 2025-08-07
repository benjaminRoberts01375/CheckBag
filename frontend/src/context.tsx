import React, { ReactNode, useState, useCallback, useEffect } from "react";
import { v4 as uuidv4 } from "uuid";
import { Context, ContextType, ProcessedChartData } from "./context-object";
import Service from "./types/service.tsx";
import { CookieKeys, Timescale } from "./types/strings";
import { useNavigate } from "react-router-dom";
import ChartData from "./types/chart-data";
import ResourceUsageData from "./types/resource-usage-data";

interface Props {
	children: ReactNode;
}

const emptyChartData: ProcessedChartData = {
	quantityData: [],
	responseCodeData: [],
	countryCodeData: [],
	IPAddressData: [],
	resourceUsage: [],
};

export const ContextProvider: React.FC<Props> = ({ children }) => {
	const [services, setServices] = useState<Service[]>(new Array<Service>());
	const [timescale, setTimescale] = useState<Timescale>("hour");
	const navigate = useNavigate();

	// Chart data states for each timespan - now cached and only recalculated when needed
	const [hourData, setHourData] = useState<ProcessedChartData>(emptyChartData);
	const [dayData, setDayData] = useState<ProcessedChartData>(emptyChartData);
	const [monthData, setMonthData] = useState<ProcessedChartData>(emptyChartData);
	const [yearData, setYearData] = useState<ProcessedChartData>(emptyChartData);

	function cookieGet(key: CookieKeys): string | undefined {
		const cookieString = document.cookie.split("; ").find(cookie => cookie.startsWith(`${key}=`));

		if (cookieString) {
			return decodeURIComponent(cookieString.split("=")[1]);
		}
		return undefined;
	}

	// Combines pre-processed service data for a specific timescale
	const combinePreProcessedData = useCallback(
		(targetTimescale: Timescale): ProcessedChartData => {
			const enabledServices = services.filter(service => service.enabled);
			if (enabledServices.length === 0) {
				return emptyChartData;
			}

			const quantityData = enabledServices.map(
				service => service.getProcessedData(targetTimescale).quantityData,
			);
			const responseCodesCounter = new Map<number, number>();
			const countryCounter = new Map<string, number>();
			const ipCounter = new Map<string, number>();
			const resourceUsage = new Array<ResourceUsageData>();

			// Combine data from all enabled services
			enabledServices.forEach(service => {
				const serviceData = service.getProcessedData(targetTimescale);

				// Combine response codes
				serviceData.responseCodeData.forEach((value, key) => {
					responseCodesCounter.set(key, (responseCodesCounter.get(key) ?? 0) + value);
				});

				// Combine countries
				serviceData.countryCodeData.forEach((value, key) => {
					countryCounter.set(key, (countryCounter.get(key) ?? 0) + value);
				});

				// Combine IP addresses
				serviceData.ipAddressData.forEach((value, key) => {
					ipCounter.set(key, (ipCounter.get(key) ?? 0) + value);
				});

				// Add resource usage data
				serviceData.resourceUsage.forEach((value, key) => {
					resourceUsage.push(new ResourceUsageData(service.title, key, value));
				});
			});

			// Create chart data for response codes
			const responseCodes = Array.from(responseCodesCounter.entries())
				.map(([key, value]) => new ChartData(value, String(key)))
				.sort((a, b) => (b.label < a.label ? 1 : -1));

			// Create chart data for countries (top 10 + others)
			const countries = Array.from(countryCounter.entries()).sort((a, b) => b[1] - a[1]);
			const topCountries = countries.slice(0, 10);
			const otherCountriesCount = countries.slice(10).reduce((sum, current) => sum + current[1], 0);
			const countryData = topCountries
				.map(([key, value]) => new ChartData(value, key))
				.sort((a, b) => (b.label < a.label ? 1 : -1));
			if (otherCountriesCount > 0) {
				countryData.push(new ChartData(otherCountriesCount, "Other"));
			}

			// Create chart data for IP addresses (top 10 + others)
			const IPs = Array.from(ipCounter.entries()).sort((a, b) => b[1] - a[1]);
			const topIPs = IPs.slice(0, 10);
			const otherIPsCount = IPs.slice(10).reduce((sum, current) => sum + current[1], 0);
			const ipData = topIPs
				.map(([key, value]) => new ChartData(value, key))
				.sort((a, b) => (b.label < a.label ? 1 : -1));
			if (otherIPsCount > 0) {
				ipData.push(new ChartData(otherIPsCount, "Other"));
			}

			return {
				quantityData,
				responseCodeData: responseCodes,
				countryCodeData: countryData,
				IPAddressData: ipData,
				resourceUsage: resourceUsage.sort((a, b) => b.quantity - a.quantity),
			};
		},
		[services],
	);

	// Update all timespan data when services change
	useEffect(() => {
		setHourData(combinePreProcessedData("hour"));
		setDayData(combinePreProcessedData("day"));
		setMonthData(combinePreProcessedData("month"));
		setYearData(combinePreProcessedData("year"));
	}, [combinePreProcessedData]);

	// Get data for current timescale
	const getCurrentTimescaleData = useCallback((): ProcessedChartData => {
		switch (timescale) {
			case "hour":
				return hourData;
			case "day":
				return dayData;
			case "month":
				return monthData;
			case "year":
				return yearData;
			default:
				return emptyChartData;
		}
	}, [timescale, hourData, dayData, monthData, yearData]);

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
						throw new Error("Failed to fetch initial data:" + response.status);
					}
					const newServices: Service[] = (await response.json()).map((serviceData: any) =>
						Service.fromJSON(serviceData),
					);
					setServices(existingServices => {
						var updatedServices = [...existingServices]; // Make a copy of the existing services, not just a reference to avoid state issues

						for (const newService of newServices) {
							var existingServiceIndex = updatedServices.findIndex(
								existingService => existingService.id === newService.id,
							);

							// If the service doesn't exist, add it
							if (existingServiceIndex === -1) {
								newService.clientID = uuidv4();
								newService.enabled = true;
								updatedServices.push(newService);
							} else {
								// If the service exists, update its processed data
								const existingService = updatedServices[existingServiceIndex];
								switch (time_step) {
									case "hour":
										existingService.hourProcessed = newService.hourProcessed;
										break;
									case "day":
										existingService.dayProcessed = newService.dayProcessed;
										break;
									case "month":
										existingService.monthProcessed = newService.monthProcessed;
										break;
									case "year":
										existingService.yearProcessed = newService.yearProcessed;
										break;
								}
								updatedServices[existingServiceIndex] = existingService;
							}
						}

						if (updatedServices.length === 0) {
							navigate("/dashboard/services");
						}
						return updatedServices;
					});
				} catch (error) {
					console.error("Error fetching initial data:", error);
					navigate("/signin");
				}
			})();
		});
	}

	function serviceAdd(service: Service): void {
		const updatedServices = [...services, service];
		setServices(updatedServices);
		serverUpdateServices(updatedServices);
	}

	function serviceDelete(serviceID: string): void {
		const updatedServices = services.filter(service => service.clientID !== serviceID);
		(async () => {
			await serverUpdateServices(updatedServices);
			requestServiceData();
		})();
	}

	function serviceUpdate(service: Service): void {
		const updatedServices = services.map(existingService => {
			return existingService.clientID === service.clientID ? service : existingService;
		});
		(async () => {
			await serverUpdateServices(updatedServices);
			requestServiceData();
		})();
	}

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
				throw new Error("Failed to modify services - " + response.status);
			}
			console.log("Successfully modified services");
		} catch (error) {
			console.error("Error modifying service:", error);
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

	// Toggle the "enabled" state of a service
	function serviceToggle(serviceID: string): void {
		setServices(services => {
			return services.map(existingService => {
				if (existingService.clientID === serviceID) {
					return new Service(
						existingService.internal_address,
						existingService.external_address,
						existingService.title,
						existingService.id,
						!existingService.enabled,
					);
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
		serviceDelete,
		serviceUpdate,
		cookieGet,
		requestServiceData,
		passwordReset,
		serviceToggle,
		hourData,
		dayData,
		monthData,
		yearData,
		getCurrentTimescaleData,
	};

	return <Context.Provider value={contextShape}>{children}</Context.Provider>;
};
