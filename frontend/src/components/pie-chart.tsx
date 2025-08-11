import "../styles.css";
import graphsStyles from "./graphs.module.css";
import ChartData from "../types/chart-data";
import { PieChart } from "@mui/x-charts/PieChart";
import html2canvas from "html2canvas";
import { useRef } from "react";
import { IoIosDownload } from "react-icons/io";

interface PieChartProps {
	data: ChartData[];
	title: string;
}

const PieChartComponent = ({ data, title }: PieChartProps) => {
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
		<div ref={chartRef} id={graphsStyles["container"]}>
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
			<div id={graphsStyles["chart"]}>
				<PieChart
					series={[
						{
							data: data,
							highlightScope: { highlight: "item" },
						},
					]}
					width={175}
					height={175}
					margin={{ top: 0, right: 0, left: 0, bottom: 0 }}
				/>
			</div>
		</div>
	);
};

export default PieChartComponent;
