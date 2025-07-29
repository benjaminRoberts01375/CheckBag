import "../styles.css";
import GraphData from "../types/graph-data";
import { BarChart } from "@mui/x-charts/BarChart";

interface StackedBarChartProps {
	graphData: GraphData[];
	yAxisLabel: string;
	title: string;
}

const StackedBarChart = ({ graphData, yAxisLabel, title }: StackedBarChartProps) => {
	return (
		<div>
			<h2 className="header">{title}</h2>
			<BarChart
				xAxis={[
					{
						scaleType: "band",
						data: graphData[0] !== undefined ? graphData[0].x_values : [],
						tickPlacement: "middle",
						label: "Time",
					},
				]}
				yAxis={[
					{
						label: yAxisLabel,
					},
				]}
				series={graphData}
				height={300}
				margin={{ top: 0, right: 0, left: 0, bottom: 0 }}
				grid={{ vertical: true, horizontal: true }}
				sx={{
					"& .MuiChartsGrid-line": {
						stroke: "#000",
						strokeDasharray: "5 5",
						opacity: 1,
					},
				}}
			/>
		</div>
	);
};

export default StackedBarChart;
