import "../styles.css";
import GraphData from "../types/graph-data";
import { BarChart } from "@mui/x-charts/BarChart";
import { Timescale } from "../types/strings";
import { ThemeProvider, createTheme } from "@mui/material/styles";

interface StackedBarChartProps {
	graphData: GraphData[];
	timescale: Timescale;
}

const darkTheme = createTheme({
	palette: {
		mode: "dark",
		text: {
			primary: "#aaa",
		},
	},
});

const StackedBarChart = ({ graphData, timescale }: StackedBarChartProps) => {
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
		<ThemeProvider theme={darkTheme}>
			<div>
				<BarChart
					xAxis={[
						{
							scaleType: "band",
							data: uniqueXValues,
							valueFormatter: formatXAxis,
							tickPlacement: "middle",
						},
					]}
					series={series}
					height={300}
					grid={{ vertical: true, horizontal: true }}
					sx={{
						"& .MuiChartsGrid-line": {
							stroke: "#000000",
							strokeDasharray: "2 2",
							opacity: 1,
						},
					}}
				/>
			</div>
		</ThemeProvider>
	);
};

export default StackedBarChart;
