import "../styles.css";
import GraphData from "../types/graph-data";
import graphsStyles from "./graphs.module.css";
import { IoIosDownload } from "react-icons/io";
import { BarChart } from "@mui/x-charts/BarChart";
import html2canvas from "html2canvas";
import { useRef } from "react";

interface StackedBarChartProps {
	graphData: GraphData[];
	yAxisLabel: string;
	title: string;
}

const StackedBarChart = ({ graphData, yAxisLabel, title }: StackedBarChartProps) => {
	const chartRef = useRef<HTMLDivElement>(null);

	const exportToPNG = async () => {
		if (chartRef.current) {
			try {
				// Find and hide the button temporarily
				const saveButton = chartRef.current.querySelector(
					"#" + graphsStyles["save-button"],
				) as HTMLElement;
				if (saveButton) {
					saveButton.style.display = "none";
				}

				const canvas = await html2canvas(chartRef.current, {
					backgroundColor: null,
					scale: 5,
					useCORS: true,
				});

				// Show the button again
				if (saveButton) {
					saveButton.style.display = "";
				}

				const link = document.createElement("a");
				link.download = `${title.toLowerCase().replace(/\s+/g, "-")}-chart.png`;
				link.href = canvas.toDataURL("image/png");
				link.click();
			} catch (error) {
				console.error("Error exporting chart:", error);
				// Make sure button is visible again even if there's an error
				const saveButton = chartRef.current.querySelector(
					"#" + graphsStyles["save-button"],
				) as HTMLElement;
				if (saveButton) {
					saveButton.style.display = "";
				}
			}
		}
	};
	return (
		<div ref={chartRef}>
			<div id={graphsStyles["header"]}>
				<h2>{title}</h2>
				<button
					onClick={exportToPNG}
					id={graphsStyles["save-button"]}
					className="secondary"
					title="Download chart"
				>
					<IoIosDownload />
				</button>
			</div>
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
				skipAnimation={true}
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
