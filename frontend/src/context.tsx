import React, { ReactNode, useState, useCallback, useEffect } from "react";
import { v4 as uuidv4 } from "uuid";
import { Context, ContextType, ProcessedChartData } from "./context-object";
import Service from "./types/service.tsx";
import { CookieKeys, Timescale } from "./types/strings";
import { useNavigate } from "react-router-dom";
import ChartData from "./types/chart-data";
import ResourceUsageData from "./types/resource-usage-data";
import APIKey from "./types/api-key.tsx";

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

	// Chart data states for each time span - now cached and only recalculated when needed
	const [hourData, setHourData] = useState<ProcessedChartData>(emptyChartData);
	const [dayData, setDayData] = useState<ProcessedChartData>(emptyChartData);
	const [monthData, setMonthData] = useState<ProcessedChartData>(emptyChartData);
	const [yearData, setYearData] = useState<ProcessedChartData>(emptyChartData);

	// API keys
	const [APIKeys, setAPIKeys] = useState<APIKey[]>([]);

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
				let regionNames = new Intl.DisplayNames(["en"], { type: "region" });
				serviceData.countryCodeData.forEach((value, key) => {
					const regionName = regionNames.of(key);
					countryCounter.set(
						regionName ?? key,
						(countryCounter.get(regionName ?? key) ?? 0) + value,
					);
				});

				// Combine IP addresses
				serviceData.ipAddressData.forEach((value, key) => {
					ipCounter.set(key, (ipCounter.get(key) ?? 0) + value);
				});

				// Add resource usage data
				serviceData.resourceUsage.forEach((value, key) => {
					resourceUsage.push(
						new ResourceUsageData(service.title, service.internal_address, key, value),
					);
				});
			});

			// Create chart data for response codes
			const responseCodes = Array.from(responseCodesCounter.entries())
				.map(
					([key, value]) =>
						new ChartData(value, (statusCodeToString[key] ?? key) + " (" + value + ")"),
				)
				.sort((a, b) => (b.label < a.label ? 1 : -1));

			// Create chart data for countries (top 10 + others)
			const countries = Array.from(countryCounter.entries()).sort((a, b) => b[1] - a[1]);
			const topCountries = countries.slice(0, 10);
			const otherCountriesCount = countries.slice(10).reduce((sum, current) => sum + current[1], 0);
			const countryData = topCountries.map(
				([key, value]) => new ChartData(value, key + " (" + value + ")"),
			);
			if (otherCountriesCount > 0) {
				const totalOtherCountries = countries.slice(10).length;
				var label = "Others";
				if (totalOtherCountries == 1) {
					label = "Other";
				}
				label = "+" + totalOtherCountries + " " + label + " (" + otherCountriesCount + ")";
				countryData.push(new ChartData(otherCountriesCount, label));
			}

			// Create chart data for IP addresses (top 10 + others)
			const IPs = Array.from(ipCounter.entries()).sort((a, b) => b[1] - a[1]);
			const topIPs = IPs.slice(0, 10);
			const otherIPsCount = IPs.slice(10).reduce((sum, current) => sum + current[1], 0);
			const ipData = topIPs.map(([key, value]) => new ChartData(value, key + " (" + value + ")"));
			if (otherIPsCount > 0) {
				const totalOtherIPs = IPs.slice(10).length;
				var label = "Others";
				if (totalOtherIPs == 1) {
					label = "Other";
				}
				label = "+" + totalOtherIPs + " " + label + " (" + otherIPsCount + ")";
				ipData.push(new ChartData(otherIPsCount, label));
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

	// Update all time span data when services change
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

	function requestAPIKeys(): void {
		(async () => {
			try {
				const response = await fetch("/api/api-keys", {
					method: "GET",
					headers: {
						"Content-Type": "application/json",
					},
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to fetch API keys:" + response.status);
				}

				const newAPIKeys = await response.json();
				setAPIKeys(newAPIKeys);
			} catch (error) {
				console.error("Error fetching API keys:", error);
			}
		})();
	}

	function updateServerAPIKeys(keys: APIKey[]): void {
		(async () => {
			try {
				const response = await fetch("/api/api-keys", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify(keys),
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to set API keys:" + response.status);
				}
				response.json().then(newKeys => {
					setAPIKeys(_ => {
						return newKeys;
					});
				});
				console.log("Successfully set API keys");
			} catch (error) {
				console.error("Error setting API keys:", error);
			}
		})();
	}

	function addAPIKey(name: string): void {
		updateServerAPIKeys([...APIKeys, new APIKey(name, uuidv4())]);
	}

	function removeAPIKey(key_id: string): void {
		updateServerAPIKeys(APIKeys.filter(key => key.id !== key_id));
	}

	function requestServiceData(): void {
		const time_steps: string[] = ["hour", "day", "month", "year"];

		// Create all fetch promises
		const fetchPromises = time_steps.map(async time_step => {
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

				return { time_step, services: newServices, success: true };
			} catch (error) {
				console.error(`Error fetching data for ${time_step}:`, error);
				return { time_step, services: [], success: false };
			}
		});

		// Wait for all requests to complete and process results
		Promise.allSettled(fetchPromises).then(results => {
			const successfulResults = results
				.filter(result => result.status === "fulfilled")
				.map(result => result.value);

			// Check if any requests succeeded
			if (successfulResults.length === 0) {
				console.error("All service data requests failed");
				navigate("/signin");
				return;
			}

			// Update services state once with all the collected data
			setServices(existingServices => {
				var updatedServices = [...existingServices];

				// Process each successful result
				successfulResults.forEach(({ time_step, services: newServices }) => {
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
				});

				if (updatedServices.length === 0) {
					navigate("/dashboard/services");
				}
				return updatedServices;
			});
			requestAPIKeys();
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
					const updated_service = new Service(
						existingService.internal_address,
						existingService.external_address,
						existingService.title,
						existingService.id,
						!existingService.enabled,
					);
					updated_service.hourProcessed = existingService.hourProcessed;
					updated_service.dayProcessed = existingService.dayProcessed;
					updated_service.monthProcessed = existingService.monthProcessed;
					updated_service.yearProcessed = existingService.yearProcessed;
					return updated_service;
				}
				return existingService;
			});
		});
	}

	const statusCodeToString: Record<number, string> = {
		// 1xx Informational
		100: "100: Continue",
		101: "101: Switching Protocols",
		102: "102: Processing",
		103: "103: Early Hints",

		// 2xx Success
		200: "200: Ok",
		201: "201: Created",
		202: "202: Accepted",
		203: "203: Non-Authoritative Information",
		204: "204: No Content",
		205: "205: Reset Content",
		206: "206: Partial Content",
		207: "207: Multi-Status",
		208: "208: Already Reported",
		226: "226: IM Used",

		// 3xx Redirection
		300: "300: Multiple Choices",
		301: "301: Moved Permanently",
		302: "302: Found",
		303: "303: See Other",
		304: "304: Not Modified",
		305: "305: Use Proxy",
		306: "306: Switch Proxy",
		307: "307: Temporary Redirect",
		308: "308: Permanent Redirect",

		// 4xx Client Error
		400: "400: Bad Request",
		401: "401: Unauthorized",
		402: "402: Payment Required",
		403: "403: Forbidden",
		404: "404: Not Found",
		405: "405: Method Not Allowed",
		406: "406: Not Acceptable",
		407: "407: Proxy Authentication Required",
		408: "408: Request Timeout",
		409: "409: Conflict",
		410: "410: Gone",
		411: "411: Length Required",
		412: "412: Precondition Failed",
		413: "413: Payload Too Large",
		414: "414: URI Too Long",
		415: "415: Unsupported Media Type",
		416: "416: Range Not Satisfiable",
		417: "417: Expectation Failed",
		418: "418: I'm a teapot", // :)
		421: "421: Misdirected Request",
		422: "422: Unprocessable Entity",
		423: "423: Locked",
		424: "424: Failed Dependency",
		425: "425: Too Early",
		426: "426: Upgrade Required",
		428: "428: Precondition Required",
		429: "429: Too Many Requests",
		431: "431: Request Header Fields Too Large",
		451: "451: Unavailable For Legal Reasons",

		// 5xx Server Error
		500: "500: Internal Server Error",
		501: "501: Not Implemented",
		502: "502: Bad Gateway",
		503: "503: Service Unavailable",
		504: "504: Gateway Timeout",
		505: "505: HTTP Version Not Supported",
		506: "506: Variant Also Negotiates",
		507: "507: Insufficient Storage",
		508: "508: Loop Detected",
		510: "510: Not Extended",
		511: "511: Network Authentication Required",

		// Non-standard - NGINX
		444: "444: No Response",
		494: "494: Request Header Too Large",
		495: "495: Certificate Error",
		496: "496: No Certificate",
		497: "497: HTTP Request Sent to HTTPS Port",
		499: "499: Client Closed Request",

		// Non-standard - Cloudflare
		520: "520: Unknown Error",
		521: "521: Server is Down",
		522: "522: Connection Timed Out",
		523: "523: Origin Is Unreachable",
		524: "524: A Timeout Occurred",
		525: "525: SSL Handshake Failed",
		526: "526: Invalid SSL Certificate",
	} as const;

	const contextShape: ContextType = {
		services,
		apiKeys: APIKeys,
		addAPIKey,
		removeAPIKey,
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
