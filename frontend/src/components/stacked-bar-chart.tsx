import "../styles.css";
import graphsStyles from "./graphs.module.css";
import GraphData from "../types/graph-data";
import { BarChart } from "@mui/x-charts/BarChart";
import { Timescale } from "../types/strings";
import { ThemeProvider, Theme } from "@mui/material/styles";

interface StackedBarChartProps {
	graphData: GraphData[];
	timescale: Timescale;
	yAxisLabel: string;
	title: string;
	theme: Theme;
}

const StackedBarChart = ({
	graphData,
	timescale,
	yAxisLabel,
	title,
	theme,
}: StackedBarChartProps) => {
	const uniqueXValues = Array.from(new Set(graphData.map(d => d.xValue.getTime())))
		.sort((a, b) => a - b)
		.map(time => new Date(time));

	const uniqueTitles = Array.from(new Set(graphData.map(d => d.title)));

	const series = uniqueTitles.map(title => {
		const dataForTitle = uniqueXValues.map(xVal => {
			const relevantDataPoints = graphData.filter(
				d => d.title === title && d.xValue.getTime() === xVal.getTime(),
			);
			return relevantDataPoints.reduce((sum, current) => sum + current.data, 0);
		});

		return {
			data: dataForTitle,
			label: title,
			stack: "total",
			type: "bar" as const,
		};
	});

	const formatXAxis = (date: Date) => {
		switch (timescale) {
			case "hour":
				return date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
			case "day":
				return date.toLocaleTimeString([], { hour: "numeric", hour12: true });
			case "month":
				return date.toLocaleString("default", { month: "short", day: "numeric" });
			case "year":
				return date.toLocaleString("default", { month: "short" });
			default:
				return date.toLocaleDateString();
		}
	};

	return (
		<ThemeProvider theme={theme}>
			<h2 className={graphsStyles["header"]}>{title}</h2>
			<BarChart
				xAxis={[
					{
						scaleType: "band",
						data: uniqueXValues,
						valueFormatter: formatXAxis,
						tickPlacement: "middle",
						label: "Time",
					},
				]}
				yAxis={[
					{
						label: yAxisLabel,
					},
				]}
				series={series}
				height={300}
				margin={{ top: 0, right: 0, left: 0, bottom: 0 }}
				grid={{ vertical: true, horizontal: true }}
				sx={{
					"& .MuiChartsGrid-line": {
						stroke: "#000000",
						strokeDasharray: "5 5",
						opacity: 1,
					},
				}}
			/>
		</ThemeProvider>
	);
};

export default StackedBarChart;
