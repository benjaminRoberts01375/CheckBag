import React, { ReactNode, useState, useCallback, useEffect } from "react";
import { Context, ContextType, ProcessedChartData } from "./context-object";
import Service from "./types/service.tsx";
import { CookieKeys, Timescale } from "./types/strings";
import { useNavigate } from "react-router-dom";
import GraphData from "./types/graph-data";
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

	// Chart data states for each timespan
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

	// Process services data into chart data for a specific timescale
	const processServicesData = useCallback(
		(targetTimescale: Timescale): ProcessedChartData => {
			const enabledServices = services.filter(service => service.enabled);
			if (enabledServices.length === 0) {
				return emptyChartData;
			}

			var graphData: GraphData[] = [];
			var responseCodesCounter = new Map<number, number>();
			var countryCounter = new Map<string, number>();
			var ipCounter = new Map<string, number>();
			var resourceUsage = new Array<ResourceUsageData>();

			// Configure timespan-specific settings
			var timeStepQuantity = 0;
			var rollback: (step: number) => Date;

			switch (targetTimescale) {
				case "hour":
					rollback = (step: number) => {
						const now = new Date();
						const currentMinute = now.getMinutes();
						now.setUTCMinutes(currentMinute + step, 0, 0);
						return now;
					};
					timeStepQuantity = 60;
					break;
				case "day":
					rollback = (step: number) => {
						const now = new Date();
						const currentHour = now.getHours();
						now.setHours(currentHour + step);
						now.setUTCHours(currentHour + step, 0, 0, 0);
						return now;
					};
					timeStepQuantity = 24;
					break;
				case "month":
					rollback = (step: number): Date => {
						var now = new Date();
						now.setUTCHours(0, 0, 0, 0);
						const currentUTCDay = now.getUTCDate();
						now.setUTCDate(currentUTCDay + step);
						return now;
					};
					timeStepQuantity = 30;
					break;
				case "year":
					rollback = (step: number) => {
						const now = new Date();
						now.setUTCDate(1);
						now.setUTCHours(0, 0, 0, 0);
						const currentMonth = now.getUTCMonth();
						now.setUTCMonth(currentMonth + step);
						return now;
					};
					timeStepQuantity = 12;
					break;
			}

			// Process each service
			for (const service of enabledServices) {
				var analyticsMap = service[targetTimescale];
				const usedResource = new Map<string, number>();
				const working_graph_data = new GraphData(service.title, targetTimescale);

				// Generate data points for the entire time range (including empty ones)
				for (let i = 0; i < timeStepQuantity; i++) {
					const date = rollback(-i);
					var dateString = () => {
						switch (targetTimescale) {
							case "hour":
								return date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
							case "day":
								return date.toLocaleTimeString([], { hour: "numeric", hour12: true });
							case "month":
								return date.toLocaleString("default", {
									month: "short",
									day: "numeric",
									timeZone: "UTC",
								});
							case "year":
								return date.toLocaleString("default", {
									month: "short",
									timeZone: "UTC",
								});
							default:
								return date.toLocaleDateString();
						}
					};
					working_graph_data.x_values[timeStepQuantity - i - 1] = dateString(); // Reverse order
					working_graph_data.data[timeStepQuantity - i - 1] =
						analyticsMap.get(date.toISOString())?.quantity ?? 0; // Reverse order
					var analytic = analyticsMap.get(date.toISOString());

					if (analytic !== undefined) {
						// Count response codes
						analytic.responseCode.forEach((value, key) => {
							responseCodesCounter.set(key, (responseCodesCounter.get(key) ?? 0) + value);
						});
						// Count countries
						analytic.country.forEach((value, key) => {
							countryCounter.set(key, (countryCounter.get(key) ?? 0) + value);
						});
						// Count IP addresses
						analytic.ip.forEach((value, key) => {
							ipCounter.set(key, (ipCounter.get(key) ?? 0) + value);
						});
						// Count resources
						analytic.resource.forEach((value, key) => {
							usedResource.set(key, (usedResource.get(key) ?? 0) + value);
						});
					}
				}
				graphData.push(working_graph_data);

				// Add resource usage data for this service
				usedResource.forEach((value, key) => {
					resourceUsage.push(new ResourceUsageData(service.title, key, value));
				});
			}

			// Create chart data for response codes
			const responseCodes = Array.from(responseCodesCounter.entries()).map(([key, value]) => {
				return new ChartData(value, String(key));
			});

			// Create chart data for countries (top 10 + others)
			const countries = Array.from(countryCounter.entries()).sort((a, b) => b[1] - a[1]);
			const topCountries = countries.slice(0, 10);
			const otherCountriesCount = countries.slice(10).reduce((sum, current) => sum + current[1], 0);
			const countryData = topCountries.map(([key, value]) => {
				return new ChartData(value, key);
			});
			if (otherCountriesCount > 0) {
				countryData.push(new ChartData(otherCountriesCount, "Other"));
			}

			// Create chart data for IP addresses (top 10 + others)
			const IPs = Array.from(ipCounter.entries()).sort((a, b) => b[1] - a[1]);
			const topIPs = IPs.slice(0, 10);
			const otherIPsCount = IPs.slice(10).reduce((sum, current) => sum + current[1], 0);
			const ipData = topIPs.map(([key, value]) => {
				return new ChartData(value, key);
			});
			if (otherIPsCount > 0) {
				ipData.push(new ChartData(otherIPsCount, "Other"));
			}

			return {
				quantityData: graphData,
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
		setHourData(processServicesData("hour"));
		setDayData(processServicesData("day"));
		setMonthData(processServicesData("month"));
		setYearData(processServicesData("year"));
	}, [processServicesData]);

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
						if (finalServices.length === 0) {
							navigate("/dashboard/services");
						}
						return finalServices;
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
