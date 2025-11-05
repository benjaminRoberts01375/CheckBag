import "../styles.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import StackedBarChart from "../components/stacked-bar-chart";
import PieChartComponent from "../components/pie-chart";
import ResourceTable from "../components/resource-table";
import { createTheme } from "@mui/material/styles";
import { ThemeProvider } from "@mui/material/styles";

// Create theme once outside component - it never changes
const theme = createTheme({
	palette: {
		mode: "light",
		text: {
			primary: "#aaa",
		},
	},
});

const DashboardScreen = () => {
	const { getCurrentTimescaleData } = useList();

	// Get the processed chart data for the current timescale
	const chartData = getCurrentTimescaleData();

	return (
		<div id={DashboardStyles["container"]}>
			<title>CheckBag - Dashboard</title>
			<ThemeProvider theme={theme}>
				<div id={DashboardStyles["summary-group"]}>
					<div className={DashboardStyles["summary-card"]}>
						<h2>Temp</h2>
					</div>
					<div className={DashboardStyles["summary-card"]}>
						<h2>Temp</h2>
					</div>
					<div className={DashboardStyles["summary-card"]}>
						<h2>Temp</h2>
					</div>
					<div className={DashboardStyles["summary-card"]}>
						<h2>Temp</h2>
					</div>
				</div>
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
