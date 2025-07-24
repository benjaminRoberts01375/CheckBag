import "../styles.css";
import servicesStyles from "./services.module.css";
import { useList } from "../context-hook";
import { useEffect, useState } from "react";
import GraphData from "../types/graph-data";
import ChartData from "../types/chart-data";
import StackedBarChart from "../components/stacked-bar-chart";
import PieChartComponent from "../components/pie-chart";
import { createTheme } from "@mui/material/styles";

const DashboardScreen = () => {
	const { services, requestServiceData, timescale } = useList();
	const [quantityData, setQuantityData] = useState<GraphData[]>([]);
	const [responseCodeData, setResponseCodeData] = useState<ChartData[]>([]);
	const [countryCodeData, setCountryCodeData] = useState<ChartData[]>([]);
	const [IPAddressData, setIPAddressData] = useState<ChartData[]>([]);

	const theme = createTheme({
		palette: {
			mode: "light",
			text: {
				primary: "#aaa",
			},
		},
	});

	useEffect(() => {
		requestServiceData();
	}, []);

	useEffect(() => {
		var graphData: GraphData[] = [];

		for (const service of services.filter(service => service.enabled)) {
			var analyticsMap = service[timescale]; // Woo bracket notation
			var timeStepQuantity = 0;
			var rollback: (step: number) => Date;

			switch (timescale) {
				case "hour":
					rollback = (step: number) => {
						const now = new Date();
						const currentMinute = now.getMinutes();
						now.setMinutes(currentMinute + step);
						now.setSeconds(0);
						now.setMilliseconds(0);
						return now;
					};
					timeStepQuantity = 60;
					break;
				case "day":
					rollback = (step: number) => {
						const now = new Date();
						const currentHour = now.getHours();
						now.setHours(currentHour + step);
						now.setMinutes(0);
						now.setSeconds(0);
						now.setMilliseconds(0);
						return now;
					};
					timeStepQuantity = 24;
					break;
				case "month":
					rollback = (step: number): Date => {
						var now = new Date();
						const currentUTCDay = now.getUTCDate();
						now.setUTCDate(currentUTCDay + step);
						now.setUTCHours(0);
						now.setUTCMinutes(0);
						now.setUTCSeconds(0);
						now.setUTCMilliseconds(0);
						return now;
					};
					timeStepQuantity = 30;
					break;
				case "year":
					rollback = (step: number) => {
						const now = new Date();
						const currentMonth = now.getMonth();
						now.setUTCMonth(currentMonth + step);
						now.setUTCDate(1);
						now.setUTCHours(0);
						now.setUTCMinutes(0);
						now.setUTCSeconds(0);
						now.setUTCMilliseconds(0);
						return now;
					};
					timeStepQuantity = 12;
					break;
			}

			var responseCodesCounter = new Map<number, number>(); // Map of error codes to quantity
			var countryCounter = new Map<string, number>(); // Map of countries to quantity
			var ipCounter = new Map<string, number>(); // Map of IP addresses to quantity
			for (let i = 0; i < timeStepQuantity; i++) {
				const date = rollback(-i);
				var analytic = analyticsMap.get(date.toISOString());
				graphData.push(new GraphData(analytic?.quantity ?? 0, service.title, date));
				if (analytic !== undefined) {
					// Count response codes
					analytic.responseCode.forEach((value, key) => {
						responseCodesCounter.set(key, responseCodesCounter.get(key) ?? 0 + value);
					});
					// Count countries
					analytic.country.forEach((value, key) => {
						countryCounter.set(key, countryCounter.get(key) ?? 0 + value);
					});
					// Count IP addresses
					analytic.ip.forEach((value, key) => {
						ipCounter.set(key, ipCounter.get(key) ?? 0 + value);
					});
				}
			}
			// Create chart data for response codes
			const responseCodes = Array.from(responseCodesCounter.entries()).map(([key, value]) => {
				return new ChartData(value, String(key));
			});
			// Create chart data for countries
			const countries = Array.from(countryCounter.entries()).map(([key, value]) => {
				return new ChartData(value, key);
			});
			// Create chart data for top 10 IP addresses and other
			const IPs = Array.from(ipCounter.entries()).sort((a, b) => b[1] - a[1]);
			const topIPs = IPs.slice(0, 10);
			const otherIPsCount = IPs.slice(10).reduce((sum, current) => sum + current[1], 0);
			const ipData = topIPs.map(([key, value]) => {
				return new ChartData(value, key);
			});
			ipData.push(new ChartData(otherIPsCount, "Other"));

			// Update state
			setCountryCodeData(countries);
			setResponseCodeData(responseCodes);
			setIPAddressData(ipData);
		}
		setQuantityData(graphData);
	}, [timescale, services]);

	return (
		<div id={servicesStyles["container"]}>
			<StackedBarChart
				graphData={quantityData}
				timescale={timescale}
				yAxisLabel="Query Quantity"
				title="Query Quantity Per Service"
				theme={theme}
			/>
			<div id={servicesStyles["pie-charts"]}>
				<PieChartComponent data={responseCodeData} title="Response Codes" theme={theme} />
				<PieChartComponent data={countryCodeData} title="Countries" theme={theme} />
				<PieChartComponent data={IPAddressData} title="Top IP Addresses" theme={theme} />
			</div>
		</div>
	);
};

export default DashboardScreen;
