import "../styles.css";
import GraphData from "../types/graph-data";
import { BarChart } from "@mui/x-charts/BarChart";
import { Timescale } from "../types/strings";

interface StackedBarChartProps {
	graphData: GraphData[];
	timescale: Timescale;
}

const StackedBarChart = ({ graphData, timescale }: StackedBarChartProps) => {
	const uniqueXValues = Array.from(new Set(graphData.map(d => d.xValue.getTime())))
		.sort((a, b) => a - b)
		.map(time => new Date(time));

	const uniqueTitles = Array.from(new Set(graphData.map(d => d.title)));

	const series = uniqueTitles.map(title => {
		const dataForTitle = uniqueXValues.map(xVal => {
			// Find all data points for the current title and xValue
			const relevantDataPoints = graphData.filter(
				d => d.title === title && d.xValue.getTime() === xVal.getTime(),
			);
			// Sum the 'data' for stacking at this xValue for this title
			return relevantDataPoints.reduce((sum, current) => sum + current.data, 0);
		});

		return {
			data: dataForTitle,
			label: title, // Label for the legend and tooltip
			stack: "total", // All series with 'total' stack together
			type: "bar" as const, // Explicitly declare type for series
		};
	});

	const formatXAxis = (date: Date) => {
		switch (timescale) {
			case "hour":
				return date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
			case "day":
				return date.toLocaleTimeString([], { hour: "numeric", hour12: true });
			case "month":
				return date.toLocaleString("default", { month: "short", day: "numeric" }); // Changed for Month Day format
			case "year":
				return date.toLocaleString("default", { month: "short" });
			default:
				return date.toLocaleDateString(); // Default case
		}
	};

	return (
		<div>
			<BarChart
				xAxis={[{ scaleType: "band", data: uniqueXValues, valueFormatter: formatXAxis }]}
				series={series}
				height={300} // Set a height for the chart
				// Adjust margins more for date labels
				grid={{ vertical: true }} // Add vertical grid lines
			/>
		</div>
	);
};

export default StackedBarChart;
