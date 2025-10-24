import React from "react";
import Service from "./types/service.tsx";
import { CookieKeys, Timescale } from "./types/strings";
import GraphData from "./types/graph-data";
import ChartData from "./types/chart-data";
import ResourceUsageData from "./types/resource-usage-data";
import ApiKey from "./types/api-key.tsx";

export interface ProcessedChartData {
	quantityData: GraphData[]; // Per service
	responseCodeData: ChartData[];
	countryCodeData: ChartData[];
	IPAddressData: ChartData[];
	resourceUsage: ResourceUsageData[];
}

export interface ContextType {
	services: Service[];
	apiKeys: ApiKey[];
	addAPIKey: (name: string) => void;
	removeAPIKey: (keyID: string) => void;
	timescale: Timescale;
	setTimescale: (timescale: Timescale) => void;
	serviceAdd: (service: Service) => void;
	serviceDelete: (serviceID: string) => void;
	serviceUpdate: (service: Service) => void;
	cookieGet: (key: CookieKeys) => string | undefined;
	passwordReset: (newPassword: string) => void;
	serviceToggle: (serviceID: string) => void;

	// New chart data states for each time span
	hourData: ProcessedChartData;
	dayData: ProcessedChartData;
	monthData: ProcessedChartData;
	yearData: ProcessedChartData;

	// Function to get data for current timescale
	getCurrentTimescaleData: () => ProcessedChartData;

	// Sign in/out
	signIn: () => void;
}

export const Context = React.createContext<ContextType | undefined>(undefined);
