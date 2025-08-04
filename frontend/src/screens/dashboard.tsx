import "../styles.css";
import servicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import { useEffect } from "react";
import StackedBarChart from "../components/stacked-bar-chart";
import PieChartComponent from "../components/pie-chart";
import ResourceTable from "../components/resource-table";
import { createTheme } from "@mui/material/styles";
import { ThemeProvider } from "@mui/material/styles";

const DashboardScreen = () => {
	const { requestServiceData, getCurrentTimescaleData, services } = useList();

	// Get the processed chart data for the current timescale
	const chartData = getCurrentTimescaleData();

	const theme = createTheme({
		palette: {
			mode: "light",
			text: {
				primary: "#aaa",
			},
		},
	});

	useEffect(() => {
		// Initial request
		if (services.length === 0) {
			requestServiceData();
		}
		// Update every 10 seconds after initial request
		const interval = setInterval(requestServiceData, 10000);
		return () => clearInterval(interval);
	}, []);

	return (
		<div id={servicesStyles["container"]}>
			<ThemeProvider theme={theme}>
				<div className={DashboardStyles["graph-group"]}>
					<StackedBarChart
						graphData={chartData.quantityData}
						yAxisLabel="Query Quantity"
						title="Query Quantity Per Service"
					/>
				</div>
				<div id={DashboardStyles["pie-charts"]} className={DashboardStyles["graph-group"]}>
					<PieChartComponent data={chartData.responseCodeData} title="Response Codes" />
					<PieChartComponent data={chartData.countryCodeData} title="Top Countries" />
					<PieChartComponent data={chartData.IPAddressData} title="Top IP Addresses" />
				</div>
				<div className={DashboardStyles["graph-group"]}>
					<ResourceTable data={chartData.resourceUsage} title="Resource Usage" />
				</div>
			</ThemeProvider>
		</div>
	);
};

export default DashboardScreen;
